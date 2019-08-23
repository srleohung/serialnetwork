package serialnetwork

import (
	. "bytes"
	. "github.com/srleohung/serialnetwork/tools"
	"io/ioutil"
	"net/http"
)

var networkLogger Logger = NewLogger("network")

const (
	SERIAL_NETWORK_SERVER_HOST string = "http://localhost:9876"
	SERIAL_DEVICE_SERVER_HOST  string = "http://localhost:9877"

	RX_RESPONSE string = "/serial/rx/response"

	TX_REQUEST                 string = "/serial/tx/request"
	TX_REQUEST_AND_RX_RESPONSE string = "/serial/tx/requestandresponse"

	CONTENT_TYPE string = ""
)

func (sn *SerialDevice) rxResponseServer() {
	for {
		sn.rxResponse(SERIAL_NETWORK_SERVER_HOST+RX_RESPONSE, CONTENT_TYPE, <-sn.rxChannel)
	}
}

func (sn *SerialDevice) rxResponse(url, contentType string, bytes []byte) {
	_, err := http.Post(url, contentType, NewBuffer(bytes))
	IsError(err)
}

func (sn *SerialDevice) txRequestServer() {
	sn.getTxWroteChannel()
	http.HandleFunc(TX_REQUEST, sn.txRequest)
	http.HandleFunc(TX_REQUEST_AND_RX_RESPONSE, sn.txRequestAndRxResponse)
	err := http.ListenAndServe(":9877", nil)
	IsError(err)
}

func (sn *SerialDevice) txRequest(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !IsError(err) {
		networkLogger.Debug(txRequest)
		sn.txChannel <- txRequest
		w.Write(<-sn.txWroteChannel)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sn *SerialDevice) txRequestAndRxResponse(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !IsError(err) {
		sn.txChannel <- txRequest
		<-sn.txWroteChannel
		w.Write(<-sn.rxChannel)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ss *SerialServer) rxResponseServer() {
	http.HandleFunc(RX_RESPONSE, ss.rxResponse)
	err := http.ListenAndServe(":9876", nil)
	IsError(err)
}

func (ss *SerialServer) rxResponse(w http.ResponseWriter, r *http.Request) {
	if rxResponse, err := ioutil.ReadAll(r.Body); !IsError(err) {
		ss.rxChannel <- rxResponse
		w.Write(rxResponse)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (ss *SerialServer) txRequest(bytes []byte) {
	_, err := http.Post(SERIAL_DEVICE_SERVER_HOST+TX_REQUEST, CONTENT_TYPE, NewBuffer(bytes))
	IsError(err)
}

func (ss *SerialServer) txRequestAndRxResponse(bytes []byte) []byte {
	rxResponse, err := http.Post(SERIAL_DEVICE_SERVER_HOST+TX_REQUEST, CONTENT_TYPE, NewBuffer(bytes))
	IsError(err)
	defer rxResponse.Body.Close()
	rx, err := ioutil.ReadAll(rxResponse.Body)
	IsError(err)
	return rx
}

func (ss *SerialServer) txRequestServer() {
	for {
		ss.txRequest(<-ss.txChannel)
	}
}
