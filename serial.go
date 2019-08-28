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

func (sd *SerialDevice) init(serialDeviceConfig SerialDeviceConfig) bool {
	if serialDeviceConfig.ServerHost != "" {
		sd.serverHost = serialDeviceConfig.ServerHost
	}
	var serialConfig serial.Config
	if serialDeviceConfig.Name != "" {
		serialConfig.Name = serialDeviceConfig.Name
	}
	if serialDeviceConfig.Baud != 0 {
		serialConfig.Baud = serialDeviceConfig.Baud
	}
	if serialDeviceConfig.ReadTimeout != 0 {
		serialConfig.ReadTimeout = time.Duration(serialDeviceConfig.ReadTimeout) * time.Millisecond
	}
	if serialDeviceConfig.Size != 0 {
		serialConfig.Size = serialDeviceConfig.Size
	}
	if serialDeviceConfig.Parity != 0 {
		serialConfig.Parity = serial.Parity(serialDeviceConfig.Parity)
	}
	if serialDeviceConfig.StopBits != 0 {
		serialConfig.StopBits = serial.StopBits(serialDeviceConfig.StopBits)
	}
	sd.serialConfig = &serialConfig
	if serialDeviceConfig.RxLength == 0 {
		serialDeviceConfig.RxLength = 1
	}
	sd.port = nil
	sd.rxLength = serialDeviceConfig.RxLength
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
