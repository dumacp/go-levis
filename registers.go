package levis

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"
)

type Register struct {
	Addr  int
	Value []uint16
}

var inputs map[int][]uint16

func (m *device) AddInput(addr, length int) error {

	if addr < m.conf.read_start || addr+length > m.conf.read_end {
		return fmt.Errorf("address out of area, %d (len: %d), area: %d <-> %d",
			addr, length, m.conf.read_start, m.conf.read_end)
	}

	if len(inputs) <= 0 {
		inputs = make(map[int][]uint16)
	}
	inputs[addr] = make([]uint16, length)

	return nil
}

func (m *device) ListenInputsWithContext(ctx context.Context) chan *Register {

	ch := make(chan *Register)

	go func() {

		if ctx == nil {
			var cancel func()
			ctx, cancel = context.WithCancel(context.TODO())
			defer cancel()
		}

		defer close(ch)
		defer fmt.Println("stop listenInputs")

		t1 := time.NewTicker(300 * time.Millisecond)
		defer t1.Stop()

		for {

			select {
			case <-ctx.Done():
				return
			case <-m.contxt.Done():
				return
			case <-t1.C:
				var regs []byte
				if err := func() error {
					m.mux.Lock()
					defer m.mux.Unlock()
					if v, err := m.client.ReadInputRegisters(
						uint16(m.conf.read_start),
						uint16(m.conf.read_end-m.conf.read_start+1)); err != nil {
						return err
					} else {
						regs = v
					}
					return nil
				}(); err != nil {
					fmt.Printf("error listenInputs = %s", err)
					continue
				}

				// fmt.Printf("regs inputs: %v\n", regs)

				regsButtons := make([]uint16, 0)

				for i := range make([]int, len(regs)/2) {

					idx := 2 * i
					value := []byte{regs[idx], regs[idx+1]}
					regsButtons = append(regsButtons, binary.BigEndian.Uint16(value))
				}

				// fmt.Printf("regs uint16 inputs: %v\n", regsButtons)

				for k, v := range inputs {
					if !Equal(v, regsButtons[k:k+len(v)]) {
						select {
						case ch <- &Register{Addr: k, Value: regsButtons[k:len(v)]}:
						case <-time.After(1 * time.Second):
						}
						value := make([]uint16, len(v))
						copy(value, regsButtons[k:len(v)])
						inputs[k] = value
					}
				}
			}
		}

	}()

	return ch
}

func (m *device) ListenInputs() chan *Register {

	return m.ListenInputsWithContext(nil)
}

func (m *device) WriteRegister(addr int, value []uint16) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	// fmt.Printf("WriteRegister (addr: %d, len: %d): [%X]\n", addr, len(value), value)
	valueBytes := DecodeToBytes(value)

	data := make([][]byte, 0)

	n := len(valueBytes) / 120
	for i := range make([]int, n+1) {
		if len(valueBytes) > (i+1)*120 {
			data = append(data, valueBytes[i*120:120*(i+1)])
		} else {
			data = append(data, valueBytes[i*120:])
			break
		}
	}

	for i, v := range data {
		fmt.Printf("data: %X\n", v)
		_, err := m.client.WriteMultipleRegisters(
			uint16(addr+(i*120/2)), uint16(len(v)/2), v)
		if err != nil {
			return err
		}
	}
	fmt.Printf("WriteRegister (addr: %d, len: %d): [%X]\n", addr, len(value), value)

	return nil
}

func (m *device) ReadRegister(addr, length int) ([]uint16, error) {

	m.mux.Lock()
	defer m.mux.Unlock()
	result, err := m.client.ReadHoldingRegisters(
		uint16(addr), uint16(length))
	if err != nil {
		return nil, err
	}

	return EncodeFromBytes(result), nil
}

func (m *device) ReadBytesRegister(addr, length int) ([]byte, error) {

	m.mux.Lock()
	defer m.mux.Unlock()
	// fmt.Println("voya aqui")
	result, err := m.client.ReadHoldingRegisters(
		uint16(addr), uint16(length))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *device) WriteRawRegister(addr int, value []byte) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	if len(value)%2 != 0 {
		valueCopy = append(valueCopy, 0x00)
	}

	_, err := m.client.WriteMultipleRegisters(
		uint16(addr), uint16(len(valueCopy)/2), valueCopy)
	if err != nil {
		return err
	}

	return nil
}
