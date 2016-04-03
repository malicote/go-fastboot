package transporter

import (
	"errors"
)

/* OSX USB implementation */
type OsxUSB struct {
	UsbHandler
}

func (usb *OsxUSB) Read() ([]byte, error) {
	return nil, errors.New("Not supported.")
}

func (usb *OsxUSB) Write(bytes []byte) error {
	return errors.New("Not supported.")
}

func (usb *OsxUSB) Close() error {
	return errors.New("Not supported.")
}

func (usb *OsxUSB) WaitForDisconnect() error {
	return errors.New("Not supported.")
}
