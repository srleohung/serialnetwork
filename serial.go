package serialnetwork

import (
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
	"time"
)

var serialLogger Logger = NewLogger("serial")

func (sn *SerialDevice) getTxWroteChannel() chan []byte {
	if sn.txWroteChannel == nil {
		sn.txWroteChannel = make(chan []byte)
	}
	return sn.txWroteChannel
}

func (sn *SerialDevice) init() bool {
	sn.rxChannel = make(chan []byte)
	sn.txChannel = make(chan []byte)
	if !sn.startable {
		go sn.serialRX()
		go sn.serialTX()
		sn.startable = true
	}
	return sn.startable
}

func (sn *SerialDevice) serialRX() {
	for {
		bytes := make([]byte, sn.rxLength)
		if _, err := sn.port.Read(bytes); IsError(err) {
			err = sn.port.Close()
			IsError(err)
			sn.port = nil
		}
		select {
		case sn.rxChannel <- bytes:
		case <-time.After(time.Second):
			continue
		}
	}
}

func (sn *SerialDevice) serialTX() {
	for {
		if sn.port == nil {
			if port, err := serial.OpenPort(sn.serialConfig); IsError(err) {
				time.Sleep(1000 * time.Millisecond)
				continue
			} else {
				sn.port = port
			}
		}
		select {
		case bytes := <-sn.txChannel:
			serialLogger.Debug(bytes)
			if _, err := sn.port.Write(bytes); IsError(err) {
				err = sn.port.Close()
				IsError(err)
				sn.port = nil
				continue
			}
			if sn.txWroteChannel != nil {
				sn.txWroteChannel <- bytes
			}
		case <-time.After(time.Second):
			continue
		}
	}
}
