package tables

// AAT layout

// go:generate ../../binarygen/cmd/generator aat_common.go

type binSearchHeader struct {
	unitSize      uint16
	nUnits        uint16
	searchRange   uint16 // The value of unitSize times the largest power of 2 that is less than or equal to the value of nUnits.
	entrySelector uint16 // The log base 2 of the largest power of 2 less than or equal to the value of nUnits.
	rangeShift    uint16 // The value of unitSize times the difference of the value of nUnits minus the largest power of 2 less than or equal to the value of nUnits.
}

type aatLookup struct {
	version aatLookupVersion
	table   aatLookupTable `version-field:"version"`
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
	values []uint16 `len:""`
}

type aatLookupTable2 struct {
	binSearchHeader
	records []lookupRecord2 `len:"nUnits"`
}

type lookupRecord2 struct {
	lastGlyph  glyphID
	firstGlyph glyphID
	value      uint16
}

type aatLookupTable4 struct {
	binSearchHeader
	records []loopkupRecord4 `len:"nUnits"`
	rawData []byte           `len:"__toEnd"`
}

type loopkupRecord4 struct {
	lastGlyph  glyphID
	firstGlyph glyphID
	// offset to an array of []uint16 (or []uint32 for extended) with length last - first + 1
	offsetToValues uint16
}

type aatLookupTable6 struct {
	binSearchHeader
	records []loopkupRecord6 `len:"nUnits"`
}

type loopkupRecord6 struct {
	glyph glyphID
	value uint16
}

type aatLookupTable8 struct {
	firstGlyph glyphID
	values     []uint16 `len:"_first16"`
}

type aatLookupTable10 struct {
	unitSize   uint16
	firstGlyph glyphID
	values     []uint16 `len:"_first16"`
}

type aatExtendedLookupTable struct {
	version uint16
}

// extended versions

type aatLookupExt struct {
	version aatLookupVersion
	table   aatLookupTableExt `version-field:"version"`
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
	values []uint32 `len:""`
}

type aatLookupTableExt2 struct {
	binSearchHeader
	records []lookupRecordExt2 `len:"nUnits"`
}

type lookupRecordExt2 struct {
	lastGlyph  glyphID
	firstGlyph glyphID
	value      uint32
}

type aatLookupTableExt4 struct {
	binSearchHeader
	// the values pointed by the record are uint32
	records []loopkupRecord4 `len:"nUnits"`
	rawData []byte           `len:"__toEnd"`
}

type aatLookupTableExt6 struct {
	binSearchHeader
	records []loopkupRecordExt6 `len:"nUnits"`
}

type loopkupRecordExt6 struct {
	glyph glyphID
	value uint32
}

type aatLookupTableExt8 struct {
	firstGlyph glyphID
	values     []uint16 `len:"_first16"`
}

type aatLookupTableExt10 struct {
	unitSize   uint16
	firstGlyph glyphID
	values     []uint32 `len:"_first16"`
}
