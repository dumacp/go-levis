package levis

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/goburrow/modbus"
)

func Test_device_ListenInputs(t *testing.T) {

	config := &conf{
		buttons_start: 0,
		buttons_end:   2,
		read_start:    0,
		read_end:      40,
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
		want   chan *Register
	}{
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

			if err := m.AddInput(5, 10); err != nil {
				t.Error(err)
			}
			if got := m.ListenInputs(); reflect.DeepEqual(got, tt.want) {
				t.Errorf("device.ListenInputs() = %v, want %v", got, tt.want)
			} else {
				for v := range got {
					t.Logf("response: %v, text: %s, len: %v",
						v, DecodeToChars(v.Value), len(DecodeToChars(v.Value)))
				}
			}
		})
	}
}

func Test_device_WriteRawRegister(t *testing.T) {

	config := &conf{
		buttons_start: 0,
		buttons_end:   2,
		read_start:    0,
		read_end:      40,
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
		value []byte
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
				addr:  2,
				value: []byte{0x00, 0x16, 0x00, 0x04},
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
			if err := m.WriteRawRegister(tt.args.addr, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("device.WriteRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_device_WriteRegister(t *testing.T) {

	config := &conf{
		buttons_start: 0,
		buttons_end:   2,
		read_start:    0,
		read_end:      40,
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
		value []uint16
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test0",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  0,
				value: []uint16{0},
			},
			wantErr: false,
		},
		{
			name: "test1",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  2,
				value: []uint16{31, 66},
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  20,
				value: EncodeFromChars([]byte("RUTA CARAJILLO")),
			},
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  10,
				value: EncodeFromChars([]byte("02:11")),
			},
			wantErr: false,
		},
		{
			name: "test4",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  4,
				value: []uint16{70},
			},
			wantErr: false,
		},
		{
			name: "test5",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  200,
				value: EncodeFromChars([]byte("RUTA ORIENTAL")),
			},
			wantErr: false,
		},
		{
			name: "test6",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:  0,
				value: []uint16{1},
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
			if err := m.WriteRegister(tt.args.addr, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("device.WriteRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_device_ReadRegister(t *testing.T) {

	config := &conf{
		buttons_start: 0,
		buttons_end:   2,
		read_start:    0,
		read_end:      40,
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
		addr   int
		length int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:   20,
				length: 32,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:   110,
				length: 1,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "test3",
			fields: fields{
				conf: config,
				dev:  client,
			},
			args: args{
				addr:   110,
				length: 1,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &device{
				conf:   tt.fields.conf,
				client: tt.fields.dev,
			}
			got, err := m.ReadBytesRegister(tt.args.addr, tt.args.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("device.ReadRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("device.ReadRegister() = %s, want %v", got, tt.want)
			}
		})
	}
}
