package tables

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

func deHexStr(s string) []byte {
	s = strings.Join(strings.Split(s, " "), "")
	if len(s)%2 != 0 {
		s += "0"
	}
	out, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return out
}

func TestParseTuple(t *testing.T) {
	// imported from fonttools

	data := deHexStr("DE AD C0 00 20 00 DE AD")
	got, _, err := ParseTuple(data[2:], 2)
	assertNoErr(t, err)
	expected := Tuple{Values: []float32{-1, 0.5}}
	assertC(t, reflect.DeepEqual(got, expected), fmt.Sprintf("%v != %v", got, expected))

	// Shared tuples in the 'gvar' table of the Skia font, as printed
	// in Apple's TrueType specification.
	// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6gvar.html
	skiaGvarSharedTuplesData := deHexStr(
		"40 00 00 00 C0 00 00 00 00 00 40 00 00 00 C0 00 " +
			"C0 00 C0 00 40 00 C0 00 40 00 40 00 C0 00 40 00")

	skiaGvarSharedTuples := SharedTuples{[]Tuple{
		{Values: []Float214{1.0, 0.0}},
		{Values: []Float214{-1.0, 0.0}},
		{Values: []Float214{0.0, 1.0}},
		{Values: []Float214{0.0, -1.0}},
		{Values: []Float214{-1.0, -1.0}},
		{Values: []Float214{1.0, -1.0}},
		{Values: []Float214{1.0, 1.0}},
		{Values: []Float214{-1.0, 1.0}},
	}}
	sharedTuples, _, err := ParseSharedTuples(skiaGvarSharedTuplesData, 8, 2)
	assertNoErr(t, err)
	assert(t, reflect.DeepEqual(sharedTuples, skiaGvarSharedTuples))
}

