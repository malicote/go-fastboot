package main

import  (
	"fmt"
	"encoding/binary"
	"bytes"
	"github.com/malicote/go_fastboot/transporter"
)

type my_struct struct {
	M_8 	uint8
	M_8b	uint8
	M_16 	uint16
}

func descriptor_is_type(buff *bytes.Reader, usb_descriptor_type uint8) bool {
	descriptor_buff := []uint8{0, 0}
	n_read, err := buff.ReadAt(descriptor_buff, 0)
	if err != nil || n_read != 2 {
		return false
	}

	return descriptor_buff[1] == usb_descriptor_type
}

const (
	MS1 = 1
	MS2 = 2
)

func callback(usb_ifc_info transporter.UsbIfcInfo) bool {
	fmt.Println("My callback.")
	return false
}

func main() {
	fmt.Println("go_fastboot")

	//b := []byte{0x1, 0x1, 0x0B, 0x0, 0x3, 0x4, 0xC, 0x0}
	b := []byte{0x1}
	buff := bytes.NewReader(b)

	if descriptor_is_type(buff, MS1) {
		var ms1 my_struct
		err := binary.Read(buff, binary.LittleEndian, &ms1)
		fmt.Println("ms1: ", err, " | ", ms1)
	}
	if descriptor_is_type(buff, MS2) {
		var ms2 my_struct
		err2 := binary.Read(buff, binary.LittleEndian, &ms2)
		fmt.Println("ms2: ", err2, " | ", ms2)
	}
	transporter.UsbOpen(callback)
}