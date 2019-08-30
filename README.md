# Serial Network
A Go package to allow you easily to read and write from the serial port uses channel or API call. This package can help you to easily make the serial network to remotely control your device. This package depends on "github.com/tarm/serial".

# Usage
## Device
rxChannel for read from the serial port
txChannel for write to the serial port
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
## Server
#### Please read examples of serial device and server.
* Serial device examples - https://github.com/srleohung/serialnetwork/example/device
* Serial server examples - https://github.com/srleohung/serialnetwork/example/server

# Possible Future Work
* Return error message and device connection status to server.
* Use RWMutex to prevent error signals.
* Upgrade device automatically reconnects.
* The Init device from the server is not working properly.

# Startup Reason
This package is startup of a unmanned store projects. A host to controls many vending machine, camera, payment device. This package is for many different devices connected to one server.