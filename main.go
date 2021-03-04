package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"time"

	// "github.com/periph/conn/spi/spireg"
	// "github.com/periph/devices/mfrc522"
	// "github.com/periph/host"
	// "github.com/periph/host/rpi"
	"github.com/periph/conn/spi/spireg"
	"github.com/periph/devices/mfrc522"
	"github.com/periph/host"
	"github.com/periph/host/rpi"
)

func main() {
	var breakMe bool
	var choice int

	for breakMe == false {
		fmt.Println("Choose your option:")
		fmt.Println("1 - Verify Password")
		fmt.Println("2 - Scan Card")
		fmt.Println("3 - Flash Card")
		fmt.Println("9 - Exit")
		//choice = 4
		switch fmt.Scan(&choice); choice {
		//switch choice {
		case 1:


		case 2:
			Search()
			// if err != nil {
			// 	log.Println("Error reading card")
			// } else {
			// 	log.Println("Card ID : ", id)
			// }
		case 3:
			//Place holder for Enroll Function

		case 9:
			breakMe = true
			fmt.Println("Stoping the program - with Exit Option")

		default:
			fmt.Println("Thats Invalid choice")
		}

	}

}

func Search() {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Using SPI as an example. See package "periph.io/x/conn/v3/spi/spireg" for more details.
	p, err := spireg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	rfid, err := mfrc522.NewSPI(p, rpi.P1_22, rpi.P1_18)
	if err != nil {
		log.Fatal(err)
	}

	// Idling device on exit.
	defer rfid.Halt()

	// Setting the antenna signal strength.
	rfid.SetAntennaGain(5)

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
