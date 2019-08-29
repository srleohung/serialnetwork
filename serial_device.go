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
	rxFormats      []RxFormat
	rxFormatter    *RxFormatter
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

type RxFormat struct {
	StartByte         []byte
	EndByte           []byte
	LengthByteIndex   int
	LengthByteMissing int
	LengthFixed       int
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
		rxFormats:      nil,
		rxFormatter:    nil,
	}
}

func (d *Device) Init(config Config) error {
	return d.init(config)
}

// Serial Rx

func (d *Device) SetRxFormat(rxFormats []RxFormat) {
	d.setRxFormat(rxFormats)
}

func (d *Device) GetRxChannel() chan []byte {
	return d.rxChannel
}

func (d *Device) ResponseToServer(serverHost string) {
	d.SetServerHost(serverHost)
	go d.responseToServer()
}

func (d *Device) SetServerHost(serverHost string) {
	d.serverHost = serverHost
}

// Serial Tx

func (d *Device) GetTxChannel() chan []byte {
	return d.txChannel
}

func (d *Device) GetTxWroteChannel() chan bool {
	return d.txWroteChannel
}

func (d *Device) RequestFromServer(deviceAddr string) {
	go d.requestFromServer(deviceAddr)
}

func (d *Device) OpenPort() bool {
	return d.openPort()
}

func (d *Device) ClosePort() bool {
	return d.closePort()
}

// Test

func (d *Device) TestRxFormats(bytes []byte) []byte {
	return d.testRxFormats(bytes)
}
