package main

type Transporter interface {
	/* Returns slice of bytes read, or error */
	Read() ([]byte, error)

	/* Provide bytes to write, returns error */
	Write([]byte) error

	/* Close connection, returns error string */
	Close() error

	/* Wait for disconnect, returns error string */
	WaitForDisconnect() error
}


