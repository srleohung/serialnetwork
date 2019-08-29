package serialnetwork

type Server struct {
	rxChannel  chan []byte
	txChannel  chan []byte
	startable  bool
	serverAddr string
	deviceHost string
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

func (s *Server) Init() error {
	return s.init()
}

func (s *Server) Ping() bool {
	return s.ping()
}

func (s *Server) InitDevice(config Config) error {
	return s.initDevice(config)
}

// Serial Rx

func (s *Server) GetRxChannel() chan []byte {
	return s.rxChannel
}

func (s *Server) ResponseFromDevice(serverAddr string) {
	go s.responseFromDevice(serverAddr)
}

// Serial Tx

func (s *Server) GetTxChannel() chan []byte {
	return s.txChannel
}

func (s *Server) SetDeviceHost(deviceHost string) {
	s.deviceHost = deviceHost
}

func (s *Server) TxRequest(bytes []byte) bool {
	return s.txRequest(bytes)
}

func (s *Server) TxRequestAndRxResponse(bytes []byte) []byte {
	return s.txRequestAndRxResponse(bytes)
}

func (s *Server) RequestToDevice(deviceHost string) {
	s.SetDeviceHost(deviceHost)
	go s.requestToDevice()
}
