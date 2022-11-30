package tables

// AAT layout

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 // The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 // The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 // The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

type AATLookup interface {
	isAATLookup()
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
	Records []lookupRecord2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecord2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint16
}

type AATLoopkup4 struct {
	version uint16 `unionTag:"4"`
	binSearchHeader
	Records []loopkupRecord4 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecord4 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	Values []uint16 `offsetSize:"Offset16" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nValues()"`
}

func (lk loopkupRecord4) nValues() int { return int(lk.LastGlyph) - int(lk.FirstGlyph) + 1 }

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
