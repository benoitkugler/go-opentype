package tables

// State table header, without the actual data
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html
type aatSTHeader struct {
	stateSize  uint16 // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	classTable uint16 // Byte offset from the beginning of the state table to the class subtable.
	stateArray uint16 // Byte offset from the beginning of the state table to the state array.
	entryTable uint16 // Byte offset from the beginning of the state table to the entry subtable.
}

// Extended state table header, without the actual data
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html - State tables
type aatSTXHeader struct {
	stateSize  uint32 // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	classTable uint32 // Byte offset from the beginning of the state table to the class subtable.
	stateArray uint32 // Byte offset from the beginning of the state table to the state array.
	entryTable uint32 // Byte offset from the beginning of the state table to the entry subtable.
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
