package serialnetwork

import (
	"log"
	"testing"
	"time"
)

func TestWebSocketSerialServer(t *testing.T) {
	type args struct {
		serverHost string
		config     Config
		message    []byte
	}
	tests := []struct {
		device *Device
		server *Server
		args   args
	}{
		{
			device: NewDevice(),
			server: NewServer(),
			args: args{
				serverHost: ":9876",
				config: Config{
					Name: "/dev/ttyUSB0",
					Baud: 115200,
					// ReadTimeout: 1000 * time.Millisecond,
					Size:   8,
					Parity: ParityNone,
					/*
						ParityNone  Parity = 'N'
						ParityOdd   Parity = 'O'
						ParityEven  Parity = 'E'
						ParityMark  Parity = 'M' // parity bit is always 1
						ParitySpace Parity = 'S' // parity bit is always 0
					*/
					StopBits: Stop1,
					/*
						Stop1     StopBits = 1
						Stop1Half StopBits = 15
						Stop2     StopBits = 2
					*/
					RxBuffer: 1,
					// ServerHost: "",
				},
				message: []byte("Hello Serial Device"),
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			/* WebSocket Serial Server */
			err := tt.server.Init()
			if err != nil {
				t.Error(err)
				return
			}
			tt.server.NewWebSocketServer(tt.args.serverHost)
			rxServerChannel := tt.server.GetRxChannel()
			txServerChannel := tt.server.GetTxChannel()

			/* WebSocket Serial Device */
			if err := tt.device.Init(tt.args.config); err != nil {
				t.Error(err)
				return
			}
			if err := tt.device.NewWebSocketClient("localhost" + tt.args.serverHost); err != nil {
				t.Error(err)
				return
			}

			txServerChannel <- tt.args.message
			go func(c chan []byte) {
				for {
					log.Printf("%s", <-c)
				}
			}(rxServerChannel)
			<-time.After(5 * time.Second)
		})
	}
}
