package tables

type SingleSubstitution struct {
	substFormat singleSubsVersion
	Data        SingleSubstData `unionField:"substFormat" subsliceStart:"AtStart"`
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
	Coverage     Coverage `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of substitution subtable
	DeltaGlyphID int16    // Add to original glyph ID to get substitute glyph ID
}

// binarygen: startOffset=2
type SingleSubstData2 struct {
	Coverage           Coverage  `offsetSize:"Offset16"`    // Offset to Coverage table, from beginning of substitution subtable
	SubstituteGlyphIDs []GlyphID `arrayCount:"FirstUint16"` //[glyphCount]	Array of substitute glyph IDs — ordered by Coverage index
}

type MultipleSubstitution struct {
	substFormat    uint16     // Format identifier: format = 1
	CoverageOffset Coverage   `offsetSize:"Offset16"` // Offset to Coverage table, from beginning of substitution subtable
	Sequences      []Sequence `arrayCount:"FirstUint16"  offsetsArray:"Offset16"`
	//[sequenceCount]	Array of offsets to Sequence tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
}

type Sequence struct {
	SubstituteGlyphIDs []GlyphID `arrayCount:"FirstUint16"` // [glyphCount]	String of glyph IDs to substitute
}

type AlternateSubstitution struct {
	substFormat    uint16         //	Format identifier: format = 1
	CoverageOffset Coverage       `offsetSize:"Offset16"` //	Offset to Coverage table, from beginning of substitution subtable
	AlternateSets  []AlternateSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"`
}

type AlternateSet struct {
	AlternateGlyphIDs []GlyphID `arrayCount:"FirstUint16"` // Array of alternate glyph IDs, in arbitrary order
}

type LigatureSubstitution struct {
	substFormat  uint16        // Format identifier: format = 1
	Coverage     Coverage      `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of substitution subtable
	LigatureSets []LigatureSet `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[ligatureSetCount]	Array of offsets to LigatureSet tables. Offsets are from beginning of substitution subtable, ordered by Coverage index
}

// All ligatures beginning with the same glyph
type LigatureSet struct {
	Ligatures []Ligature `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` // [LigatureCount]	Array of offsets to Ligature tables. Offsets are from beginning of LigatureSet table, ordered by preference.
}

// Glyph components for one ligature
type Ligature struct {
	LigatureGlyph     GlyphID   //	glyph ID of ligature to substitute
	componentCount    uint16    //	Number of components in the ligature
	ComponentGlyphIDs []GlyphID `arrayCount:"ComputedField-componentCount-1"` //  [componentCount - 1]	Array of component glyph IDs — start with the second component, ordered in writing direction
}

type ContextualSubstitution struct {
	format contextualSubsVersion
	Data   ContextualSubs `unionField:"format" subsliceStart:"AtStart"`
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
	Data   ChainedContextualSubs `unionField:"format" subsliceStart:"AtStart"`
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

type ExtensionSubs Extension

type ReverseChainSingleSubstitution struct {
	substFormat              uint16     // Format identifier: format = 1
	Coverage                 Coverage   `offsetSize:"Offset16"`                             // Offset to Coverage table, from beginning of substitution subtable.
	BacktrackCoverageOffsets []Coverage `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[backtrackGlyphCount]	Array of offsets to coverage tables in backtrack sequence, in glyph sequence order.
	LookaheadCoverageOffsets []Coverage `arrayCount:"FirstUint16"  offsetsArray:"Offset16"` //[lookaheadGlyphCount]	Array of offsets to coverage tables in lookahead sequence, in glyph sequence order.
	SubstituteGlyphIDs       []GlyphID  `arrayCount:"FirstUint16"`                          //[glyphCount]	Array of substitute glyph IDs — ordered by Coverage index.
}
