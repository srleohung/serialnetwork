# Serial Network
A Go package to allow you easily to read and write from the serial port uses channel or API call. This package can help you to easily make the serial network to remotely control your device.

# Usage
## Device
rxChannel for read from the serial port
txChannel for write from the serial port
```
package main

import (
	"github.com/srleohung/serialnetwork"
	"github.com/tarm/serial"
	"log"
)

func main() {
	SerialDevice := serialnetwork.NewSerialDevice()
	SerialDevice.Init(&serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}, 1)

	rxChannel = SerialDevice.GetRxChannel()
	txChannel = SerialDevice.GetTxChannel()

	txChannel <- []byte("test")
	log.Printf("% x", <-rxChannel)
}
```
## Server
#### Please read examples of serial device and server.
##### Serial device - https://github.com/srleohung/serialnetwork/example/device
##### Serial server - https://github.com/srleohung/serialnetwork/example/server