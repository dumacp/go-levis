package levis

import (
	"sync"
	"time"

	"github.com/goburrow/modbus"
)

type device struct {
	conf    *conf
	client  modbus.Client
	mux     sync.Mutex
	handler *modbus.RTUClientHandler
	quit    chan int
}

func NewDevice(port string, speedBaud int) (Device, error) {
	return NewDeviceWithID(port, speedBaud, 1)

}

func NewDeviceWithID(port string, speedBaud int, id int) (Device, error) {
	dev := &device{
		conf: &conf{},
		mux:  sync.Mutex{},
	}

	handler := modbus.NewRTUClientHandler(port)
	handler.BaudRate = speedBaud
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = byte(id)
	handler.Timeout = 600 * time.Millisecond
	// handler.IdleTimeout = 10 * time.Millisecond

	if err := handler.Connect(); err != nil {
		return nil, err
	}

	if dev.quit != nil {
		select {
		case <-dev.quit:
		default:
			close(dev.quit)
			time.Sleep(10 * time.Millisecond)
		}
	}

	dev.quit = make(chan int)

	dev.handler = handler

	client := modbus.NewClient(handler)

	dev.client = client

	// dev.quit = make(chan int)

	return dev, nil

}

func (dev *device) Close() error {
	if dev.quit != nil {
		select {
		case _, ok := <-dev.quit:
			if ok {
				close(dev.quit)
			}
		default:
			close(dev.quit)
		}
		// time.Sleep(600 * time.Millisecond)
	}
	return dev.handler.Close()
}

func (dev *device) SetSlaveID(id int) {
	dev.handler.SlaveId = byte(id)
}

func (dev *device) Conf() Conf {
	return dev.conf
}

func (dev *device) ReadTimeout() time.Duration {
	return dev.handler.Timeout
}

type Device interface {
	SetSlaveID(id int)
	ReadTimeout() time.Duration
	ListenButtons() chan *Button
	ListenInputs() chan *Register
	WriteRegister(addr int, value []uint16) error
	ReadRegister(addr, length int) ([]uint16, error)
	ReadBytesRegister(addr, length int) ([]byte, error)
	AddButton(addr int) error
	AddInput(addr, length int) error
	Conf() Conf
	SetIndicator(addr int, value bool) error
	Close() error
}
