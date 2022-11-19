package tables

import "fmt"

type SingleSubstitution struct {
	substFormat singleSubsVersion
	data        SingleSubstData `unionField:"substFormat" subsliceStart:"AtStart"`
}

type singleSubsVersion uint16

const (
	singleSubsVersion1 singleSubsVersion = 1
	singleSubsVersion2 singleSubsVersion = 2
)

type SingleSubstData interface {
	isSingleSubstData()
}

func (SingleSubstData1) isSingleSubstData() {}
func (SingleSubstData2) isSingleSubstData() {}

// binarygen: startOffset=2
type SingleSubstData1 struct {
	coverage     Coverage `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of substitution subtable
	deltaGlyphID int16    // Add to original glyph ID to get substitute glyph ID
}

// binarygen: startOffset=2
type SingleSubstData2 struct {
	coverage           Coverage  `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of substitution subtable
	substituteGlyphIDs []glyphID `arrayCount:"FirstUint16"` //[glyphCount]	Array of substitute glyph IDs — ordered by Coverage index
}

type MultipleSubstitution struct {
	substFormat     uint16     // Format identifier: format = 1
	coverageOffset  Coverage   `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of substitution subtable
	sequenceOffsets []Offset16 `arrayCount:"FirstUint16"` //[sequenceCount]	Array of offsets to Sequence tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	sequences       []Sequence `isOpaque:"" subsliceStart:"AtStart"`
}

func (lk MultipleSubstitution) customParseSequences(src []byte) (int, error) {
	var err error
	lk.sequences = make([]Sequence, len(lk.sequenceOffsets))
	for i, offset := range lk.sequenceOffsets {
		if L := len(src); L < int(offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		lk.sequences[i], _, err = ParseSequence(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type Sequence struct {
	substituteGlyphIDs []glyphID `arrayCount:"FirstUint16"` // [glyphCount]	String of glyph IDs to substitute
}

type AlternateSubstitution struct {
	substFormat         uint16         //	Format identifier: format = 1
	coverageOffset      Coverage       `offsetSize:"Offset16"`    //	Offset to Coverage table, from beginning of substitution subtable
	alternateSetOffsets []Offset16     `arrayCount:"FirstUint16"` //[alternateSetCount]	Array of offsets to AlternateSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	alternateSets       []AlternateSet `isOpaque:"" subsliceStart:"AtStart"`
}

func (lk AlternateSubstitution) customParseAlternateSets(src []byte) (int, error) {
	var err error
	lk.alternateSets = make([]AlternateSet, len(lk.alternateSetOffsets))
	for i, offset := range lk.alternateSetOffsets {
		if L := len(src); L < int(offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		lk.alternateSets[i], _, err = ParseAlternateSet(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

type AlternateSet struct {
	alternateGlyphIDs []glyphID `arrayCount:"FirstUint16"` // Array of alternate glyph IDs, in arbitrary order
}

type LigatureSubstitution struct {
	substFormat        uint16        // Format identifier: format = 1
	coverage           Coverage      `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of substitution subtable
	ligatureSetOffsets []Offset16    `arrayCount:"FirstUint16"` //[ligatureSetCount]	Array of offsets to LigatureSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
	ligatureSets       []LigatureSet `isOpaque:"" subsliceStart:"AtStart"`
}

