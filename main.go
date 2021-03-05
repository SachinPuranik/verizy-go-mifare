package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/periph/devices/mfrc522/commands"
	"periph.io/x/conn/spi/spireg"
	"periph.io/x/devices/mfrc522"
	"periph.io/x/host"
	"periph.io/x/host/rpi"
)

var (
	rfid *mfrc522.Dev
)

func main() {
	var breakMe bool
	var choice int
	var err error

	err = nil

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Using SPI as an example. See package "periph.io/x/conn/v3/spi/spireg" for more details.
	p, errOpen := spireg.Open("/dev/spidev0.0")
	if errOpen != nil {
		log.Fatal(err)
	}
	defer p.Close()

	rfid, err = mfrc522.NewSPI(p, rpi.P1_22, rpi.P1_18)
	if err != nil {
		log.Fatal(err)
	}

	// Idling device on exit.
	defer rfid.Halt()

	// Setting the antenna signal strength.
	rfid.SetAntennaGain(5)

	for breakMe == false {
		fmt.Println("Choose your option:")
		fmt.Println("1 - Scan Card")
		fmt.Println("2 - Write test Card")
		fmt.Println("9 - Exit")
		//choice = 4
		switch fmt.Scan(&choice); choice {
		//switch choice {
		case 1:
			search()
		case 2:
			write()
		case 3:
		case 9:
			breakMe = true
			fmt.Println("Stoping the program - with Exit Option")

		default:
			fmt.Println("Thats Invalid choice")
		}

	}

}

func search() {
	timedOut := false
	cb := make(chan []byte)
	timer := time.NewTimer(10 * time.Second)

	// Stopping timer, flagging reader thread as timed out
	defer func() {
		timer.Stop()
		timedOut = true
		close(cb)
	}()

	go func() {
		log.Printf("Started %s", rfid.String())

		for {
			// Trying to read card UID.
			uid, err := rfid.ReadUID(10 * time.Second)

			// If main thread timed out just exiting.
			if timedOut {
				return
			}

			// Some devices tend to send wrong data while RFID chip is already detected
			// but still "too far" from a receiver.
			// Especially some cheap CN clones which you can find on GearBest, AliExpress, etc.
			// This will suppress such errors.
			if err != nil {
				continue
			}

			cb <- uid
			return
		}
	}()

	for {
		select {
		case <-timer.C:
			log.Fatal("Didn't receive device data")
			return
		case data := <-cb:
			log.Println("UID:", hex.EncodeToString(data))
			return
		}
	}
}

func write() {
	// Converting access key.
	// This value corresponds to first pi "numbers": 3 14 15 92 65 35.
	hexKey, _ := hex.DecodeString("030e0f5c4123")
	var key [6]byte
	copy(key[:], hexKey)

	data, err := rfid.ReadAuth(10*time.Second, byte(commands.PICC_AUTHENT1B), 1, key)

	if err != nil {
		log.Println("getting kills as auth reading failed")
		log.Panicln(err.Error())
	} else {
		log.Println("Auth Key :", data)
	}

	// Converting expected data.
	// This value corresponds to string "@~>f=Um[X{LRwA3}".
	expected, _ := hex.DecodeString("407e3e663d556d5b587b4c527741337d")

	timedOut := false
	cb := make(chan []byte)
	timer := time.NewTimer(10 * time.Second)

	// Stopping timer, flagging reader thread as timed out
	defer func() {
		timer.Stop()
		timedOut = true
		close(cb)
	}()

	go func() {
		log.Printf("Started %s", rfid.String())

		for {
			// Trying to read data from sector 1 block 0
			data, err := rfid.ReadCard(10*time.Second, byte(commands.PICC_AUTHENT1B), 1, 0, key)

			// If main thread timed out just exiting.
			if timedOut {
				return
			}

			// Some devices tend to send wrong data while RFID chip is already detected
			// but still "too far" from a receiver.
			// Especially some cheap CN clones which you can find on GearBest, AliExpress, etc.
			// This will suppress such errors.
			if err != nil {
				continue
			}

			cb <- data
		}
	}()

	for {
		select {
		case <-timer.C:
			log.Fatal("Didn't receive device data")
			return
		case data := <-cb:
			if !reflect.DeepEqual(data, expected) {
				log.Fatal("Received data is incorrect")
			} else {
				log.Println("Received data is correct")
			}

			return
		}
	}
}
