package transporter

import (
	"fmt"
	"io/ioutil"
	"strings"
	"unicode"
)

const (
	MAX_RETRIES = 5
	WAIT_FOR_DISCONNECT_TIME = 3
	MAX_USBFS_BULK_SIZE = (16 * 1024)
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

func findUsbDevice(base string) (*UsbHandler, error) {
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
		fmt.Println(file.Name())
	}

	return new(UsbHandler), nil
}

func UsbOpen() {
	findUsbDevice("/sys/bus/usb/devices")
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
