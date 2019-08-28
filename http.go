package serialnetwork

import (
	. "bytes"
	"encoding/json"
	"errors"
	. "github.com/srleohung/serialnetwork/tools"
	"io/ioutil"
	"net/http"
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
			if httpLogger.IsErr(sd.init(serialDeviceConfig)) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sd *SerialDevice) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (sd *SerialDevice) responseToServer() {
	for {
		sd.rxResponse(<-sd.rxChannel)
	}
}

func (sd *SerialDevice) requestFromServer(deviceAddr string) {
	sd.deviceAddr = deviceAddr
	http.HandleFunc(HTTP_SERIAL_INIT_PATH, sd.initAPI)
	http.HandleFunc(HTTP_SERIAL_PING_PATH, sd.ping)
	http.HandleFunc(HTTP_SERIAL_TX_PATH, sd.txRequest)
	http.HandleFunc(HTTP_SERIAL_TX_RX_PATH, sd.txRequestAndRxResponse)
	err := http.ListenAndServe(sd.deviceAddr, nil)
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
		if <-sd.txWroteChannel {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sd *SerialDevice) txRequestAndRxResponse(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		sd.txChannel <- txRequest
		if <-sd.txWroteChannel {
			w.Write(<-sd.rxChannel)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Serial server

func (ss *SerialServer) init() error {
	ss.rxChannel = make(chan []byte)
	ss.txChannel = make(chan []byte)
	if !ss.startable {
		ss.startable = true
	}
	return nil
}

func (ss *SerialServer) ping() bool {
	_, err := http.Get(ss.deviceHost + HTTP_SERIAL_PING_PATH)
	return !httpLogger.IsErr(err)
}

func (ss *SerialServer) initSerialDevice(serialDeviceConfig SerialDeviceConfig) error {
	bytes, err := json.Marshal(serialDeviceConfig)
	if httpLogger.IsErr(err) {
		return err
	}
	resp, err := http.Post(ss.deviceHost+HTTP_SERIAL_INIT_PATH, HTTP_CONTENT_JSON_TYPE, NewBuffer(bytes))
	if httpLogger.IsErr(err) {
		return errors.New("Unable to connect serial device")
	}
	switch resp.StatusCode {
	case http.StatusBadRequest:
		return errors.New("Unable to open serial port")
	case http.StatusNotFound:
		return errors.New("Unable to connect serial device")
	default:
		return nil
	}
}

func (ss *SerialServer) responseFromDevice(serverAddr string) {
	ss.serverAddr = serverAddr
	http.HandleFunc(HTTP_SERIAL_RX_PATH, ss.rxResponse)
	err := http.ListenAndServe(ss.serverAddr, nil)
	httpLogger.IsErr(err)
}

func (ss *SerialServer) requestToDevice() {
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

func (ss *SerialServer) txRequest(bytes []byte) bool {
	httpLogger.Info(ss.deviceHost)
	resp, err := http.Post(ss.deviceHost+HTTP_SERIAL_TX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	if httpLogger.IsErr(err) || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func (ss *SerialServer) txRequestAndRxResponse(bytes []byte) []byte {
	rxResponse, err := http.Post(ss.deviceHost+HTTP_SERIAL_TX_RX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
	defer rxResponse.Body.Close()
	rx, err := ioutil.ReadAll(rxResponse.Body)
	httpLogger.IsErr(err)
	return rx
}
