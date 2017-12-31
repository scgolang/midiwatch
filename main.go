package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/scgolang/midi"
)

func main() {
	var (
		name string
	)
	flag.StringVar(&name, "d", "", "MIDI device name prefix.")
	flag.Parse()

	fmt.Println("Finding device...")
	device, err := findDevice(name)
	if err != nil {
		Die(err)
	}
	fmt.Println("Found device.")

	pkts, err := device.Packets()
	if err != nil {
		Die(err)
	}
	fmt.Println("Waiting for MIDI packets...")

	var i int
	for pkt := range pkts {
		if pkt.Err != nil {
			Die(pkt.Err)
		}
		fmt.Printf("%-16d%#v\n", i, pkt.Data)
		i++
	}
}

func Die(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func findDevice(prefix string) (*midi.Device, error) {
	devices, err := midi.Devices()
	if err != nil {
		return nil, errors.Wrap(err, "listing devices")
	}
	var d *midi.Device
	for _, device := range devices {
		if strings.HasPrefix(device.Name, prefix) {
			d = device
		}
	}
	if d == nil {
		return nil, errors.New("could not find device with prefix " + prefix)
	}
	return d, errors.Wrap(d.Open(), "opening device "+d.Name)
}
