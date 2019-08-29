package serialnetwork

import (
	"github.com/tarm/serial"
	"time"
)

type Device struct {
	port           *serial.Port
	serialConfig   *serial.Config
	startable      bool
	rxChannel      chan []byte
	rxBuffer       int
	txChannel      chan []byte
	txWroteChannel chan bool
	serverHost     string
	deviceAddr     string
}

type Config struct {
	Name        string
	Baud        int
	ReadTimeout time.Duration
	Size        byte
	Parity      byte
	StopBits    byte
	RxBuffer    int
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

func NewDevice() *Device {
	return &Device{
		port:           nil,
		serialConfig:   nil,
		startable:      false,
		rxChannel:      nil,
		rxBuffer:       1,
		txChannel:      nil,
		txWroteChannel: nil,
		serverHost:     "",
		deviceAddr:     "",
	}
}

func (sd *Device) Init(config Config) error {
	return sd.init(config)
}

// Serial Rx

func (sd *Device) GetRxChannel() chan []byte {
	return sd.rxChannel
}

func (sd *Device) ResponseToServer(serverHost string) {
	sd.SetServerHost(serverHost)
	go sd.responseToServer()
}

func (sd *Device) SetServerHost(serverHost string) {
	sd.serverHost = serverHost
}

// Serial Tx

func (sd *Device) GetTxChannel() chan []byte {
	return sd.txChannel
}

func (sd *Device) GetTxWroteChannel() chan bool {
	return sd.txWroteChannel
}

func (sd *Device) RequestFromServer(deviceAddr string) {
	go sd.requestFromServer(deviceAddr)
}

func (sd *Device) OpenPort() bool {
	return sd.openPort()
}

func (sd *Device) ClosePort() bool {
	return sd.closePort()
}
