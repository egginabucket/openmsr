package libmsr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/karalabe/usb"
)

type Device struct {
	device usb.Device
	//bpc          [3]int
	PreSendDelay time.Duration
	CheckTimeout,
	SwipeTimeout time.Duration
}

type LEDMode byte

const (
	LEDAllOff LEDMode = iota
	LEDAllOn
	LEDGreenOn
	LEDYellowOn
	LEDRedOn
)

const (
	VendorID  uint16 = 0x0801
	ProductID uint16 = 0x0003
)

func (d *Device) send(msg []byte) error {
	time.Sleep(d.PreSendDelay)
	for _, pkt := range makePackets(msg) {
		_, err := d.device.Write(pkt) // HID null byte handled in karalabe/usb
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Device) receivePacket(pktChan chan []byte, errChan chan error) {
	pkt := make([]byte, 64)
	_, err := d.device.Read(pkt)
	if err != nil {
		errChan <- err
		return
	}
	pktChan <- pkt
}

func (d *Device) receive(swipeWait bool) ([]byte, error) {
	var timeout time.Duration
	if swipeWait {
		timeout = d.SwipeTimeout
	} else {
		timeout = d.CheckTimeout
	}
	pkts := make([][]byte, 0)
packets:
	for {
		pktChan := make(chan []byte, 1)
		errChan := make(chan error, 1)
		ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		go d.receivePacket(pktChan, errChan)
		select {
		case <-ctxTimeout.Done():
			return nil, ctxTimeout.Err()
		case pkt := <-pktChan:
			pkts = append(pkts, pkt)
			if pkt[0]&seqEndBit == seqEndBit {
				break packets
			}
		case err := <-errChan:
			return nil, err
		}

	}
	return parsePackets(pkts)
}

func (d *Device) receiveEncoded(swipeWait bool) (data, result []byte, err error) {
	var msg []byte
	msg, err = d.receive(swipeWait)
	if err != nil {
		return
	}
	return decode(msg)
}

func (d *Device) sendAndReceive(msg []byte, swipeWait bool) ([]byte, error) {
	err := d.send(msg)
	if err != nil {
		return nil, err
	}
	return d.receive(swipeWait)
}

func (d *Device) sendAndReceiveEncoded(msg []byte, swipeWait bool) (data, result []byte, err error) {
	err = d.send(msg)
	if err != nil {
		return
	}
	return d.receiveEncoded(swipeWait)
}

func (d *Device) sendAndCheck(msg []byte, swipeWait bool) error {
	_, _, err := d.sendAndReceiveEncoded(msg, swipeWait)
	return err
}

func (d *Device) TestCommunication() error {
	err := d.send(esc('e'))
	if err != nil {
		return err
	}
	msg, err := d.receive(false)
	if err != nil {
		return err
	}
	if msg[0] != escByte || msg[1] != 'y' {
		return errors.New("unknown response")
	}
	return nil
}

func (d *Device) TestSensor() error {
	return d.sendAndCheck(esc(0x86), true)
}

func (d *Device) TestRAM() error {
	return d.sendAndCheck(esc(0x87), false)
}

func (d *Device) SetLoCo() error {
	return d.sendAndCheck(esc('x'), false)
}

func (d *Device) SetHiCo() error {
	return d.sendAndCheck(esc('y'), false)
}

// TODO
func (d *Device) IsHiCo() (bool, error) {
	msg, err := d.sendAndReceive(esc('d'), false)
	if err != nil {
		return false, err
	}
	switch msg[1] {
	case 'H', 'h':
		return true, nil
	case 'L', 'l':
		return false, nil
	}
	fmt.Println(msg)
	return false, fmt.Errorf("unknown response %X", msg[1])
}

func (d *Device) SetBitsPerInch(t1, t2, t3 int) error {
	invalidBPI := errors.New("invalid BPI")
	cmd := esc('b')
	switch t1 {
	case 0:
	case 75, 210:
		cmd = append(cmd, byte(t1))
	default:
		return invalidBPI
	}
	switch t2 {
	case 0:
	case 75:
		cmd = append(cmd, 0xA0)
	case 210:
		cmd = append(cmd, 0xA1)
	default:
		return invalidBPI
	}
	switch t3 {
	case 0:
	case 75:
		cmd = append(cmd, 0xC0)
	case 210:
		cmd = append(cmd, 0xC1)
	default:
		return invalidBPI
	}
	return d.sendAndCheck(cmd, false)
}

func (d *Device) SetBitsPerChar(t1, t2, t3 int) error {
	return d.sendAndCheck(esc('o', byte(t1), byte(t2), byte(t3)), false)
}

func (d *Device) Erase(t1, t2, t3 bool) error {
	var mask byte
	if t1 {
		mask |= 1
	}
	if t2 {
		mask |= 1 << 1
	}
	if t3 {
		mask |= 1 << 2
	}
	return d.sendAndCheck(esc('c', mask), true)
}

func (d *Device) WriteRawTracks(t1, t2, t3 []byte) error {
	return d.sendAndCheck(append(esc('n'), encodeRawTracks(t1, t2, t3)...), true)
}

func (d *Device) WriteISOTracks(t1, t2, t3 []byte) error {
	return d.sendAndCheck(append(esc('w'), encodeISOTracks(t1, t2, t3)...), true)
}

func (d *Device) ReadISOTracks() ([]byte, error) {
	return d.sendAndReceive(esc('r'), true)
}

func (d *Device) ReadRawTracks() ([3][]byte, error) {
	data, _, err := d.sendAndReceiveEncoded(esc('m'), true)
	if err != nil {
		return [3][]byte{}, err
	}
	return decodeTracks(data)
}

func (d *Device) Model() (string, error) {
	msg, err := d.sendAndReceive(esc('t'), false)
	return string(msg), err
}

func (d *Device) FirmwareVersion() (string, error) {
	msg, err := d.sendAndReceive(esc('v'), false)
	return string(msg), err
}

func (d *Device) SetLED(mode LEDMode) error {
	if mode > LEDRedOn {
		return errors.New("invalid LED mode")
	}
	return d.send(esc(0x81 + byte(mode)))
}

func (d *Device) Reset() error {
	return d.send(esc('a'))
}

func (d *Device) Close() error {
	err := d.Reset()
	if err != nil {
		return err
	}
	return d.device.Close()
}

func NewDevice(d usb.Device) *Device {
	return &Device{
		device:       d,
		PreSendDelay: 10 * time.Millisecond,
		CheckTimeout: 150 * time.Millisecond,
		SwipeTimeout: 30 * time.Second,
	}
}
