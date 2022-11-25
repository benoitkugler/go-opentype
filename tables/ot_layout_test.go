package tables

import (
	"fmt"
	"reflect"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

func TestParseOTLayout(t *testing.T) {
	for _, filename := range td.WithOTLayout {
		fp := readFontFile(t, filename)
		gsub, _, err := ParseLayout(readTable(t, fp, "GSUB"))
		assertNoErr(t, err)
		assert(t, len(gsub.lookupList.Lookups) == len(gsub.lookupList.Lookups))
		assert(t, len(gsub.lookupList.Lookups) > 0)

		gpos, _, err := ParseLayout(readTable(t, fp, "GPOS"))
		assertNoErr(t, err)
		assert(t, len(gpos.lookupList.Lookups) == len(gpos.lookupList.Lookups))
		assert(t, len(gpos.lookupList.Lookups) > 0)

		_, _, err = ParseGDEF(readTable(t, fp, "GDEF"))
		assertNoErr(t, err)
	}
}

func TestGSUB(t *testing.T) {
	for _, filename := range filenames(t, "toys/gsub") {
		fp := readFontFile(t, filename)
		gsub, _, err := ParseLayout(readTable(t, fp, "GSUB"))
		assertNoErr(t, err)
		assert(t, len(gsub.lookupList.Lookups) == len(gsub.lookupList.Lookups))
		assert(t, len(gsub.lookupList.Lookups) > 0)

		for _, lookup := range gsub.lookupList.Lookups {
			assert(t, lookup.lookupType > 0)
			_, err = lookup.AsGSUBLookups()
			assertNoErr(t, err)
		}
	}
}

func TestGPOS(t *testing.T) {
	for _, filename := range filenames(t, "toys/gpos") {
		fp := readFontFile(t, filename)
		gpos, _, err := ParseLayout(readTable(t, fp, "GPOS"))
		assertNoErr(t, err)
		assert(t, len(gpos.lookupList.Lookups) == len(gpos.lookupList.Lookups))
		assert(t, len(gpos.lookupList.Lookups) > 0)

		for _, lookup := range gpos.lookupList.Lookups {
			assert(t, lookup.lookupType > 0)
			_, err = lookup.AsGPOSLookups()
			assertNoErr(t, err)
		}
	}
}

func TestGSUBIndic(t *testing.T) {
	filepath := "toys/gsub/GSUBChainedContext2.ttf"
	fp := readFontFile(t, filepath)
	gsub, _, err := ParseLayout(readTable(t, fp, "GSUB"))
	assertNoErr(t, err)

	expectedScripts := []Script{
		{
			// Tag: MustNewTag("beng"),
			DefaultLangSys: LangSys{
				RequiredFeatureIndex: 0xFFFF,
				FeatureIndices:       []uint16{0, 2},
			},
			langSysRecords: []tagOffsetRecord{},
			LangSys:        []LangSys{},
		},
		{
			// Tag: MustNewTag("bng2"),
			DefaultLangSys: LangSys{
				RequiredFeatureIndex: 0xFFFF,
				FeatureIndices:       []uint16{1, 2},
			},
			langSysRecords: []tagOffsetRecord{},
			LangSys:        []LangSys{},
		},
	}

	expectedFeatures := []Feature{
		{
			// Tag: MustNewTag("init"),
			LookupListIndices: []uint16{0},
		},
		{
			// Tag: MustNewTag("init"),
			LookupListIndices: []uint16{1},
		},
		{
			// Tag: MustNewTag("blws"),
			LookupListIndices: []uint16{2},
		},
	}

	expectedLoopkup1 := []GSUBLookup{
		SingleSubs{
			Data: SingleSubstData1{
				format:       1,
				Coverage:     Coverage1{1, []GlyphID{6, 7}},
				DeltaGlyphID: 3,
			},
		},
	}

	expectedLoopkup2 := []GSUBLookup{
		SingleSubs{
			Data: SingleSubstData1{
				format:       1,
				Coverage:     Coverage1{1, []GlyphID{6, 7}},
				DeltaGlyphID: 3,
			},
		},
	}

	expectedLoopkup3 := []GSUBLookup{
		ChainedContextualSubs{
			Data: ChainedContextualSubs2{
				format:   2,
				Coverage: Coverage1{1, []GlyphID{5}},
				BacktrackClassDef: ClassDef1{
					format:          1,
					StartGlyphID:    2,
					ClassValueArray: []uint16{1},
				},
				InputClassDef: ClassDef1{
					format:          1,
					StartGlyphID:    5,
					ClassValueArray: []uint16{1},
				},
				LookaheadClassDef: ClassDef2{
					format:            2,
					ClassRangeRecords: []ClassRangeRecord{},
				},
				ChainedClassSeqRuleSet: []ChainedClassSequenceRuleSet{
					{},
					{
						[]ChainedClassSequenceRule{
							{
								BacktrackSequence: []uint16{1},
								inputGlyphCount:   1,
								InputSequence:     []uint16{},
								LookaheadGlyph:    []uint16{},
								SeqLookupRecords: []SequenceLookupRecord{
									{SequenceIndex: 0, LookupListIndex: 3},
								},
							},
						},
					},
				},
			},
		},
	}

	expectedLoopkup4 := []GSUBLookup{
		SingleSubs{
			Data: SingleSubstData1{
				format:       1,
				Coverage:     Coverage1{1, []GlyphID{5}},
				DeltaGlyphID: 6,
			},
		},
	}
	expectedLookups := [][]GSUBLookup{
		expectedLoopkup1, expectedLoopkup2, expectedLoopkup3, expectedLoopkup4,
	}

	if exp, got := expectedScripts, gsub.scriptList.Scripts; !reflect.DeepEqual(exp, got) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
	if exp, got := expectedFeatures, gsub.featureList.Features; !reflect.DeepEqual(exp, got) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
	for i, lk := range gsub.lookupList.Lookups {
		got, err := lk.AsGSUBLookups()
		assertNoErr(t, err)
		exp := expectedLookups[i]
		assertC(t, reflect.DeepEqual(got, exp), fmt.Sprintf("lookup %d expected \n%v\n, got \n%v\n", i, exp, got))
	}
}

func TestGSUBLigature(t *testing.T) {
	filepath := "toys/gsub/GSUBLigature.ttf"
	fp := readFontFile(t, filepath)
	gsub, _, err := ParseLayout(readTable(t, fp, "GSUB"))
	assertNoErr(t, err)

	lookups, err := gsub.lookupList.Lookups[0].AsGSUBLookups()
	assertNoErr(t, err)
	lookup := lookups[0]

	expected := LigatureSubs{
		substFormat: 1,
		Coverage:    Coverage1{1, []GlyphID{3, 4, 7, 8, 9}},
		LigatureSets: []LigatureSet{
			{ // glyph="3"
				[]Ligature{
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{5}, LigatureGlyph: 6},
				},
			},
			{ // glyph="4"
				[]Ligature{
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{4}, LigatureGlyph: 31},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{7}, LigatureGlyph: 32},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{11}, LigatureGlyph: 34},
				},
			},
			{ // glyph="7"
				[]Ligature{
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{4}, LigatureGlyph: 37},
					{componentCount: 3, ComponentGlyphIDs: []GlyphID{7, 7}, LigatureGlyph: 40},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{7}, LigatureGlyph: 8},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{8}, LigatureGlyph: 40},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{11}, LigatureGlyph: 38},
				},
			},
			{ // glyph="8"
				[]Ligature{
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{7}, LigatureGlyph: 40},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{11}, LigatureGlyph: 42},
				},
			},
			{ // glyph="9"
				[]Ligature{
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{4}, LigatureGlyph: 44},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{7}, LigatureGlyph: 45},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{9}, LigatureGlyph: 10},
					{componentCount: 2, ComponentGlyphIDs: []GlyphID{11}, LigatureGlyph: 46},
				},
			},
		},
	}

	assertC(t, reflect.DeepEqual(lookup, expected), fmt.Sprintf("expected %v, got %v", expected, lookup))
}

