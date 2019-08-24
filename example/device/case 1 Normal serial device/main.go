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
var rxLength int = 1

var message []byte = []byte("test")

var rxChannel chan []byte
var txChannel chan []byte
var txWroteChannel chan []byte

func main() {
	// ***** Init serial device *****
	SerialDevice := serialnetwork.NewSerialDevice()
	SerialDevice.Init(serialConfig, rxLength)

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
		logger.Infof("rxChannel % x", <-rxChannel)
	}()
	txChannel <- message
	logger.Infof("txWroteChannel % x", <-txWroteChannel)
}
