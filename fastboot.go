package main

import  (
	"fmt"
	"github.com/malicote/go_fastboot/transporter"
)

func main() {
	fmt.Println("go_fastboot")

	transporter.UsbOpen()
}