func (ls LigatureSubstitution) customParseLigatureSets(src []byte) (int, error) {
	var err error
	ls.ligatureSets = make([]LigatureSet, len(ls.ligatureSetOffsets))
	for i, offset := range ls.ligatureSetOffsets {
		if L := len(src); L < int(offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		ls.ligatureSets[i], _, err = ParseLigatureSet(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

// All ligatures beginning with the same glyph
type LigatureSet struct {
	ligatureOffsets []Offset16 `arrayCount:"FirstUint16"` // [LigatureCount]	Array of offsets to Ligature tables. Offsets are from beginning of LigatureSet table, ordered by preference.
	Ligatures       []Ligature `isOpaque:"" subsliceStart:"AtStart"`
}

func (ls LigatureSet) customParseLigatures(src []byte) (int, error) {
	var err error
	ls.Ligatures = make([]Ligature, len(ls.ligatureOffsets))
	for i, offset := range ls.ligatureOffsets {
		if L := len(src); L < int(offset) {
			return 0, fmt.Errorf("EOF: expected length: %d, got %d", offset, L)
		}
		ls.Ligatures[i], _, err = ParseLigature(src[offset:])
		if err != nil {
			return 0, err
		}
	}
	return len(src), nil
}

// Glyph components for one ligature
type Ligature struct {
	ligatureGlyph     glyphID   //	glyph ID of ligature to substitute
	componentCount    uint16    //	Number of components in the ligature
	componentGlyphIDs []glyphID `arrayCount:"ComputedField-componentCount-1"` //  [componentCount - 1]	Array of component glyph IDs — start with the second component, ordered in writing direction
}

type ContextualSubstitution struct {
	format contextualSubsVersion
	data   ContextualSubs `unionField:"format" subsliceStart:"AtStart"`
}

type contextualSubsVersion uint16

const (
	contextualSubsVersion1 contextualSubsVersion = 1
	contextualSubsVersion2 contextualSubsVersion = 2
	contextualSubsVersion3 contextualSubsVersion = 3
)

type ContextualSubs interface {
	isContextualSubs()
}

// binarygen: startOffset=2
type ContextualSubs1 SequenceContextFormat1

// binarygen: startOffset=2
type ContextualSubs2 SequenceContextFormat2

// binarygen: startOffset=2
type ContextualSubs3 SequenceContextFormat3

func (ContextualSubs1) isContextualSubs() {}
func (ContextualSubs2) isContextualSubs() {}
func (ContextualSubs3) isContextualSubs() {}

type ChainedContextualSubstitution struct {
	format chainedContextualSubsVersion
	data   ChainedContextualSubs `unionField:"format" subsliceStart:"AtStart"`
}

type chainedContextualSubsVersion uint16

const (
	chainedContextualSubsVersion1 chainedContextualSubsVersion = 1
	chainedContextualSubsVersion2 chainedContextualSubsVersion = 2
	chainedContextualSubsVersion3 chainedContextualSubsVersion = 3
)

type ChainedContextualSubs interface {
	isChainedContextualSubs()
}

// binarygen: startOffset=2
type ChainedContextualSubs1 ChainedSequenceContextFormat1

// binarygen: startOffset=2
type ChainedContextualSubs2 ChainedSequenceContextFormat2

// binarygen: startOffset=2
type ChainedContextualSubs3 ChainedSequenceContextFormat3

func (ChainedContextualSubs1) isChainedContextualSubs() {}
func (ChainedContextualSubs2) isChainedContextualSubs() {}
func (ChainedContextualSubs3) isChainedContextualSubs() {}

type ExtensionSubstitution struct {
	substFormat         uint16   //	Format identifier. Set to 1.
	extensionLookupType uint16   //	Lookup type of subtable referenced by extensionOffset (that is, the extension subtable).
	extensionOffset     Offset32 //	Offset to the extension subtable, of lookup type extensionLookupType, relative to the start of the ExtensionSubstFormat1 subtable.
	rawData             []byte   `subsliceStart:"AtStart" arrayCount:"ToEnd"`
}

type ReverseChainSingleSubstitution struct {
	substFormat              uint16     // Format identifier: format = 1
	coverage                 Coverage   `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of substitution subtable.
	backtrackCoverageOffsets []Offset16 `arrayCount:"FirstUint16"` //[backtrackGlyphCount]	Array of offsets to coverage tables in backtrack sequence, in glyph sequence order.
	lookaheadCoverageOffsets []Offset16 `arrayCount:"FirstUint16"` //[lookaheadGlyphCount]	Array of offsets to coverage tables in lookahead sequence, in glyph sequence order.
	substituteGlyphIDs       []glyphID  `arrayCount:"FirstUint16"` //[glyphCount]	Array of substitute glyph IDs — ordered by Coverage index.
}
