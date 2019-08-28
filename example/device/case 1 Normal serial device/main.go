package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
	"time"
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

var message []byte = []byte("1")

var rxChannel chan []byte
var txChannel chan []byte
var txWroteChannel chan bool

func main() {
	// ***** Init serial device *****
	SerialDevice := serialnetwork.NewSerialDevice()
	err := SerialDevice.Init(serialDeviceConfig)
	if err != nil {
		logger.Emerg(err)
	}

	// ***** Get channel *****
	if rxChannel = SerialDevice.GetRxChannel(); rxChannel != nil {
		logger.Info("Got RxChannel")
	}
	if txChannel = SerialDevice.GetTxChannel(); txChannel != nil {
		logger.Info("Got TxChannel")
	}
	/*
		If you want to check your tx message have been sent,
		you can use Tx wrote channel and listen to it.
	*/
	if txWroteChannel = SerialDevice.GetTxWroteChannel(); txWroteChannel != nil {
		logger.Info("Got TxWroteChannel")
	}

	// ***** Test channel *****
	go func() {
		for {
			txChannel <- message
			if <-txWroteChannel {
				logger.Info("The message was successfully written to the serial port.")
			} else {
				logger.Warning("The message failed to write to the serial port.")
			}
			time.Sleep(1 * time.Second)
		}
	}()
	logger.Infof("rxChannel % x", <-rxChannel)
}
