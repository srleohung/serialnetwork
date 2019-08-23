package serialnetwork

type SerialServer struct {
	rxChannel      chan []byte
	txChannel      chan []byte
	startable      bool
	serverHostPort string
	deviceHost     string
}

func NewSerialServer(serverHostPort, deviceHost string) *SerialServer {
	return &SerialServer{
		rxChannel:      nil,
		txChannel:      nil,
		serverHostPort: serverHostPort,
		deviceHost:     deviceHost,
		startable:      false,
	}
}

func (ss *SerialServer) Init() bool {
	return ss.init()
}

// Serial Rx

func (ss *SerialServer) GetRxChannel() chan []byte {
	return ss.rxChannel
}

func (ss *SerialServer) RxResponseServer() {
	ss.rxResponseServer()
}

// Serial Tx

func (ss *SerialServer) GetTxChannel() chan []byte {
	return ss.txChannel
}

func (ss *SerialServer) TxRequest(bytes []byte) []byte {
	return ss.txRequest(bytes)
}

func (ss *SerialServer) TxRequestAndRxResponse(bytes []byte) []byte {
	return ss.txRequestAndRxResponse(bytes)
}

func (ss *SerialServer) TxRequestServer() {
	ss.txRequestServer()
}
