package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
)

var logger Logger = NewLogger("main")

func main() {
	SerialServer := serialnetwork.NewSerialServer()
	go SerialServer.RxResponseServer()

	rxChannel := SerialServer.GetRxChannel()

	SerialServer.TxRequest([]byte("101"))
	logger.Info(<-rxChannel)
}
