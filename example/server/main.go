package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
)

var serverHostPost string = ":9876"

var deviceHost string = "http://localhost:9877"

var logger Logger = NewLogger("main")

var message []byte = []byte("9876")

var serialDeviceConfig serialnetwork.SerialDeviceConfig = serialnetwork.SerialDeviceConfig{
	Name:       "/dev/tty.SLAB_USBtoUART",
	Baud:       115200,
	Size:       8,
	RxLength:   1,
	ServerHost: "http://localhost:9876",
}

func main() {
	// Init service
	SerialServer := serialnetwork.NewSerialServer()
	SerialServer.Init()

	// Get channel
	var rxChannel chan []byte
	if rxChannel = SerialServer.GetRxChannel(); rxChannel != nil {
		logger.Info("Got RxChannel")
	}
	var txChannel chan []byte
	if txChannel = SerialServer.GetTxChannel(); txChannel != nil {
		logger.Info("Got TxChannel")
	}

	// Test service
	// SerialServer.SetDeviceHost(deviceHost)
	// logger.Infof("TxRequest % x", SerialServer.TxRequest(message))
	// logger.Infof("TxRequestAndRxResponse % x", SerialServer.TxRequestAndRxResponse(message))

	// Start channel handler service
	SerialServer.RxResponseServer(serverHostPost)
	SerialServer.TxRequestServer(deviceHost)

	// Init serial device
	SerialServer.InitSerialDevice(serialDeviceConfig)

	// Test channel
	txChannel <- message
	for {
		logger.Infof("rxChannel % x", <-rxChannel)
	}
}
