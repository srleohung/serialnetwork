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

func main() {
	SerialDevice := serialnetwork.NewSerialDevice(serialConfig, 1)
	SerialDevice.Init()

	rxChannel := SerialDevice.GetRxChannel()
	txChannel := SerialDevice.GetTxChannel()
	txWroteChannel := SerialDevice.GetTxWroteChannel()

	go SerialDevice.RxResponseServer()
	go SerialDevice.TxRequestServer()

	txChannel <- []byte("101")
	logger.Info(<-txWroteChannel)
	logger.Info(<-rxChannel)
}
