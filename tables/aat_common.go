package tables

import (
	"encoding/binary"
	"fmt"
)

// AAT layout

// State table header, without the actual data
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html
type aatSTHeader struct {
	stateSize  uint16 // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	classTable uint16 // Byte offset from the beginning of the state table to the class subtable.
	stateArray uint16 // Byte offset from the beginning of the state table to the state array.
	entryTable uint16 // Byte offset from the beginning of the state table to the entry subtable.
}

// Extended state table, including the data
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6Tables.html - State tables
// binarygen: argument=entryDataSize int
type AATStateTableExt struct {
	stateSize  uint32          // Size of a state, in bytes. The size is limited to 8 bits, although the field is 16 bits for alignment.
	Class      AATLookup       `offsetSize:"Offset32"` // Byte offset from the beginning of the state table to the class subtable.
	stateArray Offset32        // Byte offset from the beginning of the state table to the state array.
	entryTable Offset32        // Byte offset from the beginning of the state table to the entry subtable.
	States     [][]uint16      `isOpaque:""`
	Entries    []AATStateEntry `isOpaque:""`
}

func (state *AATStateTableExt) parseStates(src []byte, _, _ int) (int, error) {
	if state.stateArray > state.entryTable {
		return 0, fmt.Errorf("invalid AAT state table (%d > %d)", state.stateArray, state.entryTable)
	}
	if L := len(src); L < int(state.entryTable) {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", state.entryTable, L)
	}

	statesArray := src[state.stateArray:state.entryTable]
	states, err := parseUint16s(statesArray, len(statesArray)/2)
	if err != nil {
		return 0, err
	}

	nC := int(state.stateSize)
	// Ensure pre-defined classes fit.
	if nC < 4 {
		return 0, fmt.Errorf("invalid number of classes in state table: %d", nC)
	}
	state.States = make([][]uint16, len(states)/nC)
	for i := range state.States {
		state.States[i] = states[i*nC : (i+1)*nC]
	}
	return 0, nil
}

func (state *AATStateTableExt) parseEntries(src []byte, _, entryDataSize int) (int, error) {
	// find max index
	var maxi uint16
	for _, l := range state.States {
		for _, stateIndex := range l {
			if stateIndex > maxi {
				maxi = stateIndex
			}
		}
	}

	src = src[state.entryTable:] // checked in parseStates
	count := int(maxi) + 1
	entrySize := 4 + entryDataSize
	if L := len(src); L < count*entrySize {
		return 0, fmt.Errorf("EOF: expected length: %d, got %d", count*entrySize, L)
	}
	state.Entries = make([]AATStateEntry, count)
	for i := range state.Entries {
		state.Entries[i].NewState = binary.BigEndian.Uint16(src[i*entrySize:])
		state.Entries[i].Flags = binary.BigEndian.Uint16(src[i*entrySize+2:])
		copy(state.Entries[i].data[:], src[i*entrySize+4:(i+1)*entrySize])
	}

	// the own header data stop at the entryTable offset
	return 16, nil
}

type AATStateEntry struct {
	NewState uint16
	Flags    uint16  // Table specific.
	data     [4]byte // Table specific.
}

// AsMorxContextual reads the internal data for entries in morx contextual subtable.
// The returned indexes use 0xFFFF as empty value.
func (e AATStateEntry) AsMorxContextual() (markIndex, currentIndex uint16) {
	markIndex = binary.BigEndian.Uint16(e.data[:])
	currentIndex = binary.BigEndian.Uint16(e.data[2:])
	return
}

// AsMorxInsertion reads the internal data for entries in morx insertion subtable.
// The returned indexes use 0xFFFF as empty value.
func (e AATStateEntry) AsMorxInsertion() (currentIndex, markedIndex uint16) {
	currentIndex = binary.BigEndian.Uint16(e.data[:])
	markedIndex = binary.BigEndian.Uint16(e.data[2:])
	return
}

