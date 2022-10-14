package tables

//go:generate ../../binarygen/cmd/generator aat_mortx.go

// Morx is the extended glyph metamorphosis table
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6morx.html
type Morx struct {
	version uint16      // Version number of the extended glyph metamorphosis table (either 2 or 3)
	unused  uint16      // Set to 0
	nChains uint32      // Number of metamorphosis chains contained in this table.
	chains  []morxChain `len:"nChains"`
}

// mortx chain : set of subtables
type morxChain struct {
	flags           uint32
	chainLength     uint32
	nFeatureEntries uint32
	nSubtable       uint32
	features        []aatFeature `len:"nFeatureEntries"`
	subtablesData   []byte       `len:"__to_chainLength"`
}

type aatFeature struct {
	featureType    uint16
	featureSetting uint16
	enableFlags    uint32 // Flags for the settings that this feature and setting enables.
	disableFlags   uint32 // Complement of flags for the settings that this feature and setting disable.
}

type MorxChainSubtable struct {
	length uint32 // Total subtable length, including this header.

	// Coverage flags and subtable type.
	coverage byte
	ignored  [2]byte
	version  morxSubtableVersion

	subFeatureFlags uint32 // The 32-bit mask identifying which subtable this is (the subtable being executed if the AND of this value and the processed defaultFlags is nonzero)

	subtableContent morxSubtable `version-field:"version"`
}

// morxSubtableVersion indicates the kind of 'morx' subtable.
// See the constants.
type morxSubtableVersion uint8

const (
	morxSubtableVersionRearrangement morxSubtableVersion = iota
	morxSubtableVersionContextual
	morxSubtableVersionLigature
	_ // reserved
	morxSubtableVersionNonContextual
	morxSubtableVersionInsertion
)

type morxSubtable interface {
	isMorxSubtable()
}

func (morxSubtableRearrangement) isMorxSubtable() {}
func (morxSubtableContextual) isMorxSubtable()    {}
func (morxSubtableLigature) isMorxSubtable()      {}
func (morxSubtableNonContextual) isMorxSubtable() {}
func (morxSubtableInsertion) isMorxSubtable()     {}

type morxSubtableRearrangement struct {
	aatSTXHeader
	rawData []byte `len:"__toEnd"`
}

type morxSubtableContextual struct {
	aatSTXHeader
	// Byte offset from the beginning of the state subtable to the beginning of the substitution tables :
	// each value of the array is itself an offet to a aatLookupTable, and the number of
	// items is computed from the header
	substitutionTableOffset uint32 `offset-size:"uint32"`
	rawData                 []byte `len:"__toEnd"`
}

type morxSubtableLigature struct {
	aatSTXHeader
	ligActionOffset uint32 // Byte offset from stateHeader to the start of the ligature action table.
	componentOffset uint32 // Byte offset from stateHeader to the start of the component table.
	ligatureOffset  uint32 // Byte offset from stateHeader to the start of the actual ligature lists.
	rawData         []byte `len:"__toEnd"`
}

type morxSubtableNonContextual struct {
	table aatLookup
}

type morxSubtableInsertion struct {
	aatSTXHeader
	insertionActionOffset uint32 //	Byte offset from stateHeader to the start of the insertion glyph table.
	rawData               []byte `len:"__toEnd"`
}
