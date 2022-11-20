package tables

// The following are types shared by GSUB and GPOS tables

// Coverage specifies all the glyphs affected by a substitution or
// positioning operation described in a subtable.
// Conceptually is it a []GlyphIndex, but it may be implemented for efficiently.
// See https://learn.microsoft.com/typography/opentype/spec/chapter2#lookup-table
type Coverage struct {
	format coverageVersion
	data   coverageData `unionField:"format"`
}

type coverageVersion uint16

const (
	coverageVersion1 coverageVersion = 1
	coverageVersion2 coverageVersion = 2
)

type coverageData interface {
	isCoverage()
}

func (coverageData1) isCoverage() {}
func (coverageData2) isCoverage() {}

type coverageData1 struct {
	glyphs []glyphID `arrayCount:"FirstUint16"`
}

type coverageData2 struct {
	ranges []rangeRecord `arrayCount:"FirstUint16"`
}

type rangeRecord struct {
	startGlyphID       glyphID // First glyph ID in the range
	endGlyphID         glyphID // Last glyph ID in the range
	startCoverageIndex uint16  // Coverage Index of first glyph ID in range
}

type ClassDef struct {
	format classDefVersion
	data   classDefData `unionField:"format"`
}

type classDefVersion uint16

const (
	classDefVersion1 classDefVersion = 1
	classDefVersion2 classDefVersion = 2
)

type classDefData interface {
	isClassDefData()
}

func (classDefData1) isClassDefData() {}
func (classDefData2) isClassDefData() {}

type classDefData1 struct {
	startGlyphID    glyphID  // First glyph ID of the classValueArray
	classValueArray []uint16 `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type classDefData2 struct {
	classRangeRecords []classRangeRecord `arrayCount:"FirstUint16"` //[glyphCount]	Array of Class Values — one per glyph ID
}

type classRangeRecord struct {
	startGlyphID glyphID // First glyph ID in the range
	endGlyphID   glyphID // Last glyph ID in the range
	class        uint16  // Applied to all glyphs in the range
}

// Lookups

type SequenceLookupRecord struct {
	sequenceIndex   uint16 // Index (zero-based) into the input glyph sequence
	lookupListIndex uint16 // Index (zero-based) into the LookupList
}

// binarygen: startOffset=2
// Format identifier: format = 1
type SequenceContextFormat1 struct {
	Coverage   Coverage          `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of SequenceContextFormat1 table
	seqRuleSet []SequenceRuleSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[seqRuleSetCount]	Array of offsets to SequenceRuleSet tables, from beginning of SequenceContextFormat1 table (offsets may be NULL)
}

