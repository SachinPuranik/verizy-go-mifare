package cardscanner

import (
	//"github.com/warthog618/gpiod"
	//"github.com/ecc1/spi"

	"errors"
	"log"
	"strconv"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

//Card -
type Card struct {
	// spiDevice     *spi.Device
	speed         int
	spiDeviceAddr string
	// chip          *gpiod.Chip
}

//CardReaderIO - Interface for Scanner
type CardReaderIO interface {
	Capture() error
	Release()
	VerifyPassword() bool
	Scan() (string, error)
	Flash(data []byte) error
	RequestMode(int) (tagType int, err error)
	ReadWithAnticoll() (uid string, err error)
}

//NewCardScanner - Params - spiDeviceAddr > /dev/spidev0.0, speed > 100000
func NewCardScanner(spiDeviceAddr string, speed int) CardReaderIO {
	return &Card{spiDeviceAddr: spiDeviceAddr, speed: speed}

}

//Capture -
func (c *Card) Capture() error {

	if err := rpio.Open(); err != nil {
		panic(err)
	}

	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		panic(err)
	}

	rpio.SpiChipSelect(0) // Select CE0 slave

	c.init()
	return nil
}

func (c *Card) init() {

	pin := rpio.Pin(NRSTPD)
	pin.Mode(rpio.Output) // Alternative syntax
	pin.Write(rpio.High)  // Alternative syntax

	// GPIO.setmode(GPIO.BOARD)
	// GPIO.setup(self.NRSTPD, GPIO.OUT)
	// GPIO.output(self.NRSTPD, 1)

	c.DeviceReset()

	c.writeToDevice(TModeReg, 0x8D)
	c.writeToDevice(TPrescalerReg, 0x3E)
	c.writeToDevice(TReloadRegL, 30)
	c.writeToDevice(TReloadRegH, 0)

	c.writeToDevice(TxAutoReg, 0x40)
	c.writeToDevice(ModeReg, 0x3D)
	c.AntennaOn()
}

//DeviceReset -
func (c *Card) DeviceReset() {
	c.writeToDevice(CommandReg, PCD_RESETPHASE)
}

//Release -
func (c *Card) Release() {
	rpio.SpiEnd(rpio.Spi0)
	rpio.Close()
}

//VerifyPassword -
func (c *Card) VerifyPassword() bool {
	return true
}

//Scan -
func (c *Card) Scan() (string, error) {
	return "asasas", nil
}

//Flash -
func (c *Card) Flash(data []byte) error {
	return nil
}

func (c *Card) writeToDevice(addr int, val int) ([]byte, error) {
	var outBuf []byte
	var responseBytes []byte
	log.Println("Addr:(", addr, ")   val:(", val, ")")

	// outBuf = make([]byte, 5)

	bAddr := []byte(strconv.Itoa(((addr << 1) & 0x7E)))
	bVal := []byte(strconv.Itoa(val))
	outBuf = append(outBuf, bAddr[:]...)
	outBuf = append(outBuf, bVal[:]...)
	rpio.SpiExchange(outBuf)
	responseBytes = outBuf
	return responseBytes, nil
}

func (c *Card) setBitMask(reg int, mask int) {
	tmp, _ := c.readFromDevice(reg)
	c.writeToDevice(reg, int(tmp)|mask)
}

func (c *Card) clearBitMask(reg int, mask int) {
	tmp, _ := c.readFromDevice(reg)
	c.writeToDevice(reg, int(tmp)&(^mask))
}

// func (c *Card) readFromDevice(addr int) (byte, error) {
// 	log.Println("Addr:(", addr, ")   val:(", 0, ")")
// 	responseBytes, err := c.writeToDevice(addr, 0)
// 	if err != nil {
// 		log.Println("readFromDevice :", err.Error())
// 	}
// 	return responseBytes[1], err
// }

func (c *Card) readFromDevice(addr int) (byte, error) {
	var outBuf []byte
	var responseBytes []byte

	log.Println("Addr:(", addr, ")   val:(0)")

	responseBytes = nil

	bAddr := []byte(strconv.Itoa(((addr << 1) & 0x7E)))
	outBuf = append(outBuf, bAddr[:]...)
	//err := c.spiDevice.Read(outBuf) //Transfer(outBuf)
	rpio.SpiExchange(outBuf)
	responseBytes = outBuf
	return responseBytes[1], nil
}

