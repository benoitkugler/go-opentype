package tables

// Morx is the extended glyph metamorphosis table
// See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6morx.html
type Morx struct {
	version uint16      // Version number of the extended glyph metamorphosis table (either 2 or 3)
	unused  uint16      // Set to 0
	nChains uint32      // Number of metamorphosis chains contained in this table.
	Chains  []MorxChain `arrayCount:"ComputedField-nChains"`
}

// MorxChain is a set of subtables
type MorxChain struct {
	Flags           uint32              // The default specification for subtables.
	chainLength     uint32              // Total byte count, including this header; must be a multiple of 4.
	nFeatureEntries uint32              // Number of feature subtable entries.
	nSubtable       uint32              // The number of subtables in the chain.
	Features        []AATFeature        `arrayCount:"ComputedField-nFeatureEntries"`
	Subtables       []MorxChainSubtable `arrayCount:"ComputedField-nSubtable"`
}

type AATFeature struct {
	FeatureType    uint16
	FeatureSetting uint16
	EnableFlags    uint32 // Flags for the settings that this feature and setting enables.
	DisableFlags   uint32 // Complement of flags for the settings that this feature and setting disable.
}

type MorxChainSubtable struct {
	length uint32 // Total subtable length, including this header.

	// Coverage flags and subtable type.
	coverage byte
	ignored  [2]byte
	version  MorxSubtableVersion

	subFeatureFlags uint32 // The 32-bit mask identifying which subtable this is (the subtable being executed if the AND of this value and the processed defaultFlags is nonzero)

	Content     MorxSubtable `unionField:"version"`
	postProcess [0]byte      `isOpaque:""`
}

func (mc *MorxChainSubtable) parsePostProcess(src []byte, _ int) (int, error) {
	return int(mc.length), nil
}

// MorxSubtableVersion indicates the kind of 'morx' subtable.
// See the constants.
type MorxSubtableVersion uint8

const (
	MorxSubtableVersionRearrangement MorxSubtableVersion = iota
	MorxSubtableVersionContextual
	MorxSubtableVersionLigature
	_ // reserved
	MorxSubtableVersionNonContextual
	MorxSubtableVersionInsertion
)

type MorxSubtable interface {
	isMorxSubtable()
}

func (MorxSubtableRearrangement) isMorxSubtable() {}
func (MorxSubtableContextual) isMorxSubtable()    {}
func (MorxSubtableLigature) isMorxSubtable()      {}
func (MorxSubtableNonContextual) isMorxSubtable() {}
func (MorxSubtableInsertion) isMorxSubtable()     {}

// binarygen: argument=valuesCount int
type MorxSubtableRearrangement struct {
	aatSTXHeader `arguments:"valuesCount, 0"`
}

// binarygen: argument=valuesCount int
type MorxSubtableContextual struct {
	aatSTXHeader `arguments:"valuesCount, 4"`
	// Byte offset from the beginning of the state subtable to the beginning of the substitution tables :
	// each value of the array is itself an offet to a aatLookupTable, and the number of
	// items is computed from the header
	Substitutions SubstitutionsTable `offsetSize:"Offset32" arguments:".nSubs(), valuesCount"`
}

type SubstitutionsTable struct {
	Substitutions []AATLookup `offsetsArray:"Offset32"`
}

func (ct *MorxSubtableContextual) nSubs() int {
	// find the maximum index need in the substitution array
	var maxi uint16
	for _, entry := range ct.Entries {
		markIndex, currentIndex := entry.AsMorxContextual()
		if markIndex != 0xFFFF && markIndex > maxi {
			maxi = markIndex
		}
		if currentIndex != 0xFFFF && currentIndex > maxi {
			maxi = currentIndex
		}
	}
	return int(maxi) + 1
}

// binarygen: argument=valuesCount int
type MorxSubtableLigature struct {
	aatSTXHeader    `arguments:"valuesCount, 2"`
	ligActionOffset uint32 // Byte offset from stateHeader to the start of the ligature action table.
	componentOffset uint32 // Byte offset from stateHeader to the start of the component table.
	ligatureOffset  uint32 // Byte offset from stateHeader to the start of the actual ligature lists.
	rawData         []byte `arrayCount:"ToEnd"`
}

type MorxSubtableNonContextual struct {
	table AATLookup
}

// binarygen: argument=valuesCount int
type MorxSubtableInsertion struct {
	aatSTXHeader          `arguments:"valuesCount, 4"`
	insertionActionOffset uint32 //	Byte offset from stateHeader to the start of the insertion glyph table.
	rawData               []byte `arrayCount:"ToEnd"`
}
