package main

import (
	"fmt"
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

/* Linux USB implementation */
type LinuxUSB struct {
	UsbHandler
}

func NewLinuxUSB(handler UsbHandler) (*LinuxUSB){
	return &LinuxUSB{ handler }
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
