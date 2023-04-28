package main

import (
	"fmt"

	"github.com/egginabucket/openmsr/pkg/libmsr"
	"github.com/karalabe/usb"
)

func main() {
	hids, err := usb.EnumerateHid(libmsr.VendorID, libmsr.ProductID)
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
	var t [3][]byte
	var p []bool
	var lrcOK bool
	defer dev.Close() // d.Close() called laters
	d := libmsr.NewDevice(dev)
	err = d.Reset()
	if err != nil {
		panic(err)
	}

	s, err = d.FirmwareVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println("Firmare version:", s)

	s, err = d.Model()
	if err != nil {
		panic(err)
	}
	fmt.Println("MSR Model:", s)

	err = d.SetHiCo()
	if err != nil {
		panic(err)
	}
	/*
		fmt.Println("READ ISO ...")
		b, err = d.ReadISOTracks()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	*/

	fmt.Println("READ RAW ...")
	t, err = d.ReadRawTracks()
	if err != nil {
		panic(err)
	}
	fmt.Println("Raw data:", t[0], t[1], t[2])
	b, p, lrcOK = libmsr.DecodeRaw(t[0], ' ', 7, 8, false)
	fmt.Println("Track 1", string(b))
	fmt.Println(p)
	fmt.Println(lrcOK)
	b, p, _ = libmsr.DecodeRaw(t[1], '0', 5, 8, false)
	fmt.Println(p)
	fmt.Println("Track 2", string(b))
	b, _, _ = libmsr.DecodeRaw(t[2], '0', 5, 8, false)
	fmt.Println("Track 3", string(b))
	/*
		b, err = d.ReadISOTracks()
		if err != nil {
			panic(err)
		}

			fmt.Println("ISO data:", string(b))
	*/
	/*
		for _, m := range []msr.LEDMode{msr.LEDRedOn, msr.LEDYellowOn, msr.LEDGreenOn, msr.LEDAllOn} {
			err = d.SetLED(m)
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Second)
			d.SetLED(msr.LEDAllOff)
			time.Sleep(time.Second)
		}
	*/
	err = d.Close()
	if err != nil {
		panic(err)
	}
}
