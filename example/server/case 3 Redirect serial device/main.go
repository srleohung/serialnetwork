package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
	"time"
)

var logger Logger = NewLogger("main")

var config1 serialnetwork.Config = serialnetwork.Config{
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
	// RxBuffer:   1,
	ServerHost: "http://localhost:9876",
}

var config2 serialnetwork.Config = serialnetwork.Config{
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
	// RxBuffer:   1,
	// ServerHost: "http://localhost:9876",
}

var rxFormat []serialnetwork.RxFormat = []serialnetwork.RxFormat{
	// Format 1
	{
		StartByte: []byte{0x01},
		EndByte:   []byte{0x09},
		// LengthByteIndex:   1,
		// LengthByteMissing: 7,
		// LengthFixed:       9,
	},
	// Format 2
	{
		StartByte:         []byte{0x11},
		LengthByteIndex:   1,
		LengthByteMissing: 7,
		// EndByte:           []byte{0x09},
		// LengthFixed:       9,
	},
	// Format 3
	{
		StartByte:   []byte{0x21},
		LengthFixed: 9,
		// EndByte:           []byte{0x09},
		// LengthByteIndex:   1,
		// LengthByteMissing: 7,
	},
}

const ServerAddr string = ":9876"
const DeviceHost string = "http://localhost:9877"

// First serial device
var rx1 chan []byte
var tx1 chan []byte

// Second serial device
var rx2 chan []byte
var tx2 chan []byte

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
			logger.Info("The server is connecting to the device.")
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
	// err = s.InitDevice(config1)
	// if err != nil {
	// 	logger.Emerg(err)
	// }

	// ***** Get first serial device channel *****
	if rx1 = s.GetRxChannel(); rx1 != nil {
		logger.Info("Got RxChannel")
	}
	if tx1 = s.GetTxChannel(); tx1 != nil {
		logger.Info("Got TxChannel")
	}

	// ***** Start channel handler service *****
	s.ResponseFromDevice(ServerAddr)
	s.RequestToDevice(DeviceHost)

	// ***** Init second serial device *****
	d := serialnetwork.NewDevice()
	d.SetRxFormat(rxFormat)
	err = d.Init(config2)
	if err != nil {
		logger.Emerg(err)
	}

	// ***** Get second serial device channel *****
	if rx2 = d.GetRxChannel(); rx2 != nil {
		logger.Info("Got RxChannel")
	}
	if tx2 = d.GetTxChannel(); tx2 != nil {
		logger.Info("Got TxChannel")
	}

	// ***** Start to redirect serial device *****
	go func(from chan []byte, to chan []byte) {
		for {
			bytes := <-from
			logger.Infof("rx1->tx2 % x", bytes)
			to <- bytes
		}
	}(rx1, tx2)
	for {
		bytes := <-rx2
		logger.Infof("rx2->tx1 % x", bytes)
		tx1 <- bytes
	}
}
