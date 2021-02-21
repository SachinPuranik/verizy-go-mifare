## Verizy - golang library for MFRC522 communication  

This is a library used for communicating with mifare card MFRC522.


## General terms
* SS - Slave select(Master Out - active low)
* CS - Chip Select (Slave In)
* MOSI - Master Out Slave In(Common)
* MISO - Master In Slave Out(Common)
* SCK - Serial clock (Common)

* CPOL - Clock polarity
* CPHA - Clock phase
* CPOL and CPHA are used to identify the clock edge to begin
* PICC - short for Proximity Integrated Circuit Card (RFID Tag itself) 
* PCD - means Proximity Coupling Device
## Referances

[SPI tutorial](https://www.corelis.com/education/tutorials/spi-tutorial/)
[raspberry Connection](https://www.raspberrypi-spy.co.uk/2018/02/rc522-rfid-tag-read-raspberry-pi/)
[RPI j8](https://radiostud.io/understanding-spi-in-raspberry-pi/)
[RPI GPIO Layout](https://www.bigmessowires.com/2018/05/26/raspberry-pi-gpio-programming-in-c/)
[GPIO Layout](https://microcontrollerslab.com/raspberry-pi-4-pinout-description-features-peripherals-applications/)
[RC522 with Audrino](https://lastminuteengineers.com/how-rfid-works-rc522-arduino-tutorial/)