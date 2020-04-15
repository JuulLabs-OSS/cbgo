package cbgo

/*
// See cutil.go for C compiler flags.
#import "bt.h"
*/
import "C"

// CentralManagerScanOpts: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/peripheral_scanning_options
type CentralManagerScanOpts struct {
	AllowDuplicates       bool
	SolicitedServiceUUIDs []UUID
}

// DfltCentralManagerScanOpts is the set of options that gets used when nil is
// passed to `Scan()`.
var DfltCentralManagerScanOpts = CentralManagerScanOpts{}

// Scan: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1518986-scanforperipheralswithservices
func (cm CentralManager) Scan(serviceUUIDs []UUID, opts *CentralManagerScanOpts) {
	arrSvcUUIDs := uuidsToStrArr(serviceUUIDs)
	defer freeStrArr(&arrSvcUUIDs)

	if opts == nil {
		opts = &DfltCentralManagerScanOpts
	}

	arrSolSvcUUIDs := uuidsToStrArr(opts.SolicitedServiceUUIDs)
	defer freeStrArr(&arrSolSvcUUIDs)

	copts := C.struct_scan_opts{
		allow_dups:    C.bool(opts.AllowDuplicates),
		sol_svc_uuids: arrSolSvcUUIDs,
	}

	C.cb_cmgr_scan(cm.ptr, &arrSvcUUIDs, &copts)
}

// StopScan: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1518984-stopscan
func (cm CentralManager) StopScan() {
	C.cb_cmgr_stop_scan(cm.ptr)
}

// IsScanning: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1620640-isscanning
func (cm CentralManager) IsScanning() bool {
	return bool(C.cb_cmgr_is_scanning(cm.ptr))
}
