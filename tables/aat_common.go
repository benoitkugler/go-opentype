package tables

// AAT layout

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 // The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 // The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 // The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

type aatLookup interface {
	isAATLookupTable()
}

func (aatLookupTable0) isAATLookupTable()  {}
func (aatLookupTable2) isAATLookupTable()  {}
func (aatLookupTable4) isAATLookupTable()  {}
func (aatLookupTable6) isAATLookupTable()  {}
func (aatLookupTable8) isAATLookupTable()  {}
func (aatLookupTable10) isAATLookupTable() {}

type aatLookupTable0 struct {
	version uint16   `unionTag:"0"`
	Values  []uint16 `arrayCount:""`
}

type aatLookupTable2 struct {
	version uint16 `unionTag:"2"`
	binSearchHeader
	Records []lookupRecord2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecord2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint16
}

type aatLookupTable4 struct {
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

type aatLookupTable6 struct {
	version uint16 `unionTag:"6"`
	binSearchHeader
	Records []loopkupRecord6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecord6 struct {
	Glyph GlyphID
	Value uint16
}

type aatLookupTable8 struct {
	version    uint16 `unionTag:"8"`
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

type aatLookupTable10 struct {
	version    uint16 `unionTag:"10"`
	unitSize   uint16
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

// extended versions

type aatLookupExt interface {
	isAATLookupTableExt()
}

func (aatLookupTableExt0) isAATLookupTableExt()  {}
func (aatLookupTableExt2) isAATLookupTableExt()  {}
func (aatLookupTableExt4) isAATLookupTableExt()  {}
func (aatLookupTableExt6) isAATLookupTableExt()  {}
func (aatLookupTableExt8) isAATLookupTableExt()  {}
func (aatLookupTableExt10) isAATLookupTableExt() {}

type aatLookupTableExt0 struct {
	version uint16   `unionTag:"0"`
	Values  []uint32 `arrayCount:""`
}

type aatLookupTableExt2 struct {
	version uint16 `unionTag:"2"`
	binSearchHeader
	Records []lookupRecordExt2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecordExt2 struct {
	LastGlyph  GlyphID
	FirstGlyph GlyphID
	Value      uint32
}

type aatLookupTableExt4 struct {
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

type aatLookupTableExt6 struct {
	version uint16 `unionTag:"6"`
	binSearchHeader
	Records []loopkupRecordExt6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecordExt6 struct {
	Glyph GlyphID
	Value uint32
}

type aatLookupTableExt8 struct {
	version    uint16 `unionTag:"8"`
	FirstGlyph GlyphID
	Values     []uint16 `arrayCount:"FirstUint16"`
}

type aatLookupTableExt10 struct {
	version    uint16 `unionTag:"10"`
	unitSize   uint16
	FirstGlyph GlyphID
	Values     []uint32 `arrayCount:"FirstUint16"`
}
