package transporter

/* From ch9.h */
const (
	USB_DT_DEVICE 		= 0x1
	USB_DT_CONFIG		= 0x2
	USB_DT_STRING		= 0x3
	USB_DT_INTERFACE 	= 0x4
	USB_DT_ENDPOINT		= 0x5

	USB_DT_DEVICE_SIZE	= 18
	USB_DT_CONFIG_SIZE	= 9
	USB_DT_INTERFACE_SIZE	= 9
	USB_DT_ENDPOINT_SIZE	= 7

	USB_ENDPOINT_XFERTYPE_MASK	= 0x3
	USB_ENDPOINT_XFER_BULK		= 2
	USB_ENDPOINT_DIR_MASK		= 0x80

	USB_DT_SS_ENDPOINT_COMP		= 0x30
	USB_DT_SS_EP_COMP_SIZE		= 6
)

/* Little endian */

/* Due to golang struct parsing w/ encoding/bytes, the first letter
   must be capitalized.
 */

type usb_descriptor_header struct {
	BLength		uint8
	BDescriptorType	uint8
}

type usb_device_descriptor struct {
	BLength			uint8
	BDescriptorType		uint8

	BcdUSB			uint16
	BDeviceClass		uint8
	BDeviceSubClass		uint8
	BDeviceProtocol		uint8
	BMaxPacketSize0		uint8
	IdVendor		uint16
	IdProduct		uint16
	BcdDevice		uint16
	IManufacturer		uint8
	IProduct		uint8
	ISerialNumber		uint8
	BNumConfigurations 	uint8
}

type usb_config_descriptor struct {
	BLength			uint8
	BDescriptorType		uint8

	WTotalLength		uint16
	BNumInterfaces		uint8
	BConfigurationValue	uint8
	IConfiguration		uint8
	BmAttributes		uint8
	BMaxPower		uint8
}

type usb_interface_descriptor struct {
	BLength			uint8
	BDescriptorType		uint8

	BInterfaceNumber	uint8
	BAlternateSetting	uint8
	BNumEndpoints		uint8
	BInterfaceClass		uint8
	BInterfaceSubClass	uint8
	BInterfaceProtocol	uint8
	IInterface		uint8
}

type usb_endpoint_descriptor struct {
	BLength			uint8
	BDescriptorType		uint8

	BEndpointAddress	uint8
	BmAttributes		uint8
	WMaxPacketSize		uint16

	BRefresh		uint8
	BSynchAddress		uint8
}

