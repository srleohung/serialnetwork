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
	Device := serialnetwork.NewDevice()
	Device.Init(&serial.Config{Name: "/dev/ttyUSB0", Baud: 9600}, 1)

	rxChannel = Device.GetRxChannel()
	txChannel = Device.GetTxChannel()

	txChannel <- []byte("test")
	log.Printf("% x", <-rxChannel)
}
```
## Server
#### Please read examples of serial device and server.
* Serial device examples - https://github.com/srleohung/serialnetwork/example/device
* Serial server examples - https://github.com/srleohung/serialnetwork/example/server

# Possible Future Work
* Return error message and device connection status to server.
* Use RWMutex to prevent error signals.
* Upgrade device automatically reconnects.

# Startup Reason
This package is startup of a unmanned store projects. A host to controls many vending machine (e.g. Fuji, O2O, TCN, XY), camera, payment device (e.g. HK octopus, Alipay, WeChat Pay, QR code reader, rfid reader ...). This package is for many different devices connected to one server.