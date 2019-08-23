package serialnetwork

type SerialServer struct {
	rxChannel chan []byte
	txChannel chan []byte
}

func NewSerialServer() *SerialServer {
	return &SerialServer{
		rxChannel: make(chan []byte),
		txChannel: make(chan []byte),
	}
}

func (ss *SerialServer) GetRxChannel() chan []byte {
	return ss.rxChannel
}

func (ss *SerialServer) GetTxChannel() chan []byte {
	return ss.txChannel
}

func (ss *SerialServer) RxResponseServer() {
	ss.rxResponseServer()
}

func (ss *SerialServer) TxRequest(bytes []byte) {
	ss.txRequest(bytes)
}

func (ss *SerialServer) TxRequestAndRxResponse(bytes []byte) []byte {
	return ss.txRequestAndRxResponse(bytes)
}
