package transporter

// TODO: use proper logging (debug, fatal, etc)

import (
	"bytes"
	"encoding/binary"
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
	n_read, err := buff.ReadAt(descriptor_buff, 0)
	if err != nil || n_read != 2 {
		return false
	}

	/* Per usb_descriptor_header, the type is byte offset 1 */
	return descriptor_buff[1] == usb_descriptor_type
}

func filter_usb_device(sysfs_name string,
			dev_data []byte,
			writable bool,
			callback ifc_match_func) (ept_in_id int, ept_out_id int, ifc_id int, err error) {

	var device_descriptor usb_device_descriptor

	buff := bytes.NewReader(dev_data)
	if descriptor_is_type(buff, USB_DT_DEVICE) {
		binary.Read(buff, binary.LittleEndian, &device_descriptor)
		/* TODO: parse everything then call back */
		callback(UsbIfcInfo{})
		return 0, 0, 0, nil
	}
	return 0, 0, 0, nil
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
