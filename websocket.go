package serialnetwork

import (
	"github.com/gorilla/websocket"
	. "github.com/srleohung/serialnetwork/tools"
	"net/http"
)

var webSocketLogger Logger = NewLogger("websocket")

const (
	WEBSOCKET_NEW_CONNECTION_PATH string = "/serialnetwork/websocket/connect"
)

func (s *Server) newWebSocketServer(serverAddr string) {
	s.serverAddr = serverAddr
	s.upgrader = &websocket.Upgrader{CheckOrigin: s.checkOrigin}
	http.HandleFunc(WEBSOCKET_NEW_CONNECTION_PATH, s.newWebSocketConnection)
	err := http.ListenAndServe(s.serverAddr, nil)
	webSocketLogger.IsErr(err)
}

func (s *Server) checkOrigin(r *http.Request) bool {
	return true
}

func (s *Server) disconnect(c *websocket.Conn) {
	webSocketLogger.Debug("Disconnect.")
	c.Close()
}

func (s *Server) newWebSocketConnection(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if webSocketLogger.IsErr(err) {
		return
	}
	defer s.disconnect(c)
	go s.webSocketWriter(c)
	s.webSocketReader(c)
}

func (s *Server) webSocketReader(c *websocket.Conn) {
	webSocketLogger.Debug("Start running the websocket reader.")
	for {
		_, rx, err := c.ReadMessage()
		if webSocketLogger.IsErr(err) {
			break
		}
		s.rxChannel <- rx
	}
}

func (s *Server) webSocketWriter(c *websocket.Conn) {
	webSocketLogger.Debug("Start running the websocket writer.")
	for {
		if err := c.WriteMessage(websocket.TextMessage, <-s.txChannel); webSocketLogger.IsErr(err) {
			break
		}
	}
}

func (d *Device) newWebSocketClient(serverHost string) error {
	c, _, err := websocket.DefaultDialer.Dial("ws://"+serverHost+WEBSOCKET_NEW_CONNECTION_PATH, nil)
	if webSocketLogger.IsErr(err) {
		return err
	}
	d.connection = c
	go d.webSocketWriter(c)
	go d.webSocketReader(c)
	return nil
}

func (d *Device) webSocketReader(c *websocket.Conn) {
	webSocketLogger.Debug("Start running the websocket reader.")
	for {
		_, tx, err := d.connection.ReadMessage()
		if webSocketLogger.IsErr(err) {
			break
		}
		d.txChannel <- tx
	}
}

func (d *Device) webSocketWriter(c *websocket.Conn) {
	webSocketLogger.Debug("Start running the websocket writer.")
	for {
		if err := d.connection.WriteMessage(websocket.TextMessage, <-d.rxChannel); webSocketLogger.IsErr(err) {
			break
		}
	}
}
