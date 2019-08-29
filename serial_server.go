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

func (ss *Server) Init() error {
	return ss.init()
}

func (ss *Server) Ping() bool {
	return ss.ping()
}

func (ss *Server) InitDevice(config Config) error {
	return ss.initDevice(config)
}

// Serial Rx

func (ss *Server) GetRxChannel() chan []byte {
	return ss.rxChannel
}

func (ss *Server) ResponseFromDevice(serverAddr string) {
	go ss.responseFromDevice(serverAddr)
}

// Serial Tx

func (ss *Server) GetTxChannel() chan []byte {
	return ss.txChannel
}

func (ss *Server) SetDeviceHost(deviceHost string) {
	ss.deviceHost = deviceHost
}

func (ss *Server) TxRequest(bytes []byte) bool {
	return ss.txRequest(bytes)
}

func (ss *Server) TxRequestAndRxResponse(bytes []byte) []byte {
	return ss.txRequestAndRxResponse(bytes)
}

func (ss *Server) RequestToDevice(deviceHost string) {
	ss.SetDeviceHost(deviceHost)
	go ss.requestToDevice()
}
