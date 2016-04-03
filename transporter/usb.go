package transporter

type UsbIfcInfo struct {
	dev_vendor 	uint16
	dev_product 	uint16

	dev_class	uint8
	dev_subclass	uint8
	dev_protocol	uint8

	ifc_class	uint8
	ifc_subclass	uint8
	ifc_protocol	uint8

	has_bulk_in	bool
	hsa_bulk_out	bool

	writable	bool

	serial_number	string
	device_path	string
}

// TODO: add description
type ifc_match_func func(usb_ifc_info UsbIfcInfo) (int)

