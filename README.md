# Serial Network
A Go package to allow you easily to read and write from the serial port uses channel or API call. This package can help you to easily make the serial network to remotely control your device. This package depends on github.com/tarm/serial and github.com/gorilla/websocket.

# Usage

## Device

#### RxChannel for read from the serial port. TxChannel for write to the serial port.
```go
package main

import (
	"github.com/srleohung/serialnetwork"
	"log"
)

func main() {
	d := serialnetwork.NewDevice()

    // Initialize and establish a serial port connection
	d.Init(serialnetwork.Config{Name: "/dev/ttyUSB0", Baud: 115200, RxBuffer: 1})

    // Get serial read and write channel
	rx := d.GetRxChannel()
	tx := d.GetTxChannel()

    // Use channels to read and write
	tx <- []byte("test")
	log.Printf("% x", <-rx)
}
```

#### Read the serial port by format
```go
// Set the serial read format
d.SetRxFormat([]serialnetwork.RxFormat{
	{StartByte: []byte{0x01}, EndByte: []byte{0x09}},
	{StartByte: []byte{0x01}, LengthByteIndex: 1, LengthByteMissing: 7},
	{StartByte: []byte{0x01}, LengthFixed: 9},
})

// Test sample read format
message := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
bytes := d.TestRxFormats(message)
if len(bytes) != len(message) {
	log.Print("Unknown message in Rx format.")
}

// Test serial read format
message = <-d.GetRxChannel()
bytes = d.TestRxFormats(message)
if len(bytes) != len(message) {
	log.Print("Unknown message in Rx format.")
}
```
* Example
  * https://github.com/srleohung/serialnetwork/blob/master/example/device/case%203%20Serial%20device%20read%20by%20format/main.go

#### Formater option
* StartByte 
  * Type: []byte
* EndByte
  * Type: []byte
* LengthByteIndex
  * Type: int
* LengthByteMissing
  * Type: int
* LengthFixed
  * Type: int
* LengthHighByteIndex
  * Type: int
* LengthDecoder
  * Type: string
  * Option: bcd

#### Encoding
* Block check character (BCC) 
  * import "github.com/srleohung/serialnetwork/encoding/bcc"
    * Encode
* Binary coded decimal (BCD) 
  * import "github.com/srleohung/serialnetwork/encoding/bcd"
    * Encode
    * Decode

## Server

#### Establish a network socket connection server for serial server
```go
s.NewWebSocketServer(":9876")
```

#### Establish a network socket connection client for serial device
```go
d.NewWebSocketClient("localhost:9876")
```

#### Read examples of serial device and server for more information
* Serial device example
  * https://github.com/srleohung/serialnetwork/blob/master/example/device/case%204%20Remote%20serial%20device%20use%20websocket/main.go
* Serial server example
  * https://github.com/srleohung/serialnetwork/blob/master/example/server/case%204%20Remote%20serial%20device%20use%20websocket/main.go

# Possible Future Work
* Return error messages and device connection status to the server.
* Use RWMutex to prevent error signals.
* Upgrade device automatically reconnects.
* The Init device on the server is not working properly.
* Add other encodings for use.
* Upgrade websocket to automatically reconnect.

# Startup Reason
This package is startup of a unmanned store projects. A host to controls many vending machine, camera, payment device. This package is for many different devices connected to one server.