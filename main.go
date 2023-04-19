package main

import (
	"fmt"
	"time"

	"github.com/egginabucket/openmsr/msr"
	"github.com/karalabe/usb"
)

func main() {
	hids, err := usb.EnumerateHid(msr.VendorID, msr.ProductID)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(hids), "hids")
	hid := hids[0]
	fmt.Println("Product:", hid.Product)
	fmt.Println("Manufacturer:", hid.Manufacturer)

	dev, err := hid.Open()
	if err != nil {
		panic(err)
	}
	var s string
	var b []byte
	defer dev.Close() // d.Close() called laters
	d := msr.NewDevice(dev)
	err = d.Reset()
	if err != nil {
		panic(err)
	}

	s, err = d.FirmwareVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("Firmare version:", s)

	s, err = d.MSRModel()
	if err != nil {
		panic(err)
	}
	fmt.Println("MSR Model:", s)

	err = d.SetHiCo()
	if err != nil {
		panic(err)
	}

	b, err = d.ReadISOTracks()
	if err != nil {
		panic(err)
	}
	fmt.Println("ISO data:", string(b))
	for _, m := range []msr.LEDMode{msr.LEDRedOn, msr.LEDYellowOn, msr.LEDGreenOn, msr.LEDAllOn} {
		err = d.SetLED(m)
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
		d.SetLED(msr.LEDAllOff)
		time.Sleep(1 * time.Second)
	}
	err = d.Close()
	if err != nil {
		panic(err)
	}
}
