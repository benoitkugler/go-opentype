package tables

//go:generate ../../binarygen/cmd/generator kernx.go

type kernSubtable0 struct {
	nPairs        uint16         //	The number of kerning pairs in this subtable.
	searchRange   uint16         //	The largest power of two less than or equal to the value of nPairs, multiplied by the size in bytes of an entry in the subtable.
	entrySelector uint16         //	This is calculated as log2 of the largest power of two less than or equal to the value of nPairs. This value indicates how many iterations of the search loop have to be made. For example, in a list of eight items, there would be three iterations of the loop.
	rangeShift    uint16         //	The value of nPairs minus the largest power of two less than or equal to nPairs. This is multiplied b
	pairs         []kernx0Record `len:"nPairs"`
}

type kerxSubtable0 struct {
	nPairs        uint32
	searchRange   uint32
	entrySelector uint32
	rangeShift    uint32
	pairs         []kernx0Record `len:"nPairs"`
}

type kernx0Record struct {
	left, right glyphID
	value       int16
}
