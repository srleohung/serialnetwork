package serialnetwork

import (
	. "bytes"
	"encoding/json"
	. "github.com/srleohung/serialnetwork/tools"
	"github.com/tarm/serial"
	"io/ioutil"
	"net/http"
	"time"
)

var httpLogger Logger = NewLogger("http")

const (
	HTTP_SERIAL_INIT_PATH  string = "/serialnetwork/serial/init"
	HTTP_SERIAL_PING_PATH  string = "/serialnetwork/serial/ping"
	HTTP_SERIAL_RX_PATH    string = "/serialnetwork/serial/rx"
	HTTP_SERIAL_TX_PATH    string = "/serialnetwork/serial/tx"
	HTTP_SERIAL_TX_RX_PATH string = "/serialnetwork/serial/tx/rx"
	HTTP_CONTENT_TYPE      string = ""
	HTTP_CONTENT_JSON_TYPE string = "application/json"
)

// Serial device

func (sd *SerialDevice) initAPI(w http.ResponseWriter, r *http.Request) {
	if config, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		var serialDeviceConfig SerialDeviceConfig
		if err = json.Unmarshal(config, &serialDeviceConfig); !httpLogger.IsErr(err) {
			if serialDeviceConfig.ServerHost != "" {
				sd.serverHost = serialDeviceConfig.ServerHost
			}
			var serialConfig serial.Config
			if serialDeviceConfig.Name != "" {
				serialConfig.Name = serialDeviceConfig.Name
			}
			if serialDeviceConfig.Baud != 0 {
				serialConfig.Baud = serialDeviceConfig.Baud
			}
			if serialDeviceConfig.ReadTimeout != 0 {
				serialConfig.ReadTimeout = time.Duration(serialDeviceConfig.ReadTimeout) * time.Millisecond
			}
			if serialDeviceConfig.Size != 0 {
				serialConfig.Size = serialDeviceConfig.Size
			}
			if serialDeviceConfig.Parity != 0 {
				serialConfig.Parity = serial.Parity(serialDeviceConfig.Parity)
			}
			if serialDeviceConfig.StopBits != 0 {
				serialConfig.StopBits = serial.StopBits(serialDeviceConfig.StopBits)
			}
			sd.serialConfig = &serialConfig
			if serialDeviceConfig.RxLength == 0 {
				serialDeviceConfig.RxLength = 1
			}
			if !sd.init(serialConfig, serialDeviceConfig.RxLength) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sd *SerialDevice) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (sd *SerialDevice) responseToServer(serverHost string) {
	for {
		sd.rxResponse(<-sd.rxChannel)
	}
}

func (sd *SerialDevice) requestFromServer(deviceHostPort string) {
	sd.deviceHostPort = deviceHostPort
	sd.getTxWroteChannel()
	http.HandleFunc(HTTP_SERIAL_INIT_PATH, sd.initAPI)
	http.HandleFunc(HTTP_SERIAL_PING_PATH, sd.ping)
	http.HandleFunc(HTTP_SERIAL_TX_PATH, sd.txRequest)
	http.HandleFunc(HTTP_SERIAL_TX_RX_PATH, sd.txRequestAndRxResponse)
	err := http.ListenAndServe(sd.deviceHostPort, nil)
	httpLogger.IsErr(err)
}

func (sd *SerialDevice) rxResponse(bytes []byte) {
	_, err := http.Post(sd.serverHost+HTTP_SERIAL_RX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
}

func (sd *SerialDevice) txRequest(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		sd.txChannel <- txRequest
		w.Write(<-sd.txWroteChannel)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sd *SerialDevice) txRequestAndRxResponse(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		sd.txChannel <- txRequest
		<-sd.txWroteChannel
		w.Write(<-sd.rxChannel)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Serial server

func (ss *SerialServer) init() bool {
	ss.rxChannel = make(chan []byte)
	ss.txChannel = make(chan []byte)
	if !ss.startable {
		ss.startable = true
	}
	return ss.startable
}

func (ss *SerialServer) ping() bool {
	_, err := http.Get(ss.deviceHost + HTTP_SERIAL_PING_PATH)
	return httpLogger.IsErr(err)
}

func (ss *SerialServer) initSerialDevice(serialDeviceConfig SerialDeviceConfig) {
	if serialDeviceConfigJson, err := json.Marshal(serialDeviceConfig); !httpLogger.IsErr(err) {
		_, err := http.Post(ss.deviceHost+HTTP_SERIAL_INIT_PATH, HTTP_CONTENT_JSON_TYPE, NewBuffer(serialDeviceConfigJson))
		httpLogger.IsErr(err)
	}
}

func (ss *SerialServer) responseFromDevice(serverHostPort string) {
	ss.serverHostPort = serverHostPort
	http.HandleFunc(HTTP_SERIAL_RX_PATH, ss.rxResponse)
	err := http.ListenAndServe(ss.serverHostPort, nil)
	httpLogger.IsErr(err)
}

func (ss *SerialServer) requestToDevice(deviceHost string) {
	for {
		ss.txRequest(<-ss.txChannel)
	}
}

func (ss *SerialServer) rxResponse(w http.ResponseWriter, r *http.Request) {
	if rxResponse, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		ss.rxChannel <- rxResponse
		w.Write(rxResponse)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ss *SerialServer) txRequest(bytes []byte) []byte {
	httpLogger.Info(ss.deviceHost)

	txRequest, err := http.Post(ss.deviceHost+HTTP_SERIAL_TX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
	defer txRequest.Body.Close()
	tx, err := ioutil.ReadAll(txRequest.Body)
	httpLogger.IsErr(err)
	return tx
}

func (ss *SerialServer) txRequestAndRxResponse(bytes []byte) []byte {
	rxResponse, err := http.Post(ss.deviceHost+HTTP_SERIAL_TX_RX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
	defer rxResponse.Body.Close()
	rx, err := ioutil.ReadAll(rxResponse.Body)
	httpLogger.IsErr(err)
	return rx
}
