package tables

import (
	"encoding/binary"
	"fmt"
)

type SinglePos struct {
	posFormat singlePosVersion
	Data      SinglePosData `unionField:"posFormat" subsliceStart:"AtStart"`
}

type singlePosVersion uint16

const (
	singlePosVersion1 singlePosVersion = 1
	singlePosVersion2 singlePosVersion = 2
)

type SinglePosData interface {
	isSinglePosData()
}

func (SinglePosData1) isSinglePosData() {}
func (SinglePosData2) isSinglePosData() {}

// binarygen: startOffset=2
type SinglePosData1 struct {
	Coverage    Coverage    `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of SinglePos subtable.
	valueFormat ValueFormat //	Defines the types of data in the ValueRecord.
	ValueRecord ValueRecord `isOpaque:"" subsliceStart:"AtStart"` //	Defines positioning value(s) — applied to all glyphs in the Coverage table.
}

func (sp *SinglePosData1) customParseValueRecord(src []byte) (read int, err error) {
	sp.ValueRecord, read, err = parseValueRecord(sp.valueFormat, src, 6)
	return read, err
}

// binarygen: startOffset=2
type SinglePosData2 struct {
	Coverage     Coverage      `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of SinglePos subtable.
	valueFormat  ValueFormat   // Defines the types of data in the ValueRecords.
	valueCount   uint16        // Number of ValueRecords — must equal glyphCount in the Coverage table.
	ValueRecords []ValueRecord `isOpaque:"" subsliceStart:"AtStart"` //[valueCount]	Array of ValueRecords — positioning values applied to glyphs.
}

func (sp *SinglePosData2) customParseValueRecords(src []byte) (read int, err error) {
	offset := 8
	sp.ValueRecords = make([]ValueRecord, sp.valueCount)
	for i := range sp.ValueRecords {
		sp.ValueRecords[i], offset, err = parseValueRecord(sp.valueFormat, src, offset)
		if err != nil {
			return 0, err
		}
	}

	return offset, err
}

type PairPos struct {
	posFormat pairPosVersion
	Data      PairPosData `unionField:"posFormat" subsliceStart:"AtStart"`
}

type pairPosVersion uint16

const (
	pairPosVersion1 pairPosVersion = 1
	pairPosVersion2 pairPosVersion = 2
)

type PairPosData interface {
	isPairPosData()
}

func (PairPosData1) isPairPosData() {}
func (PairPosData2) isPairPosData() {}

// binarygen: startOffset=2
type PairPosData1 struct {
	Coverage Coverage `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of PairPos subtable.

	valueFormat1  ValueFormat // Defines the types of data in valueRecord1 — for the first glyph in the pair (may be zero).
	valueFormat2  ValueFormat // Defines the types of data in valueRecord2 — for the second glyph in the pair (may be zero).
	PairSetOffset []PairSet   `arrayCount:"FirstUint16" offsetsArray:"Offset16" arguments:"valueFormat1, valueFormat2"` //[pairSetCount] Array of offsets to PairSet tables. Offsets are from beginning of PairPos subtable, ordered by Coverage Index.
}

// binarygen: argument=valueFormat1  ValueFormat
// binarygen: argument=valueFormat2  ValueFormat
type PairSet struct {
	pairValueCount   uint16            // Number of PairValueRecords
	PairValueRecords []PairValueRecord `isOpaque:"" subsliceStart:"AtStart"` // [pairValueCount] Array of PairValueRecords, ordered by glyph ID of the second glyph.
}

func (ps *PairSet) customParsePairValueRecords(src []byte, fmt1, fmt2 ValueFormat) (int, error) {
	out := make([]PairValueRecord, ps.pairValueCount)
	offsetR := 2
	var err error
	for i := range out {
		if L := len(src); L < 2+offsetR {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", 2+offsetR, L)
		}
		out[i].SecondGlyph = GlyphID(binary.BigEndian.Uint16(src[offsetR:]))
		out[i].ValueRecord1, offsetR, err = parseValueRecord(fmt1, src, offsetR+2)
		if err != nil {
			return 0, fmt.Errorf("invalid pair set table: %s", err)
		}
		out[i].ValueRecord2, offsetR, err = parseValueRecord(fmt2, src, offsetR)
		if err != nil {
			return 0, fmt.Errorf("invalid pair set table: %s", err)
		}
	}
	return offsetR, nil
}

// binarygen: startOffset=2
type PairPosData2 struct {
	Coverage Coverage `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of PairPos subtable.

	valueFormat1 ValueFormat //	Defines the types of data in valueRecord1 — for the first glyph in the pair (may be zero).
	valueFormat2 ValueFormat //	Defines the types of data in valueRecord2 — for the second glyph in the pair (may be zero).

	classDef1     ClassDef       `offsetSize:"Offset16"` // Offset to ClassDef table, from beginning of PairPos subtable — for the first glyph of the pair.
	classDef2     ClassDef       `offsetSize:"Offset16"` // Offset to ClassDef table, from beginning of PairPos subtable — for the second glyph of the pair.
	class1Count   uint16         //	Number of classes in classDef1 table — includes Class 0.
	class2Count   uint16         //	Number of classes in classDef2 table — includes Class 0.
	class1Records []Class1Record `isOpaque:"" subsliceStart:"AtStart"` //[class1Count]	Array of Class1 records, ordered by classes in classDef1.
}

