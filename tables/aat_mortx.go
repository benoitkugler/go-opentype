package tables

type mortxSubtable struct {
	tableLength uint32
	coverage    byte
	padding     uint16
	kind        byte
	flags       uint32
	subtable    []byte `len-size:"-"`
}
