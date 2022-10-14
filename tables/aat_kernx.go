package tables

//go:generate ../../binarygen/cmd/generator aat_kernx.go

// Kerx is the extendered kerning table
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6kerx.html
type Kerx struct {
	version uint16 // The version number of the extended kerning table (currently 2, 3, or 4).
	padding uint16 // Unused; set to zero.
	nTables uint32 // The number of subtables included in the extended kerning table.
	rawData []byte `len:"__toEnd"`
}

type kernSubtableW struct {
	length     uint32 // The length of this subtable in bytes, including this header.
	coverage   byte   // Circumstances under which this table is used.
	version    kernSTVersion
	tupleCount uint16       // The tuple count. This value is only used with variation fonts and should be 0 for all other fonts. The subtable's tupleCount will be ignored if the 'kerx' table version is less than 4.
	data       kernSubtable `version-field:"version"`
}

type kernSubtable interface {
	isKernSubtable()
}

type kernSTVersion byte

const (
	kernSTVersion0 kernSTVersion = iota
	kernSTVersion1
	kernSTVersion2
	kernSTVersion3
)

type kernSubtable0 struct {
	nPairs        uint16         //	The number of kerning pairs in this subtable.
	searchRange   uint16         //	The largest power of two less than or equal to the value of nPairs, multiplied by the size in bytes of an entry in the subtable.
	entrySelector uint16         //	This is calculated as log2 of the largest power of two less than or equal to the value of nPairs. This value indicates how many iterations of the search loop have to be made. For example, in a list of eight items, there would be three iterations of the loop.
	rangeShift    uint16         //	The value of nPairs minus the largest power of two less than or equal to nPairs. This is multiplied b
	pairs         []kernx0Record `len:"nPairs"`
	rawData       []byte         `len:"__toEnd"` // used for variable fonts
}

type kernSubtable1 struct {
	aatSTHeader
	valueTable uint16 // Offset in bytes from the beginning of the subtable to the beginning of the kerning table.
	rawData    []byte `len:"__startToEnd"`
}

type kernSubtable2 struct {
	rowWidth           uint16          // The width, in bytes, of a row in the subtable.
	left               aatLookupTable8 `offset-size:"16"`
	right              aatLookupTable8 `offset-size:"16"`
	kerningArrayOffset uint16
	rawData            []byte `len:"__startToEnd"` // used to resolve kerning pairs
}

// extended versions

type kerxSubtableW struct {
	length     uint32 // The length of this subtable in bytes, including this header.
	coverage   uint16 // Circumstances under which this table is used.
	padding    byte   // unused
	version    kerxSTVersion
	tupleCount uint32       // The tuple count. This value is only used with variation fonts and should be 0 for all other fonts. The subtable's tupleCount will be ignored if the 'kerx' table version is less than 4.
	data       kerxSubtable `version-field:"version"`
}

type kerxSTVersion byte

const (
	kerxSTVersion0 kerxSTVersion = iota
	kerxSTVersion1
	kerxSTVersion2
	_
	kerxSTVersion4
	_
	kerxSTVersion6
)

type kerxSubtable interface {
	isKerxSubtable()
}

type kerxSubtable0 struct {
	nPairs        uint32
	searchRange   uint32
	entrySelector uint32
	rangeShift    uint32
	pairs         []kernx0Record `len:"nPairs"`
	rawData       []byte         `len:"__toEnd"` // used for variable fonts
}

type kernx0Record struct {
	left, right glyphID
	value       int16
}
