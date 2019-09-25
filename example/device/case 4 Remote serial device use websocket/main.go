package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
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
	RxBuffer: 1,
	// ServerHost: "",
}

const ServerHost string = "localhost:9876"

func main() {
	// ***** Init serial device *****
	d := serialnetwork.NewDevice()
	if err := d.Init(config); err != nil {
		logger.Emerg(err)
	}

	// ***** Start channel handler service *****
	if err := d.NewWebSocketClient(ServerHost); err != nil {
		logger.Emerg(err)
	}

	// ***** Run forever *****
	forever := make(chan bool)
	<-forever
}
