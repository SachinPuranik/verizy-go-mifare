package cardscanner

import "github.com/tarm/serial"

type mySerial struct {
	port *serial.Port
	cfg  *serial.Config
}

type Card struct {
	mySerial
}

//ScannerIO - Interface for Scanner
type CardReaderIO interface {
	Capture() error
	Release()
	VerifyPassword() bool
	Scan() (string, error)
	Flash(data []byte) error
}

//NewCardScanner
func NewCardScanner(node string, spd int) (CardReaderIO) {
	// - /dev/spidev0.0
	return &Card{}

}

func (c *Card) Capture() error {
	return nil
}

func (c *Card) Release() {

}

func (c *Card) VerifyPassword() bool {
	return true
}

func (c *Card) Scan() (string, error) {
	return "asasas", nil
}

func (c *Card) Flash(data []byte) error {
	return nil
}
