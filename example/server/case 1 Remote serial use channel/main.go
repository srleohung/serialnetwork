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

const ServerAddr string = ":9876"
const DeviceHost string = "http://localhost:9877"

var rxChannel chan []byte
var txChannel chan []byte

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

	// ***** Get channel *****
	if rxChannel = SerialServer.GetRxChannel(); rxChannel != nil {
		logger.Info("Got RxChannel")
	}
	if txChannel = SerialServer.GetTxChannel(); txChannel != nil {
		logger.Info("Got TxChannel")
	}

	// ***** Start channel handler service *****
	SerialServer.ResponseFromDevice(ServerAddr)
	SerialServer.RequestToDevice(DeviceHost)

	// ***** Test channel *****
	txChannel <- message
	for {
		logger.Infof("rxChannel % x", <-rxChannel)
	}
}
