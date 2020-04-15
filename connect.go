package cbgo

/*
// See cutil.go for C compiler flags.
#import "bt.h"
*/
import "C"

// CentralManagerConnectOpts: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/peripheral_connection_options
type CentralManagerConnectOpts struct {
	NotifyOnConnection      bool
	NotifyOnDisconnection   bool
	NotifyOnNotification    bool
	EnableTransportBridging bool
	RequiresANCS            bool
	StartDelay              int
}

// DfltCentralManagerConnectOpts is the set of options that gets used when nil
// is passed to `Connect()`.
var DfltCentralManagerConnectOpts = CentralManagerConnectOpts{
	NotifyOnConnection:    true,
	NotifyOnDisconnection: true,
	NotifyOnNotification:  true,
}

// Connect: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1518766-connectperipheral
func (cm CentralManager) Connect(prph Peripheral, opts *CentralManagerConnectOpts) {
	if opts == nil {
		opts = &DfltCentralManagerConnectOpts
	}

	copts := C.struct_connect_opts{
		notify_on_connection:      C.bool(opts.NotifyOnConnection),
		notify_on_disconnection:   C.bool(opts.NotifyOnDisconnection),
		notify_on_notification:    C.bool(opts.NotifyOnNotification),
		enable_transport_bridging: C.bool(opts.EnableTransportBridging),
		requires_ancs:             C.bool(opts.RequiresANCS),
		start_delay:               C.int(opts.StartDelay),
	}

	C.cb_cmgr_connect(cm.ptr, prph.ptr, &copts)
}

// CancelConnect: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1518952-cancelperipheralconnection
func (cm CentralManager) CancelConnect(prph Peripheral) {
	C.cb_cmgr_cancel_connect(cm.ptr, prph.ptr)
}
