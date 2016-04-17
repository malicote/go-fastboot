package transporter

// TODO: use proper logging (debug, fatal, etc)

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unicode"
)

const (
	MAX_RETRIES = 5
	WAIT_FOR_DISCONNECT_TIME = 3
	MAX_USBFS_BULK_SIZE = (16 * 1024)
	USB_SYSFS_BASE = "/sys/bus/usb/devices"
	USB_DEV_BUS_BASE = "/dev/bus/usb"
)

type UsbHandler struct {
	fname string
	desc int32
	ep_in uint8
	ep_out uint8
}

/* Linux USB implementation of Transporter interface */
type LinuxUSB struct {
	UsbHandler
}


func NewLinuxUSB(handler UsbHandler) (*LinuxUSB){
	return &LinuxUSB{ handler }
}

/* Returns true if the buff descriptor header is of type usb_descriptor_type.
 */
func descriptor_is_type(buff *bytes.Reader, usb_descriptor_type uint8) bool {
	descriptor_buff := []uint8{0, 0}
	n_read, err := buff.Read(descriptor_buff)
	if err != nil || n_read != 2 {
		return false
	}

	/* Rollback internal offset */
	err1 := buff.UnreadByte()
	err2 := buff.UnreadByte()

	if err1 != nil || err2 != nil {
		fmt.Println("Error rolling back reader!")
		return false
	}

	/* Per usb_descriptor_header, the type is byte offset 1 */
	return descriptor_buff[1] == usb_descriptor_type
}

func seek_to_next_descriptor_of_type(buff *bytes.Reader, usb_descriptor_type uint8) error {
	for {
		if descriptor_is_type(buff, usb_descriptor_type) {
			return nil
		}

		size, err := buff.ReadByte()
		if err != nil {
			return err
		} else if size == 0 {
			return errors.New("Descriptor reports 0 size.")
		}

		/* Todo: skipping past the end is OK for Seek, do
		 * we need to catch the error here or let the bytes
		 * library do it?
		 */
		_, err = buff.Seek(int64(size) - 1, 1)
		if err != nil{
			return err
		}

		/* Out of buffer */
		if buff.Len() == 0 {
			return errors.New( "Ran out of buffer while parsing.")
		}
	}
}


func filter_usb_device(sysfs_name string,
			dev_data []byte,
			writable bool,
			callback ifc_match_func) (ept_in_id int, ept_out_id int, ifc_id int, err error) {

	/* Parse the USB structure read from /dev and send any device descriptions
	   to the callback to determine if it's a match.
	 */

	var device_descriptor usb_device_descriptor

	buff := bytes.NewReader(dev_data)

	if !descriptor_is_type(buff, USB_DT_DEVICE) {
		return -1, -1, -1, errors.New("Data provided is not a device descriptor.")
	}

	err = binary.Read(buff, binary.LittleEndian, &device_descriptor)
	if err != nil {
		return -1, -1, -1, err
	}

	// fmt.Printf("Data descriptor: %+#v\n", device_descriptor)

	var config_descriptor usb_config_descriptor
	if !descriptor_is_type(buff, USB_DT_CONFIG) {
		return -1, -1, -1, errors.New("Device has no config descriptor.")
	}

	err = binary.Read(buff, binary.LittleEndian, &config_descriptor)
	if err != nil {
		return -1, -1, -1, err
	}

	// fmt.Printf("\nConfig descriptor: %+#v\n", config_descriptor)

	/* Read device serial number (if there is one).
	 * We read the serial number from sysfs, since it's faster and more
	 * reliable than issuing a control pipe read, and also won't
	 * cause problems for devices which don't like getting descriptor
	 * requests while they're in the middle of flashing.
	 */
	serial_number := ""
	if device_descriptor.ISerialNumber != 0 {
		serial_path := path.Join(USB_SYSFS_BASE, sysfs_name, "serial")
		serial, err := ioutil.ReadFile(serial_path)
		if err != nil {
			fmt.Println("Error reading serial number: ", err)
		} else {
			serial_number = strings.TrimSpace(string(serial))
		}
	}


	/* Check each interface */
	for i := uint8(0); i < config_descriptor.BNumInterfaces; i++ {
		err = seek_to_next_descriptor_of_type(buff, USB_DT_INTERFACE)
		if err != nil {
			return -1, -1, -1, err
		}

		var interface_descriptor usb_interface_descriptor

		err = binary.Read(buff, binary.LittleEndian, &interface_descriptor)
		if err != nil {
			return -1, -1, -1, err
		}

		in, out := -1, -1

		for e := uint8(0); e < interface_descriptor.BNumEndpoints; e++ {
			err = seek_to_next_descriptor_of_type(buff, USB_DT_ENDPOINT)
			if err != nil {
				/* OK, just means there's nothing left to parse */
				break
			}

			var endpoint_descriptor usb_endpoint_descriptor
			err = binary.Read(buff, binary.LittleEndian, &endpoint_descriptor)
			if err != nil {
				/* OK, just don't parse this one */
				break
			}

			if endpoint_descriptor.BmAttributes & USB_ENDPOINT_XFERTYPE_MASK != USB_ENDPOINT_XFER_BULK {
				/* Skip, bulk only */
				continue
			}

			if endpoint_descriptor.BEndpointAddress & USB_ENDPOINT_DIR_MASK != 0{
				in = int(endpoint_descriptor.BEndpointAddress)
			} else {
				out = int(endpoint_descriptor.BEndpointAddress)
			}

			/* Skip USB 3.0 SS Endpoint Companion descriptor */
			if descriptor_is_type(buff, USB_DT_SS_ENDPOINT_COMP) {
				_, err = buff.Seek(USB_DT_SS_EP_COMP_SIZE, 1)
				if err != nil {
					/* OK, I think */
					break
				}
			}
		}

		info := UsbIfcInfo{
			DevVendor: 	device_descriptor.IdVendor,
			DevProduct: 	device_descriptor.IdProduct,
			DevClass: 	device_descriptor.BDeviceClass,
			DevSubclass: 	device_descriptor.BDeviceSubClass,
			DevProtocol: 	device_descriptor.BDeviceProtocol,
			Writeable: 	writable,
			SerialNumber:	serial_number,
			DevicePath:	"usb:" + sysfs_name,
			IfcClass:	interface_descriptor.BInterfaceClass,
			IfcSubclass:	interface_descriptor.BInterfaceSubClass,
			IfcProtocol:	interface_descriptor.BInterfaceProtocol,
			HasBulkIn:	(in != -1),
			HasBulkOut:	(out != -1),
		}

		if callback(info) {
			return in, out, int(interface_descriptor.BInterfaceNumber), nil
		}
	}

	return -1, -1, -1, errors.New("No device found.")
}

