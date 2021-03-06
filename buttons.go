package levis

import (
	"fmt"
	"time"

	"github.com/dumacp/go-logs/pkg/logs"
)

type Button struct {
	Addr  int
	Value int
}

var buttons map[int]int

func (m *device) AddButton(addr int) error {

	if addr > m.conf.buttons_end || addr < m.conf.buttons_start {
		return fmt.Errorf("address in out of area, %d, area: %d <-> %d",
			addr, m.conf.buttons_start, m.conf.buttons_end)
	}

	if len(buttons) <= 0 {
		buttons = make(map[int]int)
	}
	buttons[addr] = 0x00

	return nil
}

func (m *device) ListenButtons() chan *Button {

	ch := make(chan *Button)

	go func() {

		defer close(ch)

		t1 := time.NewTicker(100 * time.Millisecond)
		defer t1.Stop()

		for {

			select {
			case <-m.quit:
				return
			case <-t1.C:

				var regs []byte
				if v, err := func() ([]byte, error) {

					m.mux.Lock()
					defer m.mux.Unlock()

					res, err := m.client.ReadCoils(
						uint16(m.conf.buttons_start), uint16(m.conf.buttons_end-m.conf.buttons_start+1))
					if err != nil {
						return nil, err
					}

					fmt.Printf("regs: %v\n", regs)
					return res, nil
				}(); err != nil {
					logs.LogError.Println(err)
					continue

				} else {
					regs = v
				}

				regsButtons := make([]int, 0)

				for i := range regs {

					for j := range make([]int, 8) {
						regsButtons = append(regsButtons, int((regs[i]>>(j))&0x01))
					}
				}

				fmt.Printf("regsButtons: %v\n", regsButtons)

				for k, v := range buttons {
					if v != regsButtons[k] {
						select {
						case ch <- &Button{Addr: k, Value: regsButtons[k]}:
						case <-time.After(100 * time.Millisecond):
						}

						buttons[k] = regsButtons[k]
					}
				}
			}
		}
	}()

	return ch
}

func (m *device) SetIndicator(addr int, value bool) error {

	m.mux.Lock()
	defer m.mux.Unlock()

	val := uint16(0x0000)
	if value {
		val = 0xFF00
	}
	if _, err := m.client.WriteSingleCoil(uint16(addr), val); err != nil {
		return err
	}

	return nil
}
