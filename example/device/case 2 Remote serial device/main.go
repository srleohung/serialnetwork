package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
)

var logger Logger = NewLogger("main")

var serialConfig serial.Config = serial.Config{
	Name:     "",
	Baud:     9600,
	Size:     8,
	Parity:   serial.ParityNone,
	StopBits: serial.Stop1,
}
var serverHost string = "http://localhost:9876"
var deviceHostPort string = ":9877"
var rxLength int = 1

func main() {
	// ***** Init serial device *****
	SerialDevice := serialnetwork.NewSerialDevice()
	/*
		You can call initialization from server api.
		If you want, you don't need to run initialize(SerialDevice.Init(serialConfig, rxLength)).
	*/
	SerialDevice.Init(serialConfig, rxLength)

	// ***** Start channel handler service *****
	/*
		If you use api calls to control and you don't need to automatic response,
		please don't run this function.
	*/
	SerialDevice.ResponseToServer(serverHost)
	SerialDevice.RequestFromServer(deviceHostPort)
}
