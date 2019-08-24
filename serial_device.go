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

type SerialDeviceConfig struct {
	Name        string
	Baud        int
	ReadTimeout int
	Size        byte
	Parity      byte
	StopBits    byte
	RxLength    int
	ServerHost  string
}

func NewSerialDevice() *SerialDevice {
	return &SerialDevice{
		port:           nil,
		serialConfig:   nil,
		startable:      false,
		rxChannel:      nil,
		rxLength:       1,
		txChannel:      nil,
		txWroteChannel: nil,
		serverHost:     "",
		deviceHostPort: "",
	}
}

func (sd *SerialDevice) Init(serialConfig serial.Config, rxLength int) bool {
	return sd.init(serialConfig, rxLength)
}

// Serial Rx

func (sd *SerialDevice) GetRxChannel() chan []byte {
	return sd.rxChannel
}

func (sd *SerialDevice) RxResponseServer(serverHost string) {
	sd.SetServerHost(serverHost)
	go sd.rxResponseServer(serverHost)
}

func (sd *SerialDevice) SetServerHost(serverHost string) {
	sd.serverHost = serverHost
}

// Serial Tx

func (sd *SerialDevice) GetTxChannel() chan []byte {
	return sd.txChannel
}

func (sd *SerialDevice) GetTxWroteChannel() chan []byte {
	return sd.getTxWroteChannel()
}

func (sd *SerialDevice) TxRequestServer(deviceHostPort string) {
	go sd.txRequestServer(deviceHostPort)
}
