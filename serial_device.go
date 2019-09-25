package serialnetwork

import (
	"github.com/gorilla/websocket"
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
	connection     *websocket.Conn
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
	StartByte           []byte
	EndByte             []byte
	LengthByteIndex     int
	LengthByteMissing   int
	LengthFixed         int
	LengthHighByteIndex int
	LengthDecoder       string
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

// Initialize the serial port and channel

func (d *Device) Init(config Config) error {
	return d.init(config)
}

// Serial port control

func (d *Device) OpenPort() bool {
	return d.openPort()
}

func (d *Device) ClosePort() bool {
	return d.closePort()
}

// Get serial read channel

func (d *Device) GetRxChannel() chan []byte {
	return d.rxChannel
}

// Get serial write channel

func (d *Device) GetTxChannel() chan []byte {
	return d.txChannel
}

func (d *Device) GetTxWroteChannel() chan bool {
	return d.txWroteChannel
}

// Set the serial read format

func (d *Device) SetRxFormat(rxFormats []RxFormat) {
	d.setRxFormat(rxFormats)
}

// Test serial read format

func (d *Device) TestRxFormats(bytes []byte) []byte {
	return d.testRxFormats(bytes)
}

// Establish a network socket connection

func (d *Device) NewWebSocketClient(serverHost string) error {
	d.SetServerHost(serverHost)
	return d.newWebSocketClient(serverHost)
}

// Establish an api server connection

func (d *Device) ResponseToServer(serverHost string) {
	d.SetServerHost(serverHost)
	go d.responseToServer()
}

func (d *Device) SetServerHost(serverHost string) {
	d.serverHost = serverHost
}

func (d *Device) RequestFromServer(deviceAddr string) {
	go d.requestFromServer(deviceAddr)
}