func TestParseGvar(t *testing.T) {
	// imported from fonttools

	gvarData := deHexStr("0001 0000 " + //   0: majorVersion=1 minorVersion=0
		"0002 0000 " + //   4: axisCount=2 sharedTupleCount=0
		"0000001C " + //   8: offsetToSharedTuples=28
		"0003 0000 " + //  12: glyphCount=3 flags=0
		"0000001C " + //  16: offsetToGlyphVariationData=28
		"0000 0000 000C 002F " + //  20: offsets=[0,0,12,47], times 2: [0,0,24,94],
		//                 //           +offsetToGlyphVariationData: [28,28,52,122]
		//
		// 28: Glyph variation data for glyph //0, ".notdef"
		// ------------------------------------------------
		// (no variation data for this glyph)
		//
		// 28: Glyph variation data for glyph //1, "space"
		// ----------------------------------------------
		"8001 000C " + //  28: tupleVariationCount=1|TUPLES_SHARE_POINT_NUMBERS, offsetToData=12(+28=40)
		"000A " + //  32: tvHeader[0].variationDataSize=10
		"8000 " + //  34: tvHeader[0].tupleIndex=EMBEDDED_PEAK
		"0000 2CCD " + //  36: tvHeader[0].peakTuple={wght:0.0, wdth:0.7}
		"00 " + //  40: all points
		"03 01 02 03 04 " + //  41: deltaX=[1, 2, 3, 4]
		"03 0b 16 21 2C " + //  46: deltaY=[11, 22, 33, 44]
		"00 " + //  51: padding
		//
		// 52: Glyph variation data for glyph //2, "I"
		// ------------------------------------------
		"8002 001c " + //  52: tupleVariationCount=2|TUPLES_SHARE_POINT_NUMBERS, offsetToData=28(+52=80)
		"0012 " + //  56: tvHeader[0].variationDataSize=18
		"C000 " + //  58: tvHeader[0].tupleIndex=EMBEDDED_PEAK|INTERMEDIATE_REGION
		"2000 0000 " + //  60: tvHeader[0].peakTuple={wght:0.5, wdth:0.0}
		"0000 0000 " + //  64: tvHeader[0].intermediateStart={wght:0.0, wdth:0.0}
		"4000 0000 " + //  68: tvHeader[0].intermediateEnd={wght:1.0, wdth:0.0}
		"0016 " + //  72: tvHeader[1].variationDataSize=22
		"A000 " + //  74: tvHeader[1].tupleIndex=EMBEDDED_PEAK|PRIVATE_POINTS
		"C000 3333 " + //  76: tvHeader[1].peakTuple={wght:-1.0, wdth:0.8}
		"00 " + //  80: all points
		"07 03 01 04 01 " + //  81: deltaX.len=7, deltaX=[3, 1, 4, 1,
		"05 09 02 06 " + //  86:                       5, 9, 2, 6]
		"07 03 01 04 01 " + //  90: deltaY.len=7, deltaY=[3, 1, 4, 1,
		"05 09 02 06 " + //  95:                       5, 9, 2, 6]
		"06 " + //  99: 6 points
		"05 00 01 03 01 " + // 100: runLen=5(+1=6); delta-encoded run=[0, 1, 4, 5,
		"01 01 " + // 105:                                    6, 7]
		"05 f8 07 fc 03 fe 01 " + // 107: deltaX.len=5, deltaX=[-8,7,-4,3,-2,1]
		"05 a8 4d 2c 21 ea 0b " + // 114: deltaY.len=5, deltaY=[-88,77,44,33,-22,11]
		"00") // 121: padding
	assert(t, len(gvarData) == 122)

	gvarDataEmptyVariations := deHexStr("0001 0000 " + //  0: majorVersion=1 minorVersion=0
		"0002 0000 " + //  4: axisCount=2 sharedTupleCount=0
		"0000001c " + //  8: offsetToSharedTuples=28
		"0003 0000 " + // 12: glyphCount=3 flags=0
		"0000001c " + // 16: offsetToGlyphVariationData=28
		"0000 0000 0000 0000") // 20: offsets=[0, 0, 0, 0]
	assert(t, len(gvarDataEmptyVariations) == 28)

	sharedTuplesExpected := SharedTuples{}
	variationsHeadersExpected := [][]TupleVariationHeader{
		0: {},
		1: {
			{
				variationDataSize: 0x000A,
				tupleIndex:        0x8000,
				peakTuple:         Tuple{[]float32{0, 0.7000122}},
			},
		},
		2: {
			{
				variationDataSize: 0x0012,
				tupleIndex:        0xC000,
				peakTuple:         Tuple{[]float32{0.5, 0}},
				intermediateTuples: [2]Tuple{
					{[]float32{0, 0}},
					{[]float32{1, 0}},
				},
			},
			{
				variationDataSize: 0x0016,
				tupleIndex:        0xA000,
				peakTuple:         Tuple{[]float32{-1, 0.7999878}},
			},
		},
	}

	gvarEmptyVariationsExpected := make([][]TupleVariationHeader, 3)

	out, _, err := ParseGvar(gvarData)
	assertNoErr(t, err)
	assert(t, reflect.DeepEqual(sharedTuplesExpected, out.SharedTuples))
	assert(t, len(variationsHeadersExpected) == len(out.GlyphVariationDatas))
	for i, exp := range variationsHeadersExpected {
		got := out.GlyphVariationDatas[i].TupleVariationHeaders
		assertC(t, fmt.Sprintf("%v", exp) == fmt.Sprintf("%v", got), fmt.Sprintf("%v != %v", exp, got))
	}

	out, _, err = ParseGvar(gvarDataEmptyVariations)
	assertNoErr(t, err)
	assert(t, len(gvarEmptyVariationsExpected) == len(out.GlyphVariationDatas))
	for i, exp := range gvarEmptyVariationsExpected {
		assert(t, reflect.DeepEqual(exp, out.GlyphVariationDatas[i].TupleVariationHeaders))
	}
}

func TestParseGvar2(t *testing.T) {
	for _, filepath := range []string{
		"common/Commissioner-VF.ttf",
		"common/Mada-VF.ttf",
	} {
		fp := readFontFile(t, filepath)
		_, _, err := ParseGvar(readTable(t, fp, "gvar"))
		assertNoErr(t, err)
	}
}

func TestParseHvar(t *testing.T) {
	for _, filepath := range []string{
		"common/Commissioner-VF.ttf",
		"common/Selawik-VF.ttf",
	} {
		fp := readFontFile(t, filepath)
		_, _, err := ParseHVAR(readTable(t, fp, "HVAR"))
		assertNoErr(t, err)
	}
}

func TestParseAvar(t *testing.T) {
	for _, filepath := range td.WithAvar {
		fp := readFontFile(t, filepath)
		_, _, err := ParseAvar(readTable(t, fp, "avar"))
		assertNoErr(t, err)
	}
}

func TestParseMVAR(t *testing.T) {
	for _, filepath := range td.WithMVAR {
		fp := readFontFile(t, filepath)
		_, _, err := ParseMVAR(readTable(t, fp, "MVAR"))
		assertNoErr(t, err)
	}
}

func TestParseFvar(t *testing.T) {
	for _, item := range td.WithFvar {
		fp := readFontFile(t, item.Path)
		fvar, _, err := ParseFvar(readTable(t, fp, "fvar"))
		assertNoErr(t, err)
		assert(t, len(fvar.Axis) == item.AxisCount)
	}
}