func TestGPOSCursive(t *testing.T) {
	filepath := "toys/gpos/GPOSCursive.ttf"
	fp := readFontFile(t, filepath)
	gpos, _, err := ParseLayout(readTable(t, fp, "GPOS"))
	assertNoErr(t, err)

	if len(gpos.lookupList.Lookups) != 4 || len(gpos.lookupList.Lookups[0].subtableOffsets) != 1 {
		t.Fatalf("invalid gpos lookups: %v", gpos.lookupList)
	}

	lookups, err := gpos.lookupList.Lookups[0].AsGPOSLookups()
	assertNoErr(t, err)

	cursive, ok := lookups[0].(CursivePos)
	assertC(t, ok, fmt.Sprintf("unexpected type for lookup %T", lookups[0]))

	got := cursive.EntryExits
	expected := []EntryExit{
		{AnchorFormat1{1, 405, 45}, AnchorFormat1{1, 0, 0}},
		{AnchorFormat1{1, 452, 500}, AnchorFormat1{1, 0, 0}},
	}
	assertC(t, reflect.DeepEqual(expected, got), fmt.Sprintf("expected %v, got %v", expected, got))
}

func TestBits(t *testing.T) {
	var buf8 [8]int8
	uint16As2Bits(buf8[:], 0x123F)
	if exp := [8]int8{0, 1, 0, -2, 0, -1, -1, -1}; buf8 != exp {
		t.Fatalf("expected %v, got %v", exp, buf8)
	}

	var buf4 [4]int8
	uint16As4Bits(buf4[:], 0x123F)
	if exp := [4]int8{1, 2, 3, -1}; buf4 != exp {
		t.Fatalf("expected %v, got %v", exp, buf4)
	}

	var buf2 [2]int8
	uint16As8Bits(buf2[:], 0x123F)
	if exp := [2]int8{18, 63}; buf2 != exp {
		t.Fatalf("expected %v, got %v", exp, buf2)
	}
}

