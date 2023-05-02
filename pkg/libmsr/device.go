package libmsr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/karalabe/usb"
)

type Device struct {
	device       usb.Device
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

// TestCommunication verifies the connection with the device.
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
		return errors.New("libmsr.Device.TestCommunication: unknown response")
	}
	return nil
}

// TestSensor verifies that the device's card sensing circuit is working.
// Does not return until a card is sensed or d.SwipeTimeout is reached.
func (d *Device) TestSensor() error {
	return d.sendAndCheck(esc(0x86), true)
}

// TestRAM verifies that the device's onboard RAM is working.
func (d *Device) TestRAM() error {
	return d.sendAndCheck(esc(0x87), false)
}

// SetLoCo sets the device to write Lo-Co cards.
func (d *Device) SetLoCo() error {
	return d.sendAndCheck(esc('x'), false)
}

// SetHiCo sets the device to write Hi-Co cards.
func (d *Device) SetHiCo() error {
	return d.sendAndCheck(esc('y'), false)
}

// IsHiCo checks the device's current write coercivity.
func (d *Device) IsHiCo() (bool, error) {
	msg, err := d.sendAndReceive(esc('d'), false)
	if err != nil {
		return false, err
	}
	switch msg[1] {
	case 'h':
		return true, nil
	case 'l':
		return false, nil
	}
	return false, fmt.Errorf("libmsr.Device.IsHiCo: unknown response %X", msg[1])
}

// SetBitsPerInch sets the density of each track in BPI.
func (d *Device) SetBitsPerInch(t1, t2, t3 int) error {
	invalidBPI := errors.New("libmsr.Device.SitBitsPerInch: invalid BPI")
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

// SetBitsPerChar sets the number of bits (including parity) for each track.
func (d *Device) SetBitsPerChar(t1, t2, t3 int) error {
	return d.sendAndCheck(esc('o', byte(t1), byte(t2), byte(t3)), false)
}

// Erase clears the selected tracks on a card.
// Does not return until a card is swiped or d.SwipeTimeout is reached.
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

// WriteRawTracks writes raw data to a card.
// Data can be encoded with EncodeRaw.
// Does not return until a card is swiped or d.SwipeTimeout is reached.
func (d *Device) WriteRawTracks(t1, t2, t3 []byte) error {
	return d.sendAndCheck(append(esc('n'), encodeRawTracks(t1, t2, t3)...), true)
}

// WriteISOTracks writes ISO data to a card.
// Does not return until a card is swiped or d.SwipeTimeout is reached.
func (d *Device) WriteISOTracks(t1, t2, t3 []byte) error {
	return d.sendAndCheck(append(esc('w'), encodeISOTracks(t1, t2, t3)...), true)
}

// ReadISOTracks reads ISO data from a card.
// Does not return until a card is swiped or d.SwipeTimeout is reached.
func (d *Device) ReadISOTracks() ([]byte, error) {
	return d.sendAndReceive(esc('r'), true)
}

// WriteRawTracks reads raw data from a card.
// Data can be decoded with DecodeRaw.
// Does not return until a card is swiped or d.SwipeTimeout is reached.
func (d *Device) ReadRawTracks() ([3][]byte, error) {
	data, _, err := d.sendAndReceiveEncoded(esc('m'), true)
	if err != nil {
		return [3][]byte{}, err
	}
	return decodeTracks(data)
}

// Model returns the device's reported model.
func (d *Device) Model() (string, error) {
	msg, err := d.sendAndReceive(esc('t'), false)
	return string(msg), err
}

// FirmwareVersion returns the device's reported firmware version.
func (d *Device) FirmwareVersion() (string, error) {
	msg, err := d.sendAndReceive(esc('v'), false)
	return string(msg), err
}

// SetLED sets the device's LEDs to mode.
func (d *Device) SetLED(mode LEDMode) error {
	if mode > LEDRedOn {
		return errors.New("libmsr.Device.SetLED: invalid LED mode")
	}
	return d.send(esc(0x81 + byte(mode)))
}

// Reset resets the device.
// Useful for cancelling timed-out operations.
func (d *Device) Reset() error {
	return d.send(esc('a'))
}

// Close resets and closes the USB device.
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
