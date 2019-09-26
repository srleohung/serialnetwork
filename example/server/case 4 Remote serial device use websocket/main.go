package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
)

var logger Logger = NewLogger("main")

const ServerHost string = ":9876"

var rxChannel chan []byte
var txChannel chan []byte

var message []byte = []byte("test")

func main() {
	// ***** Init service *****
	s := serialnetwork.NewServer()
	err := s.Init()
	if err != nil {
		logger.Emerg(err)
	}

	// ***** Start the websocket server *****
	s.NewWebSocketServer(ServerHost)

	// ***** Get channel *****
	if rxChannel = s.GetRxChannel(); rxChannel != nil {
		logger.Info("Got RxChannel")
	}
	if txChannel = s.GetTxChannel(); txChannel != nil {
		logger.Info("Got TxChannel")
	}

	// ***** Test channel *****
	txChannel <- message
	for {
		logger.Infof("rxChannel % x", <-rxChannel)
	}
}
