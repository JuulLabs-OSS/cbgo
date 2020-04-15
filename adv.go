package cbgo

type AdvServiceData struct {
	UUID UUID
	Data []byte
}

// AdvFields represents the contents of an advertisement received during
// scanning.
type AdvFields struct {
	LocalName        string
	ManufacturerData []byte
	TxPowerLevel     *int
	Connectable      *bool
	ServiceUUIDs     []UUID
	ServiceData      []AdvServiceData
}
