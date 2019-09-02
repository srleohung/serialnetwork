package serialnetwork

import (
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
	"time"
)

var serialLogger Logger = NewLogger("serial")

func (d *Device) init(config Config) error {
	d.port = nil
	d.rxChannel = nil
	d.txChannel = nil
	d.txWroteChannel = nil
	d.rxChannel = make(chan []byte)
	d.txChannel = make(chan []byte)
	d.txWroteChannel = make(chan bool, 1)
	if config.RxBuffer > 0 {
		d.rxBuffer = config.RxBuffer
	}
	if config.ServerHost != "" {
		d.serverHost = config.ServerHost
	}
	var serialConfig serial.Config = serial.Config{
		Name:        config.Name,
		Baud:        config.Baud,
		ReadTimeout: config.ReadTimeout,
		Size:        config.Size,
		Parity:      serial.Parity(config.Parity),
		StopBits:    serial.StopBits(config.StopBits),
	}
	d.serialConfig = &serialConfig
	if port, err := serial.OpenPort(&serialConfig); !serialLogger.IsErr(err) {
		d.port = port
	} else {
		return err
	}
	if !d.startable {
		go d.serialRX()
		go d.serialTX()
		serialLogger.Debug("Start running serialRX and serialTX.")
		d.startable = true
	}
	return nil
}

func (d *Device) serialRX() {
	defer serialLogger.Debug("Unable to run serialRX.")
	var bytes []byte
	var aByte []byte
	var buffer []byte
	for {
		if !d.openPort() {
			continue
		}
		if d.rxFormats == nil || d.rxFormatter == nil {
			bytes = make([]byte, d.rxBuffer)
			if _, err := d.port.Read(bytes); serialLogger.IsErr(err) {
				d.closePort()
			}
			d.rxChannel <- bytes
			serialLogger.Debugf("% x", bytes)
		} else {
			aByte = make([]byte, 1)
			if _, err := d.port.Read(aByte); serialLogger.IsErr(err) {
				d.closePort()
			}
			if buffer, bytes = d.rxHandler(buffer, aByte[0]); bytes != nil {
				d.rxChannel <- bytes
				serialLogger.Debugf("% x", bytes)
			}
		}
	}
	d.startable = false
}

func (d *Device) serialTX() {
	defer serialLogger.Debug("Unable to run serialTX.")
	for {
		bytes := <-d.txChannel
		if len(d.txWroteChannel) > 0 {
			<-d.txWroteChannel
		}
		if !d.openPort() {
			d.txWroteChannel <- false
			continue
		}
		serialLogger.Debugf("% x", bytes)
		if n, err := d.port.Write(bytes); serialLogger.IsErr(err) || n != len(bytes) {
			d.closePort()
			d.txWroteChannel <- false
		} else {
			d.txWroteChannel <- true
		}
	}
	d.startable = false
}

func (d *Device) openPort() bool {
	if d.port == nil {
		port, err := serial.OpenPort(d.serialConfig)
		if serialLogger.IsErr(err) {
			time.Sleep(time.Second)
			return false
		}
		d.port = port
	}
	return true
}

func (d *Device) closePort() bool {
	if d.port != nil {
		err := d.port.Close()
		d.port = nil
		return !serialLogger.IsErr(err)
	}
	return true
}

type RxFormatter struct {
	number            int
	numberOfEnd       int
	formatterNumber   int
	formatterIndexMax int
	length            int
}

func (d *Device) setRxFormat(rxFormats []RxFormat) {
	d.rxFormats = rxFormats
	var rxFormatter RxFormatter = RxFormatter{
		number:            0,
		numberOfEnd:       0,
		formatterNumber:   0,
		formatterIndexMax: len(rxFormats),
		length:            0,
	}
	d.rxFormatter = &rxFormatter
}

func (d *Device) rxHandler(bytes []byte, aByte byte) ([]byte, []byte) {
	for i := d.rxFormatter.formatterNumber; i < d.rxFormatter.formatterIndexMax; i++ {
		if len(d.rxFormats[i].StartByte) > 0 && d.rxFormatter.number < len(d.rxFormats[i].StartByte) {
			if aByte == d.rxFormats[i].StartByte[d.rxFormatter.number] {
				d.rxFormatter.formatterNumber = i
				d.rxFormatter.number++
				return append(bytes, aByte), nil
			}
			continue
		}
		if len(d.rxFormats[i].EndByte) > 0 {
			if aByte == d.rxFormats[i].EndByte[d.rxFormatter.numberOfEnd] {
				d.rxFormatter.numberOfEnd++
				if len(d.rxFormats[i].EndByte) == d.rxFormatter.numberOfEnd {
					d.rxFormatter.number = 0
					d.rxFormatter.numberOfEnd = 0
					d.rxFormatter.formatterNumber = 0
					d.rxFormatter.length = 0
					return nil, append(bytes, aByte)
				}
			}
			return append(bytes, aByte), nil
		}
		if d.rxFormats[i].LengthFixed > 0 {
			d.rxFormatter.length = d.rxFormats[i].LengthFixed
		} else {
			if d.rxFormatter.number == d.rxFormats[i].LengthByteIndex {
				d.rxFormatter.length = int(aByte) + d.rxFormats[i].LengthByteMissing
			}
		}
		d.rxFormatter.number++
		if d.rxFormatter.length != 0 && d.rxFormatter.number >= d.rxFormatter.length {
			d.rxFormatter.number = 0
			d.rxFormatter.numberOfEnd = 0
			d.rxFormatter.formatterNumber = 0
			d.rxFormatter.length = 0
			return nil, append(bytes, aByte)
		}
		return append(bytes, aByte), nil
	}
	return nil, nil
}

func (d *Device) testRxFormats(bytes []byte) []byte {
	var buffer []byte
	var output []byte
	if d.rxFormats == nil || d.rxFormatter == nil {
		return nil
	}
	for i := 0; i < len(bytes); i++ {
		buffer, output = d.rxHandler(buffer, bytes[i])
		serialLogger.Debugf("index %v", i)
		serialLogger.Debugf("buffer % x", buffer)
		serialLogger.Debugf("output % x", buffer)
	}
	return output
}
