package cbgo

/*
// See cutil.go for C compiler flags.
#import "bt.h"
*/
import "C"

import (
	"unsafe"
)

// ManagerState: https://developer.apple.com/documentation/corebluetooth/cbmanagerstate
type ManagerState int

const (
	ManagerStatePoweredOff   = ManagerState(C.CBManagerStatePoweredOff)
	ManagerStatePoweredOn    = ManagerState(C.CBManagerStatePoweredOn)
	ManagerStateResetting    = ManagerState(C.CBManagerStateResetting)
	ManagerStateUnauthorized = ManagerState(C.CBManagerStateUnauthorized)
	ManagerStateUnknown      = ManagerState(C.CBManagerStateUnknown)
	ManagerStateUnsupported  = ManagerState(C.CBManagerStateUnsupported)
)

// CentralManagerRestoreOpts: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/central_manager_state_restoration_options
type CentralManagerRestoreOpts struct {
	Peripherals            []Peripheral
	ScanServices           []UUID
	CentralManagerScanOpts *CentralManagerScanOpts // nil if none
}

// CentralManagerOpts: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/central_manager_initialization_options
type CentralManagerOpts struct {
	ShowPowerAlert    bool
	RestoreIdentifier string
}

// CentralManager: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager?language=objc
type CentralManager struct {
	ptr unsafe.Pointer
}

// DfltCentralManagerOpts is the set of options that gets used when nil is
// passed to `NewCentralManager()`.
var DfltCentralManagerOpts = CentralManagerOpts{
	ShowPowerAlert:    false,
	RestoreIdentifier: "",
}

var cmgrPtrMap = newPtrMap()

func findCentralManagerDlg(ptr unsafe.Pointer) CentralManagerDelegate {
	itf := cmgrPtrMap.find(ptr)
	if itf == nil {
		return nil
	}

	return itf.(CentralManagerDelegate)
}

// NewCentralManager creates a central manager.  Specify a nil `opts` value for
// defaults.  Don't forget to call `SetDelegate()` afterwards!
func NewCentralManager(opts *CentralManagerOpts) CentralManager {
	if opts == nil {
		opts = &DfltCentralManagerOpts
	}

	pwrAlert := C.bool(opts.ShowPowerAlert)

	restoreID := (*C.char)(nil)
	if opts.RestoreIdentifier != "" {
		restoreID = C.CString(opts.RestoreIdentifier)
		defer C.free(unsafe.Pointer(restoreID))
	}

	return CentralManager{
		ptr: unsafe.Pointer(C.cb_alloc_cmgr(pwrAlert, restoreID)),
	}
}

// SetDelegate configures a receiver for a central manager's asynchronous
// callbacks.
func (cm CentralManager) SetDelegate(d CentralManagerDelegate) {
	if d != nil {
		cmgrPtrMap.add(cm.ptr, d)
	}
	C.cb_cmgr_set_delegate(cm.ptr, C.bool(d != nil))
}

// State: https://developer.apple.com/documentation/corebluetooth/cbmanager/1648600-state
func (cm CentralManager) State() ManagerState {
	return ManagerState(C.cb_cmgr_state(cm.ptr))
}

// RetrieveConnectedPeripheralsWithServices: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1518924-retrieveconnectedperipheralswith
func (cm CentralManager) RetrieveConnectedPeripheralsWithServices(uuids []UUID) []Peripheral {
	strs := uuidsToStrArr(uuids)
	defer freeStrArr(&strs)

	var prphs []Peripheral

	prphPtrs := C.cb_cmgr_retrieve_prphs_with_svcs(cm.ptr, &strs)
	defer C.free(unsafe.Pointer(prphPtrs.objs))

	for i := 0; i < int(prphPtrs.count); i++ {
		ptr := getObjArrElem(&prphPtrs, i)
		prphs = append(prphs, Peripheral{ptr})
	}

	return prphs
}

// RetrievePeripheralsWithIdentifiers: https://developer.apple.com/documentation/corebluetooth/cbcentralmanager/1519127-retrieveperipheralswithidentifie
func (cm CentralManager) RetrievePeripheralsWithIdentifiers(uuids []UUID) []Peripheral {
	strs := uuidsToStrArr(uuids)
	defer freeStrArr(&strs)

	var prphs []Peripheral

	prphPtrs := C.cb_cmgr_retrieve_prphs(cm.ptr, &strs)
	defer C.free(unsafe.Pointer(prphPtrs.objs))

	for i := 0; i < int(prphPtrs.count); i++ {
		ptr := getObjArrElem(&prphPtrs, i)
		prphs = append(prphs, Peripheral{ptr})
	}

	return prphs
}
