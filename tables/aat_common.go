package tables

// AAT layout

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 // The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 // The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 // The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

type aatLookup struct {
	version aatLookupVersion
	table   aatLookupTable `unionField:"version"`
}

type aatLookupTable interface {
	isAATLookupTable()
}

func (aatLookupTable0) isAATLookupTable()  {}
func (aatLookupTable2) isAATLookupTable()  {}
func (aatLookupTable4) isAATLookupTable()  {}
func (aatLookupTable6) isAATLookupTable()  {}
func (aatLookupTable8) isAATLookupTable()  {}
func (aatLookupTable10) isAATLookupTable() {}

type aatLookupVersion uint16

const (
	aatLookupTableVersion0  aatLookupVersion = 0
	aatLookupTableVersion2  aatLookupVersion = 2
	aatLookupTableVersion4  aatLookupVersion = 4
	aatLookupTableVersion6  aatLookupVersion = 6
	aatLookupTableVersion8  aatLookupVersion = 8
	aatLookupTableVersion10 aatLookupVersion = 10
)

type aatLookupTable0 struct {
	values []uint16 `arrayCount:""`
}

type aatLookupTable2 struct {
	binSearchHeader
	records []lookupRecord2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecord2 struct {
	lastGlyph  GlyphID
	firstGlyph GlyphID
	value      uint16
}

type aatLookupTable4 struct {
	binSearchHeader
	records []loopkupRecord4 `arrayCount:"ComputedField-nUnits"`
	rawData []byte           `arrayCount:"ToEnd"`
}

type loopkupRecord4 struct {
	lastGlyph  GlyphID
	firstGlyph GlyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	offsetToValues uint16
}

type aatLookupTable6 struct {
	binSearchHeader
	records []loopkupRecord6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecord6 struct {
	glyph GlyphID
	value uint16
}

type aatLookupTable8 struct {
	firstGlyph GlyphID
	values     []uint16 `arrayCount:"FirstUint16"`
}

type aatLookupTable10 struct {
	unitSize   uint16
	firstGlyph GlyphID
	values     []uint16 `arrayCount:"FirstUint16"`
}

type aatExtendedLookupTable struct {
	version uint16
}

// extended versions

type aatLookupExt struct {
	version aatLookupVersion
	table   aatLookupTableExt `unionField:"version"`
}

type aatLookupTableExt interface {
	isAATLookupTableExt()
}

func (aatLookupTableExt0) isAATLookupTableExt()  {}
func (aatLookupTableExt2) isAATLookupTableExt()  {}
func (aatLookupTableExt4) isAATLookupTableExt()  {}
func (aatLookupTableExt6) isAATLookupTableExt()  {}
func (aatLookupTableExt8) isAATLookupTableExt()  {}
func (aatLookupTableExt10) isAATLookupTableExt() {}

type aatLookupTableExt0 struct {
	values []uint32 `arrayCount:""`
}

type aatLookupTableExt2 struct {
	binSearchHeader
	records []lookupRecordExt2 `arrayCount:"ComputedField-nUnits"`
}

type lookupRecordExt2 struct {
	lastGlyph  GlyphID
	firstGlyph GlyphID
	value      uint32
}

type aatLookupTableExt4 struct {
	binSearchHeader
	// the values pointed by the record are uint32
	records []loopkupRecord4 `arrayCount:"ComputedField-nUnits"`
	rawData []byte           `arrayCount:"ToEnd"`
}

type aatLookupTableExt6 struct {
	binSearchHeader
	records []loopkupRecordExt6 `arrayCount:"ComputedField-nUnits"`
}

type loopkupRecordExt6 struct {
	glyph GlyphID
	value uint32
}

type aatLookupTableExt8 struct {
	firstGlyph GlyphID
	values     []uint16 `arrayCount:"FirstUint16"`
}

type aatLookupTableExt10 struct {
	unitSize   uint16
	firstGlyph GlyphID
	values     []uint32 `arrayCount:"FirstUint16"`
}
