package serialnetwork

import (
	"reflect"
	"testing"
)

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
