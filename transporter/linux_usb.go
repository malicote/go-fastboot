package transporter

// TODO: use proper logging (debug, fatal, etc)

import (
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

func findUsbDevice(base string) (*UsbHandler, error) {
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

	files, err := ioutil.ReadDir(base)
	if err != nil {
		fmt.Println("Error reading ", base, ":", err)
	}

	for _, file := range files {
		if !is_device(file.Name()) {
			continue
		}

		dev_path, err := convertToDevFSName(file.Name())
		if err != nil {
			// TODO: verbose printing only
			fmt.Println(err)
			continue
		}

		writeable := true
		if _, err := os.OpenFile(dev_path, os.O_RDWR, 0); err != nil {
			// Check if device node is read-only, then we can provide
			// a helpful error message similar to 'adb devices'
			// TODO: verbose
			fmt.Println(err)
			writeable = false
			if _, err = os.OpenFile(dev_path, os.O_RDONLY, 0); err != nil {
				// TODO: verbose
				fmt.Println("Skipping: ", dev_path, "[", err, "]")
				continue
			}
		}


		fmt.Printf("Devfs path %s, writeable: %v\n", dev_path, writeable)

	}

	return new(UsbHandler), nil
}

// TODO: move this to its own package or figure out way to conditionally compile
func UsbOpen() {
	findUsbDevice(USB_SYSFS_BASE)
}

func (usb *LinuxUSB) Read() ([]byte, error) {
	fmt.Println("Reading from LinuxUSB: ", usb.fname)
	return nil, nil
}

func (usb *LinuxUSB) Write(bytes []byte) error {
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
