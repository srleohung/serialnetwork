package serialnetwork

import (
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
)

type SerialDevice struct {
	port           *serial.Port
	serialConfig   *serial.Config
	startable      bool
	rxChannel      chan []byte
	rxLength       int
	txChannel      chan []byte
	txWroteChannel chan []byte
}

func NewSerialDevice(serialConfig serial.Config, rxLength int) *SerialDevice {
	port, err := serial.OpenPort(&serialConfig)
	if IsError(err) {
		return nil
	}
	return &SerialDevice{
		port:           port,
		serialConfig:   &serialConfig,
		startable:      false,
		rxChannel:      nil,
		rxLength:       rxLength,
		txChannel:      nil,
		txWroteChannel: nil,
	}
}

func (sn *SerialDevice) GetRxChannel() chan []byte {
	return sn.rxChannel
}

func (sn *SerialDevice) GetTxChannel() chan []byte {
	return sn.txChannel
}

func (sn *SerialDevice) GetTxWroteChannel() chan []byte {
	return sn.getTxWroteChannel()
}

func (sn *SerialDevice) Init() bool {
	return sn.init()
}

func (sn *SerialDevice) RxResponseServer() {
	sn.rxResponseServer()
}

func (sn *SerialDevice) TxRequestServer() {
	sn.txRequestServer()
}
