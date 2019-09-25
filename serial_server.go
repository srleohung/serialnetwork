package serialnetwork

import (
	"github.com/gorilla/websocket"
)

type Server struct {
	rxChannel  chan []byte
	txChannel  chan []byte
	startable  bool
	serverAddr string
	deviceHost string
	upgrader   *websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		rxChannel:  nil,
		txChannel:  nil,
		serverAddr: "",
		deviceHost: "",
		startable:  false,
	}
}

// Initialize channel

func (s *Server) Init() error {
	return s.init()
}

// Get serial read channel

func (s *Server) GetRxChannel() chan []byte {
	return s.rxChannel
}

// Get serial write channel

func (s *Server) GetTxChannel() chan []byte {
	return s.txChannel
}

// Establish a network socket connection server

func (s *Server) NewWebSocketServer(serverAddr string) {
	go s.newWebSocketServer(serverAddr)
}

// Establish api server connection server

func (s *Server) ResponseFromDevice(serverAddr string) {
	go s.responseFromDevice(serverAddr)
}

func (s *Server) RequestToDevice(deviceHost string) {
	s.SetDeviceHost(deviceHost)
	go s.requestToDevice()
}

// Api calls the device to test the connection

func (s *Server) Ping() bool {
	return s.ping()
}

// Api calls the serial port of the initialization device

func (s *Server) InitDevice(config Config) error {
	return s.initDevice(config)
}

// Set api to call device host

func (s *Server) SetDeviceHost(deviceHost string) {
	s.deviceHost = deviceHost
}

// Api call to write a message to the device

func (s *Server) TxRequest(bytes []byte) bool {
	return s.txRequest(bytes)
}

func (s *Server) TxRequestAndRxResponse(bytes []byte) []byte {
	return s.txRequestAndRxResponse(bytes)
}