func (pp *PairPosData2) customParseClass1Records(src []byte) (int, error) {
	const headerSize = 16 // including posFormat and coverageOffset

	pp.class1Records = make([]Class1Record, pp.class1Count)

	offset := headerSize
	for i := range pp.class1Records {
		vi := make([]Class2Record, pp.class2Count)
		for j := range vi {
			var err error
			vi[j].ValueRecord1, offset, err = parseValueRecord(pp.valueFormat1, src, offset)
			if err != nil {
				return 0, err
			}
			vi[j].ValueRecord2, offset, err = parseValueRecord(pp.valueFormat2, src, offset)
			if err != nil {
				return 0, err
			}
		}
		pp.class1Records[i] = Class1Record{Class2Records: vi}
	}
	return offset, nil
}

// DeviceTableHeader is the common header for DeviceTable
// See https://learn.microsoft.com/fr-fr/typography/opentype/spec/chapter2#device-and-variationindex-tables
type DeviceTableHeader struct {
	first       uint16
	second      uint16
	deltaFormat uint16 // Format of deltaValue array data
}

type Anchor interface {
	isAnchor()
}

type EntryExit struct {
	EntryAnchor Anchor
	ExitAnchor  Anchor
}

type CursivePos struct {
	posFormat        uint16            //	Format identifier: format = 1
	coverage         Coverage          `offsetSize:"Offset16"`    //	Offset to Coverage table, from beginning of CursivePos subtable.
	entryExitRecords []entryExitRecord `arrayCount:"FirstUint16"` //[entryExitCount]	Array of EntryExit records, in Coverage index order.
	EntryExits       []EntryExit       `isOpaque:"" subsliceStart:"AtStart"`
}

type entryExitRecord struct {
	entryAnchorOffset Offset16 // Offset to entryAnchor table, from beginning of CursivePos subtable (may be NULL).
	exitAnchorOffset  Offset16 // Offset to exitAnchor table, from beginning of CursivePos subtable (may be NULL).
}

func (cp *CursivePos) customParseEntryExits(src []byte) (int, error) {
	cp.EntryExits = make([]EntryExit, len(cp.entryExitRecords))
	var err error
	for i, rec := range cp.entryExitRecords {
		if rec.entryAnchorOffset != 0 {
			if L := len(src); L < int(rec.entryAnchorOffset) {
				return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.entryAnchorOffset, L)
			}
			cp.EntryExits[i].EntryAnchor, _, err = ParseAnchor(src[rec.entryAnchorOffset:])
			if err != nil {
				return 0, err
			}
		}
		if rec.exitAnchorOffset != 0 {
			if L := len(src); L < int(rec.exitAnchorOffset) {
				return 0, fmt.Errorf("EOF: expected length: %d, got %d", rec.exitAnchorOffset, L)
			}
			cp.EntryExits[i].ExitAnchor, _, err = ParseAnchor(src[rec.exitAnchorOffset:])
			if err != nil {
				return 0, err
			}
		}
	}
	return 0, nil
}

type MarkBasePos struct {
	posFormat      uint16    // Format identifier: format = 1
	markCoverage   Coverage  `offsetSize:"Offset16"` // Offset to markCoverage table, from beginning of MarkBasePos subtable.
	baseCoverage   Coverage  `offsetSize:"Offset16"` // Offset to baseCoverage table, from beginning of MarkBasePos subtable.
	markClassCount uint16    // Number of classes defined for marks
	markArray      MarkArray `offsetSize:"Offset16"`                            // Offset to MarkArray table, from beginning of MarkBasePos subtable.
	baseArray      BaseArray `offsetSize:"Offset16" arguments:"markClassCount"` // Offset to BaseArray table, from beginning of MarkBasePos subtable.
}

type BaseArray struct {
	baseRecords []BaseRecord `arrayCount:"FirstUint16"`
	BaseAnchors [][]Anchor   `isOpaque:"" subsliceStart:"AtStart"`
}

func (ba *BaseArray) customParseBaseAnchors(src []byte, _ int) (int, error) {
	ba.BaseAnchors = make([][]Anchor, len(ba.baseRecords))
	var err error
	for i, rec := range ba.baseRecords {
		bi := make([]Anchor, len(rec.baseAnchorOffsets))
		for j, offset := range rec.baseAnchorOffsets {
			if L := len(src); L < int(offset) {
				return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
			}
			bi[j], _, err = ParseAnchor(src[offset:])
			if err != nil {
				return 0, err
			}
		}
		ba.BaseAnchors[i] = bi
	}
	return len(src), nil
}

type BaseRecord struct {
	baseAnchorOffsets []Offset16 // [markClassCount] // Array of offsets (one per mark class) to Anchor tables. Offsets are from beginning of BaseArray table, ordered by class (offsets may be NULL).
}