// AsMorxLigature reads the internal data for entries in morx ligature subtable.
func (e AATStateEntry) AsMorxLigature() (ligActionIndex uint16) {
	return binary.BigEndian.Uint16(e.data[:])
}

// AsKernxIndex reads the internal data for entries in 'kern/x' subtable format 1.
// An entry with no index returns 0xFFFF
func (e AATStateEntry) AsKernxIndex() uint16 {
	// for kern table, during parsing, we store the resolved index
	// at the same place as kerx tables
	return binary.BigEndian.Uint16(e.data[:])
}

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 // The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 // The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 // The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

// AATLookup is conceptually a map[GlyphID]uint16, but it may
// be implemented more efficiently.
type AATLookup interface {
	isAATLookup()
	// Class returns the class ID for the provided glyph, or (0, false)
	// for glyphs not covered by this class.
	Class(g GlyphID) (uint16, bool)
}

func (AATLoopkup0) isAATLookup()  {}
func (AATLoopkup2) isAATLookup()  {}
func (AATLoopkup4) isAATLookup()  {}
func (AATLoopkup6) isAATLookup()  {}
func (AATLoopkup8) isAATLookup()  {}
func (AATLoopkup10) isAATLookup() {}

type AATLoopkup0 struct {
	version uint16   `unionTag:"0"`
	Values  []uint16 `arrayCount:""`
}

type AATLoopkup2 struct {
	version uint16 `unionTag:"2"`
	binSearchHeader
	Records []LookupRecord2 `arrayCount:"ComputedField-nUnits"`
}

type LookupRecord2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint16
}

type AATLoopkup4 struct {
	version uint16 `unionTag:"4"`
	binSearchHeader
	// Do not include the termination segment
	Records []AATLookupRecord4 `arrayCount:"ComputedField-nUnits-1"`
}

type AATLookupRecord4 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	Values []uint16 `offsetSize:"Offset16" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nValues()"`
}

func (lk AATLookupRecord4) nValues() int { return int(lk.LastGlyph) - int(lk.FirstGlyph) + 1 }

type AATLoopkup6 struct {
	version uint16 `unionTag:"6"`
	binSearchHeader
	Records []loopkupRecord6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecord6 struct {
	Glyph GlyphID
	Value uint16
}

type AATLoopkup8 struct {
	version    uint16 `unionTag:"8"`
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

type AATLoopkup10 struct {
	version    uint16 `unionTag:"10"`
	unitSize   uint16
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

// extended versions

type AATLookupExt interface {
	isAATLookupExt()
}

func (AATLoopkupExt0) isAATLookupExt()  {}
func (AATLoopkupExt2) isAATLookupExt()  {}
func (AATLoopkupExt4) isAATLookupExt()  {}
func (AATLoopkupExt6) isAATLookupExt()  {}
func (AATLoopkupExt8) isAATLookupExt()  {}
func (AATLoopkupExt10) isAATLookupExt() {}

type AATLoopkupExt0 struct {
	version uint16   `unionTag:"0"`
	Values  []uint32 `arrayCount:""`
}

type AATLoopkupExt2 struct {
	version uint16 `unionTag:"2"`
	binSearchHeader
	Records []lookupRecordExt2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecordExt2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint32
}

type AATLoopkupExt4 struct {
	version uint16 `unionTag:"4"`
	binSearchHeader
	// the values pointed by the record are uint32
	Records []loopkupRecordExt4 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecordExt4 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	Values []uint32 `offsetSize:"Offset16" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nValues()"`
}

type AATLoopkupExt6 struct {
	version uint16 `unionTag:"6"`
	binSearchHeader
	Records []loopkupRecordExt6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecordExt6 struct {
	Glyph GlyphID
	Value uint32
}

type AATLoopkupExt8 struct {
	version    uint16 `unionTag:"8"`
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

type AATLoopkupExt10 struct {
	version    uint16 `unionTag:"10"`
	unitSize   uint16
	FirstGlyph GlyphID
	Values     []uint32 `arrayCount:"FirstUint16"`
}