func convertToDevFSName(dir_name string) (string, error) {
	bus_num_path := path.Join(USB_SYSFS_BASE, dir_name, "busnum")
	dev_num_path := path.Join(USB_SYSFS_BASE, dir_name, "devnum")

	busnum, err := ioutil.ReadFile(bus_num_path)
	if err != nil {
		return "", err
	}

	devnum, err := ioutil.ReadFile(dev_num_path)
	if err != nil {
		return "", err
	}

	devpath := path.Join(USB_DEV_BUS_BASE,
				fmt.Sprintf("%03s", strings.TrimSpace(string(busnum[:]))),
				fmt.Sprintf("%03s", strings.TrimSpace(string(devnum[:]))))

	return  devpath, nil
}

func findUsbDevice(base string, callback ifc_match_func) (*UsbHandler, error) {
	// TODO: can this be simplified?
	is_device := func(path string) bool {
		if strings.HasPrefix(path, ".") {
			// Is ./ or ../
			return false
		}

		for _, c := range path {
			ok := (unicode.IsDigit(c) || c == '.' || c == '-')
			if !ok {
				// Is an interface or a hub
				return false
			}
		}
		return true
	}

	dirs, err := ioutil.ReadDir(base)
	if err != nil {
		fmt.Println("Error reading ", base, ":", err)
	}

	for _, dir := range dirs {
		if !is_device(dir.Name()) {
			continue
		}

		dev_path, err := convertToDevFSName(dir.Name())
		if err != nil {
			// TODO: verbose printing only
			fmt.Println(err)
			continue
		}

		writable := true
		dev_node, err := os.OpenFile(dev_path, os.O_RDWR, 0)
		if err != nil {
			// Check if device node is read-only, then we can provide
			// a helpful error message similar to 'adb devices'
			// TODO: verbose
			fmt.Println(err)
			writable = false
			if dev_node, err = os.OpenFile(dev_path, os.O_RDONLY, 0); err != nil {
				// TODO: verbose
				fmt.Println("Skipping: ", dev_path, "[", err, "]")
				continue
			}
		}

		dev_buff := make([]byte, 1024)
		n_read, err := dev_node.Read(dev_buff)
		if err != nil || n_read == 0 {
			fmt.Printf("Failed to read %v: %v (read %v bytes)\n", dev_node.Name(), err, n_read)
			continue
		}

		_, _, _, err = filter_usb_device(dir.Name(), dev_buff[0:n_read], writable, callback)
		if err != nil {
			/* TODO: debug output */
			// fmt.Printf("Error reading %s: %s\n", dir.Name(), err)
			continue
		}

		fmt.Printf("Devfs path %s, writeable: %v\n", dev_path, writable)
	}

	return new(UsbHandler), nil
}

// TODO: move this to its own package or figure out way to conditionally compile
func UsbOpen(callback ifc_match_func) {
	findUsbDevice(USB_SYSFS_BASE, callback)
}

func (usb *LinuxUSB) Read() ([]byte, error) {
	fmt.Println("Reading from LinuxUSB: ", usb.fname)
	return nil, nil
}

func (usb *LinuxUSB) Write(buff []byte) error {
	fmt.Println("Writing to LinuxUSB: ", usb.fname)
	return nil
}

func (usb *LinuxUSB) Close() error {
	fmt.Println("Closing LinuxUSB: ", usb.fname)
	return nil
}

func (usb *LinuxUSB) WaitForDisconnect() error {
	fmt.Println("WaitingForDisconnect LinuxUSB: ", usb.fname)
	return nil
}
