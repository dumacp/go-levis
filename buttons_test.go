package levis

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goburrow/modbus"
)

func Test_device_ListenButtons(t *testing.T) {

	config := &conf{
		buttons_start: 0,
		buttons_end:   4,
		read_start:    0,
		read_end:      10,
	}

	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	if err := handler.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer handler.Close()

	client := modbus.NewClient(handler)

	type fields struct {
		conf *conf
		dev  modbus.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   chan *Button
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			fields: fields{
				conf: config,
				dev:  client,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &device{
				conf:   tt.fields.conf,
				client: tt.fields.dev,
			}

			if err := m.AddButton(1); err != nil {
				t.Error(err)
			}
			if err := m.AddButton(2); err != nil {
				t.Error(err)
			}
			if err := m.AddButton(3); err != nil {
				t.Error(err)
			}
			if err := m.AddButton(4); err != nil {
				t.Error(err)
			}

			if got := m.ListenButtons(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("device.ListenButtons() = %v, don't want %v", got, tt.want)
				return
			} else {
				for v := range got {
					t.Logf("response: %v", v)
				}
			}

		})
	}
}

func Test_device_SetIndicator(t *testing.T) {

	config := &conf{
		buttons_start: 0,
		buttons_end:   2,
		read_start:    0,
		read_end:      10,
	}

	handler := modbus.NewRTUClientHandler("/dev/ttyUSB0")
	handler.BaudRate = 19200
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second

	if err := handler.Connect(); err != nil {
		log.Fatalln(err)
	}

	defer handler.Close()

	client := modbus.NewClient(handler)

	type fields struct {
		conf *conf
		dev  modbus.Client
	}
	type args struct {
		addr  int
		value bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  11,
				value: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &device{
				conf:   tt.fields.conf,
				client: tt.fields.dev,
			}
			if err := m.SetIndicator(tt.args.addr, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("device.SetIndicator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
