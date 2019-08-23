package serialnetwork

import (
	. "bytes"
	. "github.com/srleohung/serialnetwork/tools"
	"io/ioutil"
	"net/http"
)

var httpLogger Logger = NewLogger("http")

const (
	HTTP_SERIAL_RX_PATH    string = "/serial/rx"
	HTTP_SERIAL_TX_PATH    string = "/serial/tx"
	HTTP_SERIAL_TX_RX_PATH string = "/serial/tx/rx"
	HTTP_CONTENT_TYPE      string = ""
)

// Serial device

func (sn *SerialDevice) rxResponseServer() {
	for {
		sn.rxResponse(<-sn.rxChannel)
	}
}

func (sn *SerialDevice) rxResponse(bytes []byte) {
	_, err := http.Post(sn.serverHost+HTTP_SERIAL_RX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
}

func (sn *SerialDevice) txRequestServer() {
	sn.getTxWroteChannel()
	http.HandleFunc(HTTP_SERIAL_TX_PATH, sn.txRequest)
	http.HandleFunc(HTTP_SERIAL_TX_RX_PATH, sn.txRequestAndRxResponse)
	err := http.ListenAndServe(sn.deviceHostPort, nil)
	httpLogger.IsErr(err)
}

func (sn *SerialDevice) txRequest(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		sn.txChannel <- txRequest
		w.Write(<-sn.txWroteChannel)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sn *SerialDevice) txRequestAndRxResponse(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		sn.txChannel <- txRequest
		<-sn.txWroteChannel
		w.Write(<-sn.rxChannel)
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

func (ss *SerialServer) rxResponseServer() {
	http.HandleFunc(HTTP_SERIAL_RX_PATH, ss.rxResponse)
	err := http.ListenAndServe(ss.serverHostPort, nil)
	httpLogger.IsErr(err)
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

func (ss *SerialServer) txRequestServer() {
	for {
		ss.txRequest(<-ss.txChannel)
	}
}
