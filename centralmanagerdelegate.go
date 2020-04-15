package cbgo

// CentralManagerDelegate: https://developer.apple.com/documentation/corebluetooth/cbcentralmanagerdelegate
type CentralManagerDelegate interface {
	DidConnectPeripheral(cmgr CentralManager, prph Peripheral)
	DidFailToConnectPeripheral(cmgr CentralManager, prph Peripheral, err error)
	DidDisconnectPeripheral(cmgr CentralManager, prph Peripheral, err error)
	DidUpdateState(cmgr CentralManager)
	WillRestoreState(cmgr CentralManager, opts CentralManagerRestoreOpts)
	DidDiscoverPeripheral(cmgr CentralManager, prph Peripheral, advFields AdvFields, rssi int)
}

// CentralManagerDelegateBase implements the CentralManagerDelegate interface
// with stub functions.  Embed this in your delegate type if you only want to
// define a subset of the CentralManagerDelegate interface.
type CentralManagerDelegateBase struct {
}

func (b *CentralManagerDelegateBase) DidConnectPeripheral(cmgr CentralManager, prph Peripheral) {
}
func (b *CentralManagerDelegateBase) DidFailToConnectPeripheral(cmgr CentralManager, prph Peripheral, err error) {
}
func (b *CentralManagerDelegateBase) DidDisconnectPeripheral(cmgr CentralManager, prph Peripheral, err error) {
}
func (b *CentralManagerDelegateBase) DidUpdateState(cmgr CentralManager) {
}
func (b *CentralManagerDelegateBase) WillRestoreState(cmgr CentralManager, opts CentralManagerRestoreOpts) {
}
func (b *CentralManagerDelegateBase) DidDiscoverPeripheral(cmgr CentralManager, prph Peripheral, advFields AdvFields, rssi int) {
}
