package tables

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type KernSubtable interface {
	// IsHorizontal returns true if the subtable has horizontal kerning values.
	IsHorizontal() bool
	// IsBackwards returns true if state-table based should process the glyphs backwards.
	IsBackwards() bool
	// IsCrossStream returns true if the subtable has cross-stream kerning values.
	IsCrossStream() bool
	// IsVariation returns true if the subtable has variation kerning values.
	IsVariation() bool

	// Data returns the actual kerning data
	Data() KernData
}

type OTKernSubtableHeader struct {
	version  uint16        // Kern subtable version number
	length   uint16        // Length of the subtable, in bytes (including this header).
	format   kernSTVersion // What type of information is contained in this table.
	coverage byte          // What type of information is contained in this table.
	data     KernData      `unionField:"format"`
}

// check and return the length
func (st *OTKernSubtableHeader) parseEnd(src []byte) (int, error) {
	if L, E := len(src), int(st.length); L < E {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}
	return int(st.length), nil
}

type AATKernSubtableHeader struct {
	length     uint32 // The length of this subtable in bytes, including this header.
	coverage   byte   // Circumstances under which this table is used.
	version    kernSTVersion
	tupleCount uint16   // The tuple count. This value is only used with variation fonts and should be 0 for all other fonts. The subtable's tupleCount will be ignored if the 'kerx' table version is less than 4.
	data       KernData `unionField:"version"`
}

// check and return the length
func (st *AATKernSubtableHeader) parseEnd(src []byte) (int, error) {
	if L, E := len(src), int(st.length); L < E {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", E, L)
	}
	return int(st.length), nil
}

type KernData interface {
	isKernData()
}

func (KernData0) isKernData() {}
func (KernData1) isKernData() {}
func (KernData2) isKernData() {}
func (KernData3) isKernData() {}

type kernSTVersion byte

const (
	kernSTVersion0 kernSTVersion = iota
	kernSTVersion1
	kernSTVersion2
	kernSTVersion3
)

type KernData0 struct {
	nPairs        uint16         //	The number of kerning pairs in this subtable.
	searchRange   uint16         //	The largest power of two less than or equal to the value of nPairs, multiplied by the size in bytes of an entry in the subtable.
	entrySelector uint16         //	This is calculated as log2 of the largest power of two less than or equal to the value of nPairs. This value indicates how many iterations of the search loop have to be made. For example, in a list of eight items, there would be three iterations of the loop.
	rangeShift    uint16         //	The value of nPairs minus the largest power of two less than or equal to nPairs. This is multiplied b
	Pairs         []Kernx0Record `arrayCount:"ComputedField-nPairs"`
}

type KernData1 struct {
	AATStateTable
	valueTable uint16  // Offset in bytes from the beginning of the subtable to the beginning of the kerning table.
	Values     []int16 `isOpaque:""`
}

func (kd *KernData1) parseValues(src []byte) (int, error) {
	valuesOffset := int(kd.valueTable)
	// start by resolving offset -> index
	for i := range kd.Entries {
		entry := &kd.Entries[i]
		offset := int(entry.Flags & Kern1Offset)
		if offset == 0 || offset < valuesOffset {
			binary.BigEndian.PutUint16(entry.data[:], 0xFFFF)
		} else {
			index := uint16((offset - valuesOffset) / 2)
			binary.BigEndian.PutUint16(entry.data[:], index)
		}
	}
	var err error
	kd.Values, err = parseKernx1Values(src, kd.Entries, valuesOffset, 0)
	return len(src), err
}

type KernData2 struct {
	rowWidth    uint16          // The width, in bytes, of a row in the subtable.
	left        AATLoopkup8Data `offsetSize:"Offset16" offsetRelativeTo:"Parent"`
	right       AATLoopkup8Data `offsetSize:"Offset16" offsetRelativeTo:"Parent"`
	array       Offset16        // Offset from beginning of this subtable to the start of the kerning array.
	kerningData []byte          `isOpaque:"" offsetRelativeTo:"Parent"` // indexed by Left + Right
}

func (kd *KernData2) parseKerningData(_ []byte, parentSrc []byte) (int, error) {
	kd.kerningData = parentSrc
	return 0, nil
}

type KernData3 struct {
	glyphCount      uint16  // The number of glyphs in this font.
	kernValueCount  uint8   // The number of kerning values.
	leftClassCount  uint8   // The number of left-hand classes.
	rightClassCount uint8   // The number of right-hand classes.
	flags           uint8   // Set to zero (reserved for future use).
	kernings        []int16 `arrayCount:"ComputedField-kernValueCount"`
	leftClass       []uint8 `arrayCount:"ComputedField-glyphCount"`
	rightClass      []uint8 `arrayCount:"ComputedField-glyphCount"`
	kernIndex       []uint8 `arrayCount:"ComputedField-nKernIndex()"`
}

func (kd *KernData3) nKernIndex() int { return int(kd.leftClassCount) * int(kd.rightClassCount) }

// sanitize index and class values
func (kd *KernData3) parseEnd(_ byte) (int, error) {
	for _, index := range kd.kernIndex {
		if index >= kd.kernValueCount {
			return 0, errors.New("invalid kern subtable format 3 index value")
		}
	}

	for i := range kd.leftClass {
		if kd.leftClass[i] >= kd.leftClassCount {
			return 0, errors.New("invalid kern subtable format 3 left class value")
		}
		if kd.rightClass[i] >= kd.rightClassCount {
			return 0, errors.New("invalid kern subtable format 3 right class value")
		}
	}

	return 0, nil
}
