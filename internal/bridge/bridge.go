package bridge

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const (
	usbVID = "239A"
	usbPID = "800C"
)

type Bridge interface {
	io.Closer

	Init() error
	OutChan() <-chan string
	StartScanning()
}

type bridge struct {
	ctx    context.Context
	cancel context.CancelFunc
	sd     serial.Port

	chOut chan string
}

func New(ctx context.Context) (Bridge, error) {
	b := &bridge{
		chOut: make(chan string, 64),
	}
	b.ctx, b.cancel = context.WithCancel(ctx)
	err := b.Init()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *bridge) Close() error {
	b.cancel()
	time.Sleep(200 * time.Millisecond)
	err := b.sd.Close()
	close(b.chOut)
	return err
}

func (b *bridge) Init() error {
	dpl, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return err
	}

	mode := serial.Mode{
		BaudRate:          9600,
		DataBits:          8,
		Parity:            serial.NoParity,
		StopBits:          serial.OneStopBit,
		InitialStatusBits: nil,
	}

	var path string

	for _, d := range dpl {
		if d.IsUSB && strings.ToUpper(d.VID) == usbVID && strings.ToUpper(d.PID) == usbPID {
			path = d.Name
			break
		}
	}

	if path == "" {
		return fmt.Errorf("no lora bridge found")
	}

	b.sd, err = serial.Open(path, &mode)
	if err != nil {
		return err
	}

	return nil
}

func (b *bridge) OutChan() <-chan string {
	return b.chOut
}

func (b *bridge) StartScanning() {
	scanner := bufio.NewScanner(b.sd)
	scanner.Split(bufio.ScanLines)

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
			if scanner.Scan() {
				b.chOut <- scanner.Text()
			}
		}
	}
}