//AntennaOn -
func (c *Card) AntennaOn() {
	temp, _ := c.readFromDevice(TxControlReg)
	if (temp & 0x03) == 0x00 {
		c.setBitMask(TxControlReg, 0x03)
	}
}

//AntennaOff -
func (c *Card) AntennaOff() {
	c.clearBitMask(TxControlReg, 0x03)
}

func (c *Card) writeCommandToCard(command int, sendData []byte) (responseData []byte, responseLen int, err error) {

	irqEn := 0x00
	waitIRq := 0x00
	status := MI_ERR

	responseLen = 0
	responseData = make([]byte, 0)

	// lastBits := None

	if command == PCD_AUTHENT {
		irqEn = 0x12
		waitIRq = 0x10
	} else if command == PCD_TRANSCEIVE {
		irqEn = 0x77
		waitIRq = 0x30
	}

	c.writeToDevice(CommIEnReg, irqEn|0x80)
	c.clearBitMask(CommIrqReg, 0x80)
	c.setBitMask(FIFOLevelReg, 0x80)

	c.writeToDevice(CommandReg, PCD_IDLE)

	i := 0
	log.Println("sendData Len :", len(sendData))

	for i < len(sendData) {
		log.Println("Sending byte :", sendData[i])
		c.writeToDevice(FIFODataReg, int(sendData[i]))
		i = i + 1
	}

	c.writeToDevice(CommandReg, command)

	if command == PCD_TRANSCEIVE {
		c.setBitMask(BitFramingReg, 0x80)
	}

	i = 2000
	n := byte(0)
	for {
		n, _ = c.readFromDevice(CommIrqReg)
		i = i - 1
		log.Println("read Byte : ", n)
		//if ~((i!=0) and ~(n&0x01) and ~(&waitIRqn)) (~(false)
		//C1 - Break if i is 0 - evaulate expr false - hence false ==false break
		//C2 - Break when last bit of n is 1 - evaulate expr false - hence false ==false break
		//C3 - Break when 5th or 6th bit are 1 - evaulate expr false - hence false ==false break
		if ((i != 0) && (int(n)&0x01) <= 0 && (int(n)&waitIRq) <= 0) == false {
			log.Println("WIll break loop")
			break
		}
	}

	c.clearBitMask(BitFramingReg, 0x80)

	if i != 0 {
		b, _ := c.readFromDevice(ErrorReg)
		if (b & 0x1B) == 0x00 {
			status = MI_OK

			if int(n) != 0&irqEn&0x01 {
				log.Println("MI_NOTAGERR")
				status = MI_NOTAGERR
			}

			if command == PCD_TRANSCEIVE {
				n, _ = c.readFromDevice(FIFOLevelReg)
				lastBits, _ := c.readFromDevice(ControlReg)
				lastBitsInteger := int(lastBits) & 0x07
				if lastBits != 0 {
					responseLen = (int(n)-1)*8 + lastBitsInteger
				} else {
					responseLen = int(n) * 8
				}
			}

			if n == 0 {
				n = 1
			}
			if n > MAX_LEN {
				n = MAX_LEN
			}

			i = 0
			for i < int(n) {
				resp, _ := c.readFromDevice(FIFODataReg)
				responseData = append(responseData, resp)
				i = i + 1
				log.Println("looping")
			}
		}
	} else {

		status = MI_ERR
		log.Println("MI_ERR and i is 0")
	}

	if status != MI_OK {
		err = errors.New("Some error")
		responseData = nil
		responseLen = -1
	}

	return responseData, responseLen, err
}

//RequestMode -
func (c *Card) RequestMode(reqMode int) (int, error) {

	err := error(nil)
	outBuf := make([]byte, 0)

	c.writeToDevice(BitFramingReg, 0x07)
	outBuf = append(outBuf, byte(reqMode))

	responseBytes, responseLen, err := c.writeCommandToCard(PCD_TRANSCEIVE, outBuf)

	if (err != nil) || (responseLen != 0x10) {
		return 0, err
	}
	log.Println("Response ->", responseBytes)

	return responseLen, err
}

//ReadWithAnticoll -
func (c *Card) ReadWithAnticoll() (uid string, err error) {

	return "nil", nil
}
