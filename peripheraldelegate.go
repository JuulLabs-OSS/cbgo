package cbgo

// PeripheralDelegate: https://developer.apple.com/documentation/corebluetooth/cbperipheraldelegate
type PeripheralDelegate interface {
	DidDiscoverServices(prph Peripheral, err error)
	DidDiscoverIncludedServices(prph Peripheral, svc Service, err error)
	DidDiscoverCharacteristics(prph Peripheral, svc Service, err error)
	DidDiscoverDescriptors(prph Peripheral, chr Characteristic, err error)
	DidUpdateValueForCharacteristic(prph Peripheral, chr Characteristic, err error)
	DidUpdateValueForDescriptor(prph Peripheral, dsc Descriptor, err error)
	DidWriteValueForCharacteristic(prph Peripheral, chr Characteristic, err error)
	DidWriteValueForDescriptor(prph Peripheral, dsc Descriptor, err error)
	IsReadyToSendWriteWithoutResponse(prph Peripheral)
	DidUpdateNotificationState(prph Peripheral, chr Characteristic, err error)
	DidReadRSSI(prph Peripheral, rssi int, err error)
	DidUpdateName(prph Peripheral)
	DidModifyServices(prph Peripheral, invSvcs []Service)
}

// PeripheralDelegateBase implements the PeripheralDelegate interface with stub
// functions.  Embed this in your delegate type if you only want to define a
// subset of the PeripheralDelegate interface.
type PeripheralDelegateBase struct {
}

func (b *PeripheralDelegateBase) DidDiscoverServices(prph Peripheral, err error) {
}
func (b *PeripheralDelegateBase) DidDiscoverIncludedServices(prph Peripheral, svc Service, err error) {
}
func (b *PeripheralDelegateBase) DidDiscoverCharacteristics(prph Peripheral, svc Service, err error) {
}
func (b *PeripheralDelegateBase) DidDiscoverDescriptors(prph Peripheral, chr Characteristic, err error) {
}
func (b *PeripheralDelegateBase) DidUpdateValueForCharacteristic(prph Peripheral, chr Characteristic, err error) {
}
func (b *PeripheralDelegateBase) DidUpdateValueForDescriptor(prph Peripheral, dsc Descriptor, err error) {
}
func (b *PeripheralDelegateBase) DidWriteValueForCharacteristic(prph Peripheral, chr Characteristic, err error) {
}
func (b *PeripheralDelegateBase) DidWriteValueForDescriptor(prph Peripheral, dsc Descriptor, err error) {
}
func (b *PeripheralDelegateBase) IsReadyToSendWriteWithoutResponse(prph Peripheral) {
}
func (b *PeripheralDelegateBase) DidUpdateNotificationState(prph Peripheral, chr Characteristic, err error) {
}
func (b *PeripheralDelegateBase) DidReadRSSI(prph Peripheral, rssi int, err error) {
}
func (b *PeripheralDelegateBase) DidUpdateName(prph Peripheral) {
}
func (b *PeripheralDelegateBase) DidModifyServices(prph Peripheral, invSvcs []Service) {
}
