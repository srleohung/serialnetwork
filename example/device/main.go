package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
)

var logger Logger = NewLogger("main")

var serialConfig serial.Config = serial.Config{
	Name:     "/dev/tty.SLAB_USBtoUART",
	Baud:     115200,
	Size:     8,
	Parity:   serial.ParityNone,
	StopBits: serial.Stop1,
}

var serverHost string = "http://localhost:9876"

var deviceHostPort string = ":9877"

var rxLength int = 1

var message []byte = []byte("9876")

func main() {
	// Init service
	SerialDevice := serialnetwork.NewSerialDevice()
	// SerialDevice.Init(serialConfig, rxLength)

	// Get channel
	// var rxChannel chan []byte
	// if rxChannel = SerialDevice.GetRxChannel(); rxChannel != nil {
	// 	logger.Info("Got RxChannel")
	// }
	// var txChannel chan []byte
	// if txChannel = SerialDevice.GetTxChannel(); txChannel != nil {
	// 	logger.Info("Got TxChannel")
	// }
	// var txWroteChannel chan []byte
	// if txWroteChannel = SerialDevice.GetTxWroteChannel(); txWroteChannel != nil {
	// 	logger.Info("Got TxWroteChannel")
	// }

	// Start channel handler service
	SerialDevice.RxResponseServer(serverHost)
	SerialDevice.TxRequestServer(deviceHostPort)

	// Test channel
	// txChannel <- message
	// logger.Infof("txWroteChannel % x", <-txWroteChannel)
	forever := make(chan bool)
	<-forever
}
