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

func (d *Device) initDevice(w http.ResponseWriter, r *http.Request) {
	if body, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		var config Config
		if err = json.Unmarshal(body, &config); !httpLogger.IsErr(err) {
			if httpLogger.IsErr(d.init(config)) {
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

func (d *Device) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (d *Device) responseToServer() {
	for {
		d.rxResponse(<-d.rxChannel)
	}
}

func (d *Device) requestFromServer(deviceAddr string) {
	d.deviceAddr = deviceAddr
	http.HandleFunc(HTTP_SERIAL_INIT_PATH, d.initDevice)
	http.HandleFunc(HTTP_SERIAL_PING_PATH, d.ping)
	http.HandleFunc(HTTP_SERIAL_TX_PATH, d.txRequest)
	http.HandleFunc(HTTP_SERIAL_TX_RX_PATH, d.txRequestAndRxResponse)
	err := http.ListenAndServe(d.deviceAddr, nil)
	httpLogger.IsErr(err)
}

func (d *Device) rxResponse(bytes []byte) {
	_, err := http.Post(d.serverHost+HTTP_SERIAL_RX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
}

func (d *Device) txRequest(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		d.txChannel <- txRequest
		if <-d.txWroteChannel {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (d *Device) txRequestAndRxResponse(w http.ResponseWriter, r *http.Request) {
	if txRequest, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		httpLogger.Debugf("% x", txRequest)
		d.txChannel <- txRequest
		if <-d.txWroteChannel {
			w.Write(<-d.rxChannel)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Serial server

func (s *Server) init() error {
	s.rxChannel = make(chan []byte)
	s.txChannel = make(chan []byte)
	if !s.startable {
		s.startable = true
	}
	return nil
}

func (s *Server) ping() bool {
	_, err := http.Get(s.deviceHost + HTTP_SERIAL_PING_PATH)
	return !httpLogger.IsErr(err)
}

func (s *Server) initDevice(config Config) error {
	bytes, err := json.Marshal(config)
	if httpLogger.IsErr(err) {
		return err
	}
	resp, err := http.Post(s.deviceHost+HTTP_SERIAL_INIT_PATH, HTTP_CONTENT_JSON_TYPE, NewBuffer(bytes))
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

func (s *Server) responseFromDevice(serverAddr string) {
	s.serverAddr = serverAddr
	http.HandleFunc(HTTP_SERIAL_RX_PATH, s.rxResponse)
	err := http.ListenAndServe(s.serverAddr, nil)
	httpLogger.IsErr(err)
}

func (s *Server) requestToDevice() {
	for {
		s.txRequest(<-s.txChannel)
	}
}

func (s *Server) rxResponse(w http.ResponseWriter, r *http.Request) {
	if rxResponse, err := ioutil.ReadAll(r.Body); !httpLogger.IsErr(err) {
		s.rxChannel <- rxResponse
		w.Write(rxResponse)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) txRequest(bytes []byte) bool {
	resp, err := http.Post(s.deviceHost+HTTP_SERIAL_TX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	if httpLogger.IsErr(err) || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func (s *Server) txRequestAndRxResponse(bytes []byte) []byte {
	rxResponse, err := http.Post(s.deviceHost+HTTP_SERIAL_TX_RX_PATH, HTTP_CONTENT_TYPE, NewBuffer(bytes))
	httpLogger.IsErr(err)
	defer rxResponse.Body.Close()
	rx, err := ioutil.ReadAll(rxResponse.Body)
	httpLogger.IsErr(err)
	return rx
}
