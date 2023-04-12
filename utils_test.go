package levis

import (
	"reflect"
	"testing"

	"github.com/goburrow/modbus"
)

func TestGenerateMessage(t *testing.T) {
	type args struct {
		slaveId,
		funcCode int
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     make([]byte, 22),
			},
			want: nil,
		},
		{
			name: "test2",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "test3",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "test4",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "test5",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "test6",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "test7",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
		{
			name: "test8",
			args: args{
				slaveId:  0x01,
				funcCode: modbus.FuncCodeReadHoldingRegisters,
				data:     []byte{0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateMessage(tt.args.slaveId, tt.args.funcCode, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateMessage() = %v, %X, %q, want %v", got, got, got, tt.want)
			}
		})
	}
}
