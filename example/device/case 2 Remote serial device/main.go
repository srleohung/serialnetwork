package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
)

var logger Logger = NewLogger("main")

var serialDeviceConfig serialnetwork.SerialDeviceConfig = serialnetwork.SerialDeviceConfig{
	Name: "/dev/ttyUSB0",
	Baud: 115200,
	// ReadTimeout: 1000 * time.Millisecond,
	Size:   8,
	Parity: serialnetwork.ParityNone,
	/*
		ParityNone  Parity = 'N'
		ParityOdd   Parity = 'O'
		ParityEven  Parity = 'E'
		ParityMark  Parity = 'M' // parity bit is always 1
		ParitySpace Parity = 'S' // parity bit is always 0
	*/
	StopBits: serialnetwork.Stop1,
	/*
		Stop1     StopBits = 1
		Stop1Half StopBits = 15
		Stop2     StopBits = 2
	*/
	RxBuffer: 1,
	// ServerHost: "",
}

const ServerHost string = "http://localhost:9876"
const DeviceAddr string = ":9877"

func main() {
	// ***** Init serial device *****
	SerialDevice := serialnetwork.NewSerialDevice()
	/*
		You can call initialization from server api.
		If you want, you don't need to run initialize(SerialDevice.Init(serialConfig, rxBuffer)).
	*/
	err := SerialDevice.Init(serialDeviceConfig)
	if err != nil {
		logger.Emerg(err)
	}

	// ***** Start channel handler service *****
	/*
		If you use api calls to control and you don't need to automatic response,
		please don't run this function.
	*/
	SerialDevice.ResponseToServer(ServerHost)
	SerialDevice.RequestFromServer(DeviceAddr)

	// ***** Run forever *****
	forever := make(chan bool)
	<-forever
}
