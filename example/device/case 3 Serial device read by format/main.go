package main

import (
	"github.com/srleohung/serialnetwork"
	. "github.com/srleohung/serialnetwork/tools"
	"time"
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
		StartByte:         []byte{0x011},
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

// Test format 1
var message1 []byte = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

// Test format 2
var message2 []byte = []byte{0x11, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

// Test format 3
var message3 []byte = []byte{0x21, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

// Test format error
var message4 []byte = []byte{0x31, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}

var rxChannel chan []byte
var txChannel chan []byte
var txWroteChannel chan bool

func main() {
	d := serialnetwork.NewDevice()

	// ***** Set rx format *****
	d.SetRxFormat(rxFormat)

	// ***** Test rx format *****
	bytes := d.TestRxFormats(message1)
	if len(bytes) != len(message1) {
		logger.Emerg("Unknown massage in Rx format.")
	} else {
		logger.Info("Can handle the wrong massage 1 in Rx format")
	}
	bytes = d.TestRxFormats(message2)
	if len(bytes) != len(message2) {
		logger.Emerg("Unknown massage in Rx format.")
	} else {
		logger.Info("Can handle the wrong massage 2 in Rx format")
	}
	bytes = d.TestRxFormats(message3)
	if len(bytes) != len(message3) {
		logger.Emerg("Unknown massage in Rx format.")
	} else {
		logger.Info("Can handle the wrong massage 3 in Rx format")
	}
	bytes = d.TestRxFormats(message4)
	if len(bytes) != len(message4) {
		logger.Infof("Can handle the wrong massage in Rx format. After filter length %v", len(bytes))
	} else {
		logger.Emerg("The wrong massage in Rx format cannot be processed.")
	}

	// ***** Init serial device *****
	err := d.Init(config)
	if err != nil {
		logger.Emerg(err)
	}

	// ***** Get channel *****
	if rxChannel = d.GetRxChannel(); rxChannel != nil {
		logger.Info("Got RxChannel")
	}
	if txChannel = d.GetTxChannel(); txChannel != nil {
		logger.Info("Got TxChannel")
	}
	/*
		If you want to check your tx message have been sent,
		you can use Tx wrote channel and listen to it.
	*/
	if txWroteChannel = d.GetTxWroteChannel(); txWroteChannel != nil {
		logger.Info("Got TxWroteChannel")
	}

	// ***** Start channel handler service *****
	/*
		If you use api calls to control and you don't need to automatic response,
		please don't run this function.
	*/
	/*
		d.ResponseToServer(ServerHost)
		d.RequestFromServer(DeviceAddr)
		forever := make(chan bool)
		<-forever
	*/

	// ***** Test channel *****
	go func() {
		for {
			txChannel <- message1
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
