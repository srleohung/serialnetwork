package serialnetwork

type SerialServer struct {
	rxChannel      chan []byte
	txChannel      chan []byte
	startable      bool
	serverHostPort string
	deviceHost     string
}

func NewSerialServer() *SerialServer {
	return &SerialServer{
		rxChannel:      nil,
		txChannel:      nil,
		serverHostPort: "",
		deviceHost:     "",
		startable:      false,
	}
}

func (ss *SerialServer) Init() bool {
	return ss.init()
}

func (ss *SerialServer) Ping() bool {
	return ss.ping()
}

func (ss *SerialServer) InitSerialDevice(serialDeviceConfig SerialDeviceConfig) {
	ss.initSerialDevice(serialDeviceConfig)
}

// Serial Rx

func (ss *SerialServer) GetRxChannel() chan []byte {
	return ss.rxChannel
}

func (ss *SerialServer) ResponseFromDevice(serverHostPort string) {
	go ss.responseFromDevice(serverHostPort)
}

// Serial Tx

func (ss *SerialServer) GetTxChannel() chan []byte {
	return ss.txChannel
}

func (ss *SerialServer) SetDeviceHost(deviceHost string) {
	ss.deviceHost = deviceHost
}

func (ss *SerialServer) TxRequest(bytes []byte) []byte {
	return ss.txRequest(bytes)
}

func (ss *SerialServer) TxRequestAndRxResponse(bytes []byte) []byte {
	return ss.txRequestAndRxResponse(bytes)
}

func (ss *SerialServer) RequestToDevice(deviceHost string) {
	ss.SetDeviceHost(deviceHost)
	go ss.requestToDevice(deviceHost)
}
