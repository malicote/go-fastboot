package main

import  (
	"fmt"
	"github.com/malicote/go_fastboot/transporter"
)

const (
	EXPECTED_IFC_CLASS = 0xff
	EXPECTED_IFC_SUBCLASS = 0x42
	EXPECTED_IFC_PROTOCOL = 0x03

)

func match_fastboot_with_serial(info transporter.UsbIfcInfo, local_serial string) bool {
	/* TODO: check vendor_id is specified by user */
	/* TODO: check serial number if specified by user */
	return (info.IfcClass == EXPECTED_IFC_CLASS &&
		info.IfcSubclass == EXPECTED_IFC_SUBCLASS &&
		info.IfcProtocol == EXPECTED_IFC_PROTOCOL)
}

func debug_callback(usb_ifc_info transporter.UsbIfcInfo) bool {
	fmt.Printf("**********\nReceived: %+#v\n*********\n", usb_ifc_info)
	return false
}

func list_devices_callback(usb_ifc_info transporter.UsbIfcInfo) bool {
	if match_fastboot_with_serial(usb_ifc_info, "") {
		serial := usb_ifc_info.SerialNumber
		if serial == "" {
			serial = "?????????"
		}

		/* TODO: correct printing ! */
		fmt.Printf("%s\tfastboot\n", serial)
	}
	return false
}

func main() {
	fmt.Println("go_fastboot!")

	transporter.UsbOpen(list_devices_callback)
}