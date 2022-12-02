package tables

import (
	"encoding/binary"
	"fmt"
)

// Kerx is the extended kerning table
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6kerx.html
type Kerx struct {
	version uint16         // The version number of the extended kerning table (currently 2, 3, or 4).
	padding uint16         // Unused; set to zero.
	nTables uint32         // The number of subtables included in the extended kerning table.
	Tables  []KerxSubtable `arrayCount:"ComputedField-nTables"`
}

// extended versions

// binarygen: argument=valuesCount int
type KerxSubtable struct {
	length     uint32 // The length of this subtable in bytes, including this header.
	Coverage   uint16 // Circumstances under which this table is used.
	padding    byte   // unused
	version    kerxSTVersion
	TupleCount uint32   // The tuple count. This value is only used with variation fonts and should be 0 for all other fonts. The subtable's tupleCount will be ignored if the 'kerx' table version is less than 4.
	Data       KerxData `unionField:"version" arguments:"tupleCount=.TupleCount, valuesCount=valuesCount"`
}

// check and return the subtable length
func (ks *KerxSubtable) parseEnd(src []byte) (int, error) {
	if L := len(src); L < int(ks.length) {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", ks.length, L)
	}
	return int(ks.length), nil
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

type KerxData interface {
	isKerxData()
}

func (KerxData0) isKerxData() {}
func (KerxData1) isKerxData() {}
func (KerxData2) isKerxData() {}
func (KerxData4) isKerxData() {}
func (KerxData6) isKerxData() {}

type KerxData0 struct {
	nPairs        uint32
	searchRange   uint32
	entrySelector uint32
	rangeShift    uint32
	Pairs         []Kernx0Record `arrayCount:"ComputedField-nPairs"`
	rawData       []byte         `arrayCount:"ToEnd" subsliceStart:"AtStart"` // not empty only for variable fonts
}

type Kernx0Record struct {
	Left, Right GlyphID
	value       int16
}

// binarygen: argument=tupleCount int
// binarygen: argument=valuesCount int
type KerxData1 struct {
	AATStateTableExt `arguments:"valuesCount=valuesCount, entryDataSize=2"`
	valueTable       Offset32
	Values           []int16 `isOpaque:""`
}

// From Apple 'kern' spec:
// Each pops one glyph from the kerning stack and applies the kerning value to it.
// The end of the list is marked by an odd value...
func (kx *KerxData1) parseValues(src []byte, tupleCount, _ int) (int, error) {
	// find the maximum index need in the values array
	var maxi uint16
	for _, entry := range kx.Entries {
		if index := entry.AsKernxIndex(); index != 0xFFFF && index > maxi {
			maxi = index
		}
	}

	if tupleCount == 0 {
		tupleCount = 1
	}
	nbUint16Min := tupleCount * int(maxi+1)
	if L, E := len(src), int(kx.valueTable)+2*nbUint16Min; L < E {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}

	src = src[kx.valueTable:]
	kx.Values = make([]int16, 0, nbUint16Min)
	for len(src) >= 2 { // gracefully handle missing odd value
		v := int16(binary.BigEndian.Uint16(src))
		kx.Values = append(kx.Values, v)
		src = src[2:]
		if len(kx.Values) >= nbUint16Min && v&1 != 0 {
			break
		}
	}
	return len(src), nil
}

type KerxData2 struct {
	rowWidth    uint32    // The number of bytes in each row of the kerning value array
	Left        AATLookup `offsetSize:"Offset32" offsetRelativeTo:"Parent"` // Offset from beginning of this subtable to the left-hand offset table.
	Right       AATLookup `offsetSize:"Offset32" offsetRelativeTo:"Parent"` // Offset from beginning of this subtable to right-hand offset table.
	array       Offset32  // Offset from beginning of this subtable to the start of the kerning array.
	kerningData []byte    `subsliceStart:"AtStart" arrayCount:"ToEnd"` // indexed by Left + Right
}

type KerxData4 struct{}

type KerxData6 struct{}

// -------------------------------- kern table --------------------------------

type microsoftKern struct {
	version uint16 // Table version number (0)
	nTables uint16 // Number of subtables in the kerning table.
}

type kernSubtableW struct {
	length     uint32 // The length of this subtable in bytes, including this header.
	coverage   byte   // Circumstances under which this table is used.
	version    kernSTVersion
	tupleCount uint16       // The tuple count. This value is only used with variation fonts and should be 0 for all other fonts. The subtable's tupleCount will be ignored if the 'kerx' table version is less than 4.
	data       kernSubtable `unionField:"version"`
}

type kernSubtable interface {
	isKernSubtable()
}

func (kernSubtable0) isKernSubtable() {}
func (kernSubtable1) isKernSubtable() {}
func (kernSubtable2) isKernSubtable() {}

type kernSTVersion byte

const (
	kernSTVersion0 kernSTVersion = iota
	kernSTVersion1
	kernSTVersion2
	// kernSTVersion3 // TODO:
)

type kernSubtable0 struct {
	nPairs        uint16         //	The number of kerning pairs in this subtable.
	searchRange   uint16         //	The largest power of two less than or equal to the value of nPairs, multiplied by the size in bytes of an entry in the subtable.
	entrySelector uint16         //	This is calculated as log2 of the largest power of two less than or equal to the value of nPairs. This value indicates how many iterations of the search loop have to be made. For example, in a list of eight items, there would be three iterations of the loop.
	rangeShift    uint16         //	The value of nPairs minus the largest power of two less than or equal to nPairs. This is multiplied b
	pairs         []Kernx0Record `arrayCount:"ComputedField-nPairs"`
	rawData       []byte         `arrayCount:"ToEnd"` // used for variable fonts
}

type kernSubtable1 struct {
	aatSTHeader
	valueTable uint16 // Offset in bytes from the beginning of the subtable to the beginning of the kerning table.
	rawData    []byte `arrayCount:"ToEnd" subsliceStart:"AtStart"`
}

type kernSubtable2 struct {
	rowWidth           uint16      // The width, in bytes, of a row in the subtable.
	left               AATLoopkup8 `offsetSize:"Offset16"`
	right              AATLoopkup8 `offsetSize:"Offset16"`
	kerningArrayOffset uint16
	rawData            []byte `arrayCount:"ToEnd" subsliceStart:"AtStart"` // used to resolve kerning pairs
}
