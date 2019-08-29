package serialnetwork

import (
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
	"time"
)

var serialLogger Logger = NewLogger("serial")

func (sd *Device) init(config Config) error {
	sd.rxChannel = make(chan []byte)
	sd.txChannel = make(chan []byte)
	sd.txWroteChannel = make(chan bool, 1)
	if config.RxBuffer > 0 {
		sd.rxBuffer = config.RxBuffer
	}
	if config.ServerHost != "" {
		sd.serverHost = config.ServerHost
	}
	var serialConfig serial.Config = serial.Config{
		Name:        config.Name,
		Baud:        config.Baud,
		ReadTimeout: config.ReadTimeout,
		Size:        config.Size,
		Parity:      serial.Parity(config.Parity),
		StopBits:    serial.StopBits(config.StopBits),
	}
	sd.serialConfig = &serialConfig
	sd.serialConfig = &serialConfig
	if port, err := serial.OpenPort(&serialConfig); !serialLogger.IsErr(err) {
		sd.port = port
	} else {
		return err
	}
	if !sd.startable {
		go sd.serialRX()
		go sd.serialTX()
		sd.startable = true
	}
	return nil
}

func (sd *Device) serialRX() {
	for {
		if !sd.openPort() {
			continue
		}
		bytes := make([]byte, sd.rxBuffer)
		if _, err := sd.port.Read(bytes); serialLogger.IsErr(err) {
			sd.closePort()
		}
		sd.rxChannel <- bytes
		serialLogger.Debugf("% x", bytes)
	}
}

func (sd *Device) serialTX() {
	for {
		if !sd.openPort() {
			continue
		}
		bytes := <-sd.txChannel
		serialLogger.Debugf("% x", bytes)
		if len(sd.txWroteChannel) > 0 {
			<-sd.txWroteChannel
		}
		if n, err := sd.port.Write(bytes); serialLogger.IsErr(err) || n != len(bytes) {
			sd.closePort()
			sd.txWroteChannel <- false
		} else {
			sd.txWroteChannel <- true
		}
	}
}

func (sd *Device) openPort() bool {
	if sd.port == nil {
		port, err := serial.OpenPort(sd.serialConfig)
		if serialLogger.IsErr(err) {
			time.Sleep(time.Second)
			return false
		}
		sd.port = port
	}
	return true
}

func (sd *Device) closePort() bool {
	if sd.port != nil {
		err := sd.port.Close()
		sd.port = nil
		return !serialLogger.IsErr(err)
	}
	return true
}
