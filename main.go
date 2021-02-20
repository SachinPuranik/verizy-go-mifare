package main

import (
	"fmt"
	"log"

	//"github.com/SachinPuranik/verizy-go-mifare/cardscanner"
	"github.com/SachinPuranik/verizy-go-mifare/cardscanner"
)

func main() {
	var breakMe bool
	var choice int

	scanner := cardscanner.NewCardScanner("/dev/spidev0.0", 0x0000)

	err := scanner.Capture()
	if err != nil {
		log.Fatal("Wow...Cant't handel err =>", err.Error())
	}

	defer scanner.Release()

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
			if scanner.VerifyPassword() == true {
				log.Println("Password verified")
			} else {
				log.Println("Password wrong")
			}

		case 2:
			id, err := Search(scanner)
			if err != nil {
				log.Println("Error reading card")
			} else {
				log.Println("Card ID : ", id)
			}
		case 3:
			//Place holder for Enroll Function
			writeBuf := []byte("This is test")
			err := scanner.Flash(writeBuf)
			if err != nil {
				log.Println("Error writing card")
			} else {
				log.Println("Wrote :", string(writeBuf))
			}

		case 9:
			breakMe = true
			fmt.Println("Stoping the program - with Exit Option")

		default:
			fmt.Println("Thats Invalid choice")
		}

	}

}

//Search -
func Search(scanner cardscanner.CardReaderIO) (string, error) {

	var continueRead bool
	continueRead = true

	for continueRead == true {
		_, err := scanner.RequestMode(cardscanner.PICC_REQIDL)

		if err == nil {
			log.Println("Card detected")
			continueRead = false
		} else {
			log.Println("Card RequestMode error :", err.Error())
		}
		continueRead = false
		// // Get the UID of the card with anti collision
		// uid = ""
		// uid , err := cardscanner.ReadWithAnticoll()
		// if(err == nil){
		// 	scanner.SelectTag(ui)
		//	err := scanner.AuthanticateTag(ui)
		// if err == nil{
		//     scanner.Read(8)
		//     scanner.StopCrypto1()
		// } else{}
		//     log.Println("Authentication error")
		// }
	}

	return "Nothing", nil
}
