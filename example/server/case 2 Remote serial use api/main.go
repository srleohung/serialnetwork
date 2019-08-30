package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
	"time"
)

var logger Logger = NewLogger("main")

var config serialnetwork.Config = serialnetwork.Config{
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
	RxBuffer:   1,
	ServerHost: "http://localhost:9876",
}

const DeviceHost string = "http://localhost:9877"

var message []byte = []byte("test")

func main() {
	// ***** Init service *****
	s := serialnetwork.NewServer()
	err := s.Init()
	if err != nil {
		logger.Emerg(err)
	}
	s.SetDeviceHost(DeviceHost)

	// ***** Test Connection *****
	for {
		if s.Ping() {
			logger.Warning("The server is connecting to the device.")
			break
		} else {
			logger.Warning("The server cannot connect to the device.")
			time.Sleep(1 * time.Second)
		}
	}

	// ***** Init serial device *****
	/*
		You can call initialization from server api.
		If you don't want, you don't need to run initialize(s.InitDevice(config)).
	*/
	// err = s.InitDevice(config)
	// if err != nil {
	// 	logger.Emerg(err)
	// }

	// ***** Test service *****
	logger.Infof("TxRequest % x", s.TxRequest(message))
	logger.Infof("TxRequestAndRxResponse % x", s.TxRequestAndRxResponse(message))
}
