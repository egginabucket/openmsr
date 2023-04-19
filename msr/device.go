package msr

import (
	"context"
	"errors"
	"time"

	"github.com/karalabe/usb"
)

type Device struct {
	device usb.Device
	//info   usb.DeviceInfo
	//mu sync.Mutex
}

const (
	VendorID  uint16 = 0x0801
	ProductID uint16 = 0x0003
)

var (
	PreSendDelay = 10 * time.Millisecond
	CheckTimeout = 150 * time.Millisecond
	ReadTimeout  = 30 * time.Minute
)

func (d *Device) send(msg []byte) error {
	time.Sleep(PreSendDelay)
	for _, pkt := range makePackets(msg) {
		//d.mu.Lock()
		_, err := d.device.Write(pkt) // hid null byte handled in karalabe/usb
		//d.mu.Unlock()
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

func (d *Device) receive(timeout time.Duration) ([]byte, error) {
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

func (d *Device) receiveEncoded(timeout time.Duration) (data, result []byte, err error) {
	var msg []byte
	msg, err = d.receive(timeout)
	if err != nil {
		return
	}
	return decode(msg)
}

/*
func (d *Device) sendAndWaitCustom(msg []byte, timeout time.Duration) error {
	err := d.send(msg)
	if err != nil {
		return err
	}
	_, _, err = d.receiveMessage(timeout)
	return err
}
*/

func (d *Device) sendAndCheck(msg []byte) error {
	err := d.send(msg)
	if err != nil {
		return err
	}
	_, _, err = d.receiveEncoded(CheckTimeout)
	return err
}

func (d *Device) sendAndReceive(msg []byte, timeout time.Duration) ([]byte, error) {
	err := d.send(msg)
	if err != nil {
		return nil, err
	}
	return d.receive(timeout)
}

func (d *Device) Reset() error {
	return d.send(esc('a'))
}

func (d *Device) SetLoCo() error {
	return d.sendAndCheck(esc('x'))
}

func (d *Device) SetHiCo() error {
	return d.sendAndCheck(esc('y'))
}

func (d *Device) SetBitsPerInch(track int, bpi int) error {
	if bpi != 210 && bpi != 75 {
		return errors.New("invalid BPI")
	}
	if track < 1 || track > 3 {
		return errors.New("invalid track")
	}
	var bpiByte byte
	switch track {
	case 1:
		bpiByte = byte(bpi)
	case 2:
		bpiByte = 0xA0
	case 3:
		bpiByte = 0xC0
	}
	if track != 1 && bpi == 210 {
		bpiByte++
	}
	return d.sendAndCheck(esc('b', bpiByte))
}

func (d *Device) SetBitsPerChar(t1, t2, t3 int) error {
	return d.sendAndCheck(esc('o', byte(t1), byte(t2), byte(t3)))
}

func (d *Device) Erase(t1, t2, t3 bool) error {
	var mask byte
	if t1 {
		mask |= 1 << 0
	}
	if t2 {
		mask |= 1 << 1
	}
	if t3 {
		mask |= 1 << 2
	}
	return d.sendAndCheck(esc('c', mask))
}

func (d *Device) WriteRawTracks(t1, t2, t3 []byte) error {
	return d.sendAndCheck(append(esc('n'), encodeRaw(t1, t2, t3)...))
}

func (d *Device) WriteISOTracks(t1, t2, t3 []byte) error {
	return d.sendAndCheck(append(esc('w'), encodeISO(t1, t2, t3)...))
}

func (d *Device) ReadISOTracks() ([]byte, error) {
	return d.sendAndReceive(esc('r'), ReadTimeout)
}

func (d *Device) ReadRawTracks() ([]byte, error) {
	return d.sendAndReceive(esc('m'), ReadTimeout)
}

func (d *Device) FirmwareVersion() (string, error) {
	msg, err := d.sendAndReceive(esc('v'), CheckTimeout)
	return string(msg), err
}

func (d *Device) MSRModel() (string, error) {
	msg, err := d.sendAndReceive(esc('t'), CheckTimeout)
	return string(msg), err
}

func (d *Device) SetLED(mode LEDMode) error {
	if !mode.ok() {
		return ErrInvalidLEDMode
	}
	return d.send(esc(byte(mode)))
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
		device: d,
	}
}
