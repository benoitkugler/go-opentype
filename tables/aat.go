package tables

// AAT layout

//go:generate ../../binarygen/cmd/generator aat.go

// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html - State tables
type aatStateTable struct {
	stateSize  uint16 // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	classTable uint16 // Byte offset from the beginning of the state table to the class subtable.
	stateArray uint16 // Byte offset from the beginning of the state table to the state array.
	entryTable uint16 // Byte offset from the beginning of the state table to the entry subtable.
	data       []byte `len:"_toEnd"`
}

type aatExtendedStateTable struct {
	stateSize  uint32 // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	classTable uint32 // Byte offset from the beginning of the state table to the class subtable.
	stateArray uint32 // Byte offset from the beginning of the state table to the state array.
	entryTable uint32 // Byte offset from the beginning of the state table to the entry subtable.
	data       []byte `len:"_toEnd"`
}

type aatStateEntry0 struct {
	newState uint16
	flags    uint16 // Table specific.
}

type aatStateEntry2 struct {
	newState uint16
	flags    uint16 // Table specific.
	data     [2]byte
}

type aatStateEntry4 struct {
	newState uint16
	flags    uint16 // Table specific.
	data     [4]byte
}

type aatLookupTable struct {
	version aatLookupTableVersion
}

type aatLookupTableVersion uint16

const (
	aatLookupTableVersion0  aatLookupTableVersion = 0
	aatLookupTableVersion2  aatLookupTableVersion = 2
	aatLookupTableVersion4  aatLookupTableVersion = 4
	aatLookupTableVersion6  aatLookupTableVersion = 6
	aatLookupTableVersion8  aatLookupTableVersion = 8
	aatLookupTableVersion10 aatLookupTableVersion = 10
)

type aatLookupTable0 struct {
	values []uint16 `len:""`
}

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 //		The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 //		The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 //		The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

type aatLookupTable2 struct {
	binSearchHeader
	records []classRangeRecord `len:"nUnits"`
}

type classRangeRecord struct {
	start, end glyphID
	value      uint16
}

type classFormat2 []classRangeRecord

type aatLookupTable4 struct{}

type aatLookupTable6 struct{}

type aatLookupTable8 struct{}

type aatLookupTable10 struct{}

type aatExtendedLookupTable struct {
	version uint16
}
