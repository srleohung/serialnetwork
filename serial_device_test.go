package serialnetwork

import (
	"log"
	"reflect"
	"testing"
	"time"
)

func TestSerialDeviceCommunication(t *testing.T) {
	type args struct {
		config        Config
		startCommand  [][]byte
		holdCommand   []byte
		writeInterval int // time.Millisecond
		cancelTimeout int // time.Second
		readTimeout   int // time.Second
		rxFormat      []RxFormat
	}
	tests := []struct {
		device *Device
		args   args
	}{
		{
			device: NewDevice(),
			args: args{
				config: Config{
					Name:   "/dev/ttyUSB0",
					Baud:   115200,
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
				},
				readTimeout:   15,
				cancelTimeout: 30,
				writeInterval: 100,
				startCommand: [][]byte{
					{0x01, 0x09},
				},
				holdCommand: []byte{0x01, 0x09},
				rxFormat: []RxFormat{
					{StartByte: []byte{0x01}, EndByte: []byte{0x09}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if err := tt.device.Init(tt.args.config); err != nil {
				t.Error(err)
				return
			}
			log.Print("Initialized the device.")
			rxChannel := tt.device.GetRxChannel()
			txChannel := tt.device.GetTxChannel()
			cancelChannel := make(chan time.Time, 1)
			go func(c chan []byte) {
				for _, bytes := range tt.args.startCommand {
					c <- bytes
					time.Sleep(time.Duration(tt.args.writeInterval) * time.Millisecond)
				}
				if len(tt.args.holdCommand) > 0 {
					for {
						c <- tt.args.holdCommand
						time.Sleep(time.Duration(tt.args.writeInterval) * time.Millisecond)
					}
				}
			}(txChannel)
			go func(c chan time.Time) {
				log.Print("Cancel channel start timing.")
				c <- <-time.After(time.Duration(tt.args.cancelTimeout) * time.Second)
			}(cancelChannel)
			if len(tt.args.rxFormat) > 0 {
				log.Print("Set the receiving format.")
				tt.device.SetRxFormat(tt.args.rxFormat)
			}
			for {
				select {
				case <-time.After(time.Duration(tt.args.readTimeout) * time.Second):
					t.Error("Read channel timeout.")
					return
				case <-cancelChannel:
					log.Print("Cancel channel timeout.")
					return
				case rx := <-rxChannel:
					log.Printf("% x", rx)
				}
			}
		})
	}
}

func TestRxFormats(t *testing.T) {
	type args struct {
		rxFormat []RxFormat
		bytes    [][]byte
	}
	tests := []struct {
		device *Device
		args   args
	}{
		{
			device: NewDevice(),
			args: args{
				rxFormat: []RxFormat{
					{StartByte: []byte{0x01}, EndByte: []byte{0x09}},
					{StartByte: []byte{0x01}, LengthByteIndex: 1, LengthByteMissing: 7},
					{StartByte: []byte{0x01}, LengthFixed: 9},
				},
				bytes: [][]byte{
					{
						0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt.device.SetRxFormat(tt.args.rxFormat)
			for _, bytes := range tt.args.bytes {
				if got := tt.device.TestRxFormats(bytes); !reflect.DeepEqual(got, bytes) {
					t.Errorf("got = [% x], want [% x]", got, bytes)
				}
			}
		})
	}
}
