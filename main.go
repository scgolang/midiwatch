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
		list bool
		name string
	)
	flag.BoolVar(&list, "l", false, "List MIDI devices.")
	flag.StringVar(&name, "d", "", "MIDI device name prefix.")
	flag.Parse()

	if list {
		Die(listDevices())
		return
	}
	fmt.Println("Finding device...")
	device, err := findDevice(name)
	if err != nil {
		Die(err)
	}
	fmt.Println("Found device.")

	pktslice, err := device.Packets()
	if err != nil {
		Die(err)
	}
	fmt.Println("Waiting for MIDI packets...")

	var i int
	for pkts := range pktslice {
		for _, pkt := range pkts {
			if pkt.Err != nil {
				Die(pkt.Err)
			}
			fmt.Printf("%-16d%#v\n", i, pkt.Data)
			i++
		}
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

func listDevices() error {
	devices, err := midi.Devices()
	if err != nil {
		return errors.Wrap(err, "getting devices list")
	}
	fmt.Printf("ID: (TYPE) NAME\n")
	for _, d := range devices {
		fmt.Printf("%s: (%s) \"%s\"\n", d.ID, d.Type, d.Name)
	}
	return nil
}
