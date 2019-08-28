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
	deviceAddr     string
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

const (
	Stop1     byte = 1
	Stop1Half byte = 15
	Stop2     byte = 2
)

const (
	ParityNone  byte = 'N'
	ParityOdd   byte = 'O'
	ParityEven  byte = 'E'
	ParityMark  byte = 'M' // parity bit is always 1
	ParitySpace byte = 'S' // parity bit is always 0
)

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
		deviceAddr:     "",
	}
}

func (sd *SerialDevice) Init(serialDeviceConfig SerialDeviceConfig) bool {
	return sd.init(serialDeviceConfig)
}

// Serial Rx

func (sd *SerialDevice) GetRxChannel() chan []byte {
	return sd.rxChannel
}

func (sd *SerialDevice) ResponseToServer(serverHost string) {
	sd.SetServerHost(serverHost)
	go sd.responseToServer()
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

func (sd *SerialDevice) RequestFromServer(deviceAddr string) {
	go sd.requestFromServer(deviceAddr)
}
