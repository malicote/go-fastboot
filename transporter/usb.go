package transporter

type UsbIfcInfo struct {
	DevVendor 	uint16
	DevProduct 	uint16

	DevClass	uint8
	DevSubclass	uint8
	DevProtocol	uint8

	IfcClass	uint8
	IfcSubclass	uint8
	IfcProtocol	uint8

	HasBulkIn	bool
	HasBulkOut	bool

	Writeable	bool

	SerialNumber	string
	DevicePath	string
}

// TODO: add description
/* Return true if matches, false otherwise */
type ifc_match_func func(usb_ifc_info UsbIfcInfo) bool