type SequenceRuleSet struct {
	seqRule []SequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type SequenceRule struct {
	glyphCount       uint16                 // Number of glyphs in the input glyph sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	inputSequence    []glyphID              `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Array of input glyph IDs—starting with the second glyph
	seqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of Sequence lookup records
}

// binarygen: startOffset=2
// Format identifier: format = 2
type SequenceContextFormat2 struct {
	coverage        Coverage               `offsetSize:"Offset16"`                            //	Offset to Coverage table, from beginning of SequenceContextFormat2 table
	classDef        ClassDef               `offsetSize:"Offset16"`                            //	Offset to ClassDef table, from beginning of SequenceContextFormat2 table
	classSeqRuleSet []ClassSequenceRuleSet `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[classSeqRuleSetCount]	Array of offsets to ClassSequenceRuleSet tables, from beginning of SequenceContextFormat2 table (may be NULL)
}

type ClassSequenceRuleSet struct {
	classSeqRule []ClassSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //	Array of offsets to ClassSequenceRule tables, from beginning of ClassSequenceRuleSet table
}

type ClassSequenceRule struct {
	glyphCount       uint16                 // Number of glyphs to be matched
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	inputSequence    []uint16               `arrayCount:"ComputedField-glyphCount-1"`   //[glyphCount - 1]	Sequence of classes to be matched to the input glyph sequence, beginning with the second glyph position
	seqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"` //[seqLookupCount]	Array of SequenceLookupRecords
}

// binarygen: startOffset=2
// Format identifier: format = 2
type SequenceContextFormat3 struct {
	glyphCount       uint16                 // Number of glyphs in the input sequence
	seqLookupCount   uint16                 // Number of SequenceLookupRecords
	coverageOffsets  []Coverage             `arrayCount:"ComputedField-glyphCount" offsetsArray:"Offset16"` //[glyphCount]	Array of offsets to Coverage tables, from beginning of SequenceContextFormat3 subtable
	seqLookupRecords []SequenceLookupRecord `arrayCount:"ComputedField-seqLookupCount"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

// Format identifier: format = 1
type ChainedSequenceContextFormat1 struct {
	coverageOffset    Coverage                 `offsetSize:"Offset16"`                             //	Offset to Coverage table, from beginning of ChainSequenceContextFormat1 table
	chainedSeqRuleSet []ChainedSequenceRuleSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[chainedSeqRuleSetCount]	Array of offsets to ChainedSeqRuleSet tables, from beginning of ChainedSequenceContextFormat1 table (may be NULL)
}

type ChainedSequenceRuleSet struct {
	chainedSeqRules []ChainedSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to SequenceRule tables, from beginning of the SequenceRuleSet table
}

type ChainedSequenceRule struct {
	backtrackSequence []glyphID              `arrayCount:"FirstUint16"` //[backtrackGlyphCount]	Array of backtrack glyph IDs
	inputGlyphCount   uint16                 //	Number of glyphs in the input sequence
	inputSequence     []glyphID              `arrayCount:"ComputedField-inputGlyphCount-1"` //[inputGlyphCount - 1]	Array of input glyph IDs—start with second glyph
	lookaheadSequence glyphID                `arrayCount:"FirstUint16"`                     //[lookaheadGlyphCount]	Array of lookahead glyph IDs
	seqLookupRecords  []SequenceLookupRecord `arrayCount:"FirstUint16"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

// binarygen: startOffset=2
// Format identifier: format = 2
type ChainedSequenceContextFormat2 struct {
	coverage               Coverage                      `offsetSize:"Offset16"`                            // Offset to Coverage table, from beginning of ChainedSequenceContextFormat2 table
	backtrackClassDef      ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing backtrack sequence context, from beginning of ChainedSequenceContextFormat2 table
	inputClassDef          ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing input sequence context, from beginning of ChainedSequenceContextFormat2 table
	lookaheadClassDef      ClassDef                      `offsetSize:"Offset16"`                            // Offset to ClassDef table containing lookahead sequence context, from beginning of ChainedSequenceContextFormat2 table
	chainedClassSeqRuleSet []ChainedClassSequenceRuleSet `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[chainedClassSeqRuleSetCount]	Array of offsets to ChainedClassSequenceRuleSet tables, from beginning of ChainedSequenceContextFormat2 table (may be NULL)
}

type ChainedClassSequenceRuleSet struct {
	chainedClassSeqRules []ChainedClassSequenceRule `arrayCount:"FirstUint16" offsetsArray:"Offset16"` // Array of offsets to ChainedClassSequenceRule tables, from beginning of ChainedClassSequenceRuleSet
}

type ChainedClassSequenceRule struct {
	backtrackSequence   []uint16               `arrayCount:"FirstUint16"` //[backtrackGlyphCount]	Array of backtrack-sequence classes
	inputGlyphCount     uint16                 //	Total number of glyphs in the input sequence
	inputSequence       []uint16               `arrayCount:"ComputedField-inputGlyphCount-1"` //[inputGlyphCount - 1]	Array of input sequence classes, beginning with the second glyph position
	lookaheadGlyphCount []uint16               `arrayCount:"FirstUint16"`                     //[lookaheadGlyphCount]	Array of lookahead-sequence classes
	seqLookupRecords    []SequenceLookupRecord `arrayCount:"FirstUint16"`                     //[seqLookupCount]	Array of SequenceLookupRecords
}

// binarygen: startOffset=2
// Format identifier: format = 3
type ChainedSequenceContextFormat3 struct {
	backtrackCoverageOffsets []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[backtrackGlyphCount]	Array of offsets to coverage tables for the backtrack sequence
	inputCoverageOffsets     []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[inputGlyphCount]	Array of offsets to coverage tables for the input sequence
	lookaheadCoverageOffsets []Coverage             `arrayCount:"FirstUint16" offsetsArray:"Offset16"` //[lookaheadGlyphCount]	Array of offsets to coverage tables for the lookahead sequence
	seqLookupRecords         []SequenceLookupRecord `arrayCount:"FirstUint16"`                         //[seqLookupCount]	Array of SequenceLookupRecords
}
