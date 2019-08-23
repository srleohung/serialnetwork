package serialnetwork

import (
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
	serverHost     string
	deviceHostPort string
}

func NewSerialDevice(serialConfig serial.Config, serverHost, deviceHostPort string, rxLength int) *SerialDevice {
	port, err := serial.OpenPort(&serialConfig)
	if serialLogger.IsErr(err) {
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
		serverHost:     serverHost,
		deviceHostPort: deviceHostPort,
	}
}

func (sn *SerialDevice) Init() bool {
	return sn.init()
}

// Serial Rx

func (sn *SerialDevice) GetRxChannel() chan []byte {
	return sn.rxChannel
}

func (sn *SerialDevice) RxResponseServer() {
	sn.rxResponseServer()
}

// Serial Tx

func (sn *SerialDevice) GetTxChannel() chan []byte {
	return sn.txChannel
}

func (sn *SerialDevice) GetTxWroteChannel() chan []byte {
	return sn.getTxWroteChannel()
}

func (sn *SerialDevice) TxRequestServer() {
	sn.txRequestServer()
}
