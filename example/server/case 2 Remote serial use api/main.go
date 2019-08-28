package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
)

var logger Logger = NewLogger("main")

var serialDeviceConfig serialnetwork.SerialDeviceConfig = serialnetwork.SerialDeviceConfig{
	Name: "/dev/ttyUSB0",
	Baud: 115200,
	// ReadTimeout: 1000,
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
	RxLength:   1,
	ServerHost: "http://localhost:9876",
}

const DeviceHost string = "http://localhost:9877"

var message []byte = []byte("test")

func main() {
	// ***** Init service *****
	SerialServer := serialnetwork.NewSerialServer()
	SerialServer.Init()

	// ***** Init serial device *****
	/*
		You can call initialization from server api.
		If you don't want, you don't need to run initialize(SerialServer.InitSerialDevice(serialDeviceConfig)).
	*/
	SerialServer.InitSerialDevice(serialDeviceConfig)

	// ***** Test service *****
	SerialServer.SetDeviceHost(DeviceHost)
	logger.Infof("TxRequest % x", SerialServer.TxRequest(message))
	logger.Infof("TxRequestAndRxResponse % x", SerialServer.TxRequestAndRxResponse(message))
}
