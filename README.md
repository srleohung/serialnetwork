# Serial Network
A Go package to allow you easily to read and write from the serial port uses channel or API call. This package can help you to easily make the serial network to remotely control your device. This package depends on "github.com/tarm/serial".

# Usage
## Device
#### rxChannel for read from the serial port
#### txChannel for write to the serial port
```
package main

import (
	"github.com/srleohung/serialnetwork"
	"log"
)

func main() {
	d := serialnetwork.NewDevice()
	d.Init(serialnetwork.Config{Name: "/dev/ttyUSB0", Baud: 115200, RxBuffer: 1})

	rx := d.GetRxChannel()
	tx := d.GetTxChannel()

	tx <- []byte("test")
	log.Printf("% x", <-rx)
}
```
#### read the serial port by format
```
var rxFormat []serialnetwork.RxFormat = []serialnetwork.RxFormat{
	// Example format 1
	{
		StartByte: []byte{0x01},
		EndByte:   []byte{0x09},
	},
	// Example format 2
	{
		StartByte:         []byte{0x011},
		LengthByteIndex:   1,
		LengthByteMissing: 7,
	},
	// Example format 3
	{
		StartByte:   []byte{0x21},
		LengthFixed: 9,
	},
}
d.SetRxFormat(rxFormat)
```
#### Example link - https://github.com/srleohung/serialnetwork/blob/master/example/device/case%203%20Serial%20device%20read%20by%20format/main.go
## Server
#### Please read examples of serial device and server.
* Serial device examples - https://github.com/srleohung/serialnetwork/tree/master/example/device
* Serial server examples - https://github.com/srleohung/serialnetwork/tree/master/example/server

# Possible Future Work
* Return error message and device connection status to server.
* Use RWMutex to prevent error signals.
* Upgrade device automatically reconnects.
* The Init device from the server is not working properly.

# Startup Reason
This package is startup of a unmanned store projects. A host to controls many vending machine, camera, payment device. This package is for many different devices connected to one server.