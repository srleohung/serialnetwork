package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
)

var logger Logger = NewLogger("main")

var message []byte = []byte("9876")

func main() {
	// Init service
	SerialServer := serialnetwork.NewSerialServer(":9876", "http://localhost:9877")
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
	logger.Infof("TxRequest % x", SerialServer.TxRequest(message))
	// logger.Infof("TxRequestAndRxResponse % x", SerialServer.TxRequestAndRxResponse(message))

	// Start channel handler service
	go SerialServer.RxResponseServer()
	go SerialServer.TxRequestServer()

	// Test channel
	txChannel <- message
	for {
		logger.Infof("rxChannel % x", <-rxChannel)
	}
}
