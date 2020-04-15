package main

// This example application executes the following procedure:
// 1. Scans for a peripheral with the specified name
// 2. Connects to the peripheral
// 3. Performs service discovery
// 4. Reads a characteristic (if a suitable one is found)
// 5. Reads a decriptor (if a suitable one is found)

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/JuulLabs-OSS/cbgo"
)

type MyDelegate struct {
	cbgo.CentralManagerDelegateBase
	cbgo.PeripheralDelegateBase
}

var devName string
var myPrph cbgo.Peripheral

// This channel is used to make a blocking interface out of the CoreBluetooth
// API.
var ch = make(chan error)

func block() {
	err := <-ch
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
}

func (d *MyDelegate) DidUpdateState(cmgr cbgo.CentralManager) {
	if cmgr.State() != cbgo.ManagerStatePoweredOn {
		ch <- fmt.Errorf("central manager has invalid state: have=%d want=%d: is Bluetooth turned on?",
			cmgr.State(), cbgo.ManagerStatePoweredOn)
	} else {
		// Bluetooth is powered on.  Unblock the main thread.
		ch <- nil
	}
}

func (d *MyDelegate) DidDiscoverPeripheral(cm cbgo.CentralManager, prph cbgo.Peripheral,
	advFields cbgo.AdvFields, rssi int) {

	name := advFields.LocalName
	if name == "" {
		name = "<none>"
	}
	fmt.Printf("discovered peripheral: name=%s\n", name)

	if advFields.LocalName == devName {
		// We found our peer.  Abort the scan and unblock the main thread.
		myPrph = prph
		cm.StopScan()
		ch <- nil
	}
}

func (d *MyDelegate) DidConnectPeripheral(cm cbgo.CentralManager, prph cbgo.Peripheral) {
	// Make sure our delegate gets called for events related to this
	// peripheral.
	prph.SetDelegate(d)

	// Unblock the main thread now that we're connected.
	ch <- nil
}

func (d *MyDelegate) DidFailToConnectPeripheral(cm cbgo.CentralManager, prph cbgo.Peripheral, err error) {
	ch <- fmt.Errorf("failed to connect: %v", err)
}

func (d *MyDelegate) DidDisconnectPeripheral(cm cbgo.CentralManager, prph cbgo.Peripheral, err error) {
	ch <- fmt.Errorf("peripheral disconnected: %v", err)
}

func (d *MyDelegate) DidDiscoverServices(prph cbgo.Peripheral, err error) {
	if err != nil {
		ch <- fmt.Errorf("failed to discover services: %v\n", err)
	} else {
		// Unblock the main thread now that we have discovered the peer's
		// services.
		ch <- nil
	}
}

func (d *MyDelegate) DidDiscoverCharacteristics(prph cbgo.Peripheral, svc cbgo.Service, err error) {
	if err != nil {
		ch <- fmt.Errorf("failed to discover characteristics: %v\n", err)
	} else {
		// Unblock the main thread now that we have finished discovering
		// characteristics for the current service.
		ch <- nil
	}
}

func (d *MyDelegate) DidDiscoverDescriptors(prph cbgo.Peripheral, chr cbgo.Characteristic, err error) {
	if err != nil {
		ch <- fmt.Errorf("failed to discover descriptors: %v\n", err)
	} else {
		// Unblock the main thread now that we have finished discovering
		// descriptors for the current characteristic.
		ch <- nil
	}
}

func (d *MyDelegate) DidUpdateValueForCharacteristic(prph cbgo.Peripheral, chr cbgo.Characteristic, err error) {
	if err != nil {
		ch <- fmt.Errorf("failed to read characteristic: %v", err)
	} else {
		// Unblock the main thread now that we have received the read response.
		ch <- nil
	}
}

func (d *MyDelegate) DidUpdateValueForDescriptor(prph cbgo.Peripheral, dsc cbgo.Descriptor, err error) {
	if err != nil {
		ch <- fmt.Errorf("failed to read descriptor: %v", err)
	} else {
		// Unblock the main thread now that we have received the read response.
		ch <- nil
	}
}

// discoverProfile performs full service discovery on the specified peripheral.
func discoverProfile(prph cbgo.Peripheral) {
	prph.DiscoverServices(nil)
	block()

	svcs := prph.Services()
	for _, s := range svcs {
		prph.DiscoverCharacteristics(nil, s)
		block()

		chrs := s.Characteristics()
		for _, c := range chrs {
			prph.DiscoverDescriptors(c)
			block()
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <device-name>\n", os.Args[0])
		os.Exit(1)
	}
	devName = os.Args[1]

	cm := cbgo.NewCentralManager(nil)
	cm.SetDelegate(&MyDelegate{})

	// Wait for the Bluetooth power on event.
	block()

	cm.Scan(nil, nil)
	block()

	cm.Connect(myPrph, nil)
	block()

	fmt.Printf("Connected to %s\n", myPrph.Name())

	discoverProfile(myPrph)

	var chr *cbgo.Characteristic
	var dsc *cbgo.Descriptor

	// Print the peer's profile.
	svcs := myPrph.Services()
	for i, s := range svcs {
		fmt.Printf("s   %d: %v\n", i, s.UUID())

		chrs := s.Characteristics()
		for j, c := range chrs {
			if chr == nil && c.Properties()&cbgo.CharacteristicPropertyRead != 0 {
				// This characteristic is readable.  Remember it so that we can
				// read it later.
				chr = &chrs[j]
			}
			fmt.Printf("c       %d: %v\n", j, c.UUID())

			dscs := c.Descriptors()
			for k, d := range dscs {
				if dsc == nil {
					dsc = &dscs[k]
				}
				fmt.Printf("d           %d: %v\n", k, d.UUID())
			}
		}
	}

	if chr == nil {
		fmt.Printf("no characteristics to read!\n")
	} else {
		myPrph.ReadCharacteristic(*chr)
		block()

		fmt.Printf("read characteristic:\n")
		fmt.Printf("    UUID: %v\n", chr.UUID())
		fmt.Printf("    value: %v\n", hex.EncodeToString(chr.Value()))
	}

	if dsc == nil {
		fmt.Printf("no descriptors to read!\n")
	} else {
		myPrph.ReadDescriptor(*dsc)
		block()

		fmt.Printf("read descriptor:\n")
		fmt.Printf("    UUID: %v\n", dsc.UUID())
		fmt.Printf("    value: %v\n", hex.EncodeToString(dsc.Value()))
	}
}
