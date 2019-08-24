package serialnetwork

import (
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
	"time"
)

var serialLogger Logger = NewLogger("serial")

func (sd *SerialDevice) getTxWroteChannel() chan []byte {
	if sd.txWroteChannel == nil {
		sd.txWroteChannel = make(chan []byte)
	}
	return sd.txWroteChannel
}

func (sd *SerialDevice) init(serialConfig serial.Config, rxLength int) bool {
	sd.port = nil
	sd.rxLength = rxLength
	sd.serialConfig = &serialConfig
	sd.rxChannel = make(chan []byte)
	sd.txChannel = make(chan []byte)
	if port, err := serial.OpenPort(&serialConfig); !serialLogger.IsErr(err) {
		sd.port = port
	}
	if !sd.startable {
		go sd.serialRX()
		go sd.serialTX()
		sd.startable = true
	}
	return sd.startable
}

func (sd *SerialDevice) serialRX() {
	for {
		bytes := make([]byte, sd.rxLength)
		if _, err := sd.port.Read(bytes); serialLogger.IsErr(err) {
			err = sd.port.Close()
			serialLogger.IsErr(err)
			sd.port = nil
		}
		select {
		case sd.rxChannel <- bytes:
		case <-time.After(time.Second):
			continue
		}
	}
}

func (sd *SerialDevice) serialTX() {
	for {
		if sd.port == nil {
			if port, err := serial.OpenPort(sd.serialConfig); serialLogger.IsErr(err) {
				time.Sleep(time.Second)
				continue
			} else {
				sd.port = port
			}
		}
		select {
		case bytes := <-sd.txChannel:
			serialLogger.Debugf("% x", bytes)
			if _, err := sd.port.Write(bytes); serialLogger.IsErr(err) {
				err = sd.port.Close()
				serialLogger.IsErr(err)
				sd.port = nil
				continue
			}
			if sd.txWroteChannel != nil {
				sd.txWroteChannel <- bytes
			}
		case <-time.After(time.Second):
			continue
		}
	}
}