func TestGDEFCaretList3(t *testing.T) {
	filepath := "toys/GDEFCaretList3.ttf"
	fp := readFontFile(t, filepath)
	gdef, _, err := ParseGDEF(readTable(t, fp, "GDEF"))
	assertNoErr(t, err)

	expectedLigGlyphs := [][]CaretValue3{ //  LigGlyphCount=4
		// CaretCount=1
		{
			CaretValue3{Coordinate: 620, Device: DeviceVariation{DeltaSetOuter: 3, DeltaSetInner: 205}},
		},
		// CaretCount=1
		{
			CaretValue3{Coordinate: 675, Device: DeviceVariation{DeltaSetOuter: 3, DeltaSetInner: 193}},
		},
		// CaretCount=2
		{
			CaretValue3{Coordinate: 696, Device: DeviceVariation{DeltaSetOuter: 3, DeltaSetInner: 173}},
			CaretValue3{Coordinate: 1351, Device: DeviceVariation{DeltaSetOuter: 6, DeltaSetInner: 14}},
		},
		// CaretCount=2
		{
			CaretValue3{Coordinate: 702, Device: DeviceVariation{DeltaSetOuter: 3, DeltaSetInner: 179}},
			CaretValue3{Coordinate: 1392, Device: DeviceVariation{DeltaSetOuter: 6, DeltaSetInner: 11}},
		},
	}

	for i := range expectedLigGlyphs {
		expL, gotL := expectedLigGlyphs[i], gdef.LigCaretList.LigGlyphs[i]
		assert(t, len(expL) == len(gotL.CaretValues))
		for j := range expL {
			exp, got := expL[j], gotL.CaretValues[j]
			asFormat3, ok := got.(CaretValue3)
			assert(t, ok)
			assert(t, exp.Coordinate == asFormat3.Coordinate)
			assert(t, exp.Device == asFormat3.Device)
		}
	}
}
