package tables

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/opentype"
)

func TestParseGlyf(t *testing.T) {
	for _, file := range td.WithGlyphs {
		fp := readFontFile(t, file.Path)
		head, _, err := ParseHead(readTable(t, fp, "head"))
		assertNoErr(t, err)

		maxp, _, err := ParseMaxp(readTable(t, fp, "maxp"))
		assertNoErr(t, err)

		loca, err := ParseLoca(readTable(t, fp, "loca"), int(maxp.numGlyphs), head.IndexToLocFormat == 1)
		assertNoErr(t, err)

		glyphs, err := ParseGlyf(readTable(t, fp, "glyf"), loca)
		assertNoErr(t, err)
		assert(t, len(glyphs) == len(loca)-1)
		assert(t, len(glyphs) == file.GlyphNumber)
	}
}

func assertGlyphSizeEqual(t *testing.T, g1, g2 Glyph) {
	assert(t, g1.XMin == g2.XMin)
	assert(t, g1.YMin == g2.YMin)
	assert(t, g1.XMax == g2.XMax)
	assert(t, g1.YMax == g2.YMax)
}

// do not compare flags
func assertPointEqual(t *testing.T, exp, got GlyphContourPoint) {
	assertC(t, exp.X == got.X && exp.Y == got.Y,
		fmt.Sprintf("expected contour point (%d,%d), got (%d,%d)", exp.X, exp.Y, got.X, got.Y))
}

// do not compare flags
func assertCompositeEqual(t *testing.T, exp, got CompositeGlyphPart) {
	exp.Flags, got.Flags = 0, 0
	assertC(t, exp == got, fmt.Sprintf("expected composite part %v, got %v", exp, got))
}

func TestGlyphCoordinates(t *testing.T) {
	// imported from fonttools
	glyfBin := []byte{0x0, 0x2, 0x0, 0xc, 0x0, 0x0, 0x4, 0x94, 0x5, 0x96, 0x0, 0x6, 0x0, 0xa, 0x0, 0x0, 0x41, 0x33, 0x1, 0x23, 0x1, 0x1, 0x23, 0x13, 0x21, 0x15, 0x21, 0x1, 0xf5, 0xb6, 0x1, 0xe9, 0xbd, 0xfe, 0x78, 0xfe, 0x78, 0xbb, 0xed, 0x2, 0xae, 0xfd, 0x52, 0x5, 0x96, 0xfa, 0x6a, 0x4, 0xa9, 0xfb, 0x57, 0x2, 0x2, 0xa2}
	headBin := []byte{0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x6e, 0x4f, 0x1c, 0xcf, 0x5f, 0xf, 0x3c, 0xf5, 0x20, 0x1b, 0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd7, 0xc0, 0x23, 0x1c, 0x0, 0x0, 0x0, 0x0, 0xd7, 0xc1, 0x2e, 0xf9, 0x0, 0xc, 0x0, 0x0, 0x4, 0x94, 0x5, 0x96, 0x0, 0x0, 0x0, 0x9, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0}
	locaBin := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b}
	maxpBin := []byte{0x0, 0x1, 0x0, 0x0, 0x0, 0x4, 0x0, 0xb, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x2, 0x0, 0x1, 0x61, 0x0, 0x0, 0x0, 0x0}

	num, _, err := ParseMaxp(maxpBin)
	assertNoErr(t, err)
	assert(t, num.numGlyphs == 4)

	head, _, err := ParseHead(headBin)
	assertNoErr(t, err)

	loca, err := ParseLoca(locaBin, int(num.numGlyphs), head.IndexToLocFormat == 1)
	assertNoErr(t, err)

	glyphs, err := ParseGlyf(glyfBin, loca)
	assertNoErr(t, err)
	assert(t, len(glyphs) == 4)
	assert(t, glyphs[0].Data == nil)
	assert(t, glyphs[1].Data == nil)
	assert(t, glyphs[2].Data == nil)

	glyph := glyphs[3]
	assertGlyphSizeEqual(t, glyph, Glyph{XMin: 12, YMin: 0, XMax: 1172, YMax: 1430})

	glyphData, ok := glyph.Data.(SimpleGlyph)
	assert(t, ok)
	assert(t, len(glyphData.Instructions) == 0)
	assert(t, reflect.DeepEqual(glyphData.EndPtsOfContours, []uint16{6, 10}))

	exp := [...]struct {
		x, y    int16
		overlap bool
	}{
		{x: 501, y: 1430, overlap: true},
		{x: 683, y: 1430},
		{x: 1172, y: 0},
		{x: 983, y: 0},
		{x: 591, y: 1193},
		{x: 199, y: 0},
		{x: 12, y: 0},
		{x: 249, y: 514},
		{x: 935, y: 514},
		{x: 935, y: 352},
		{x: 249, y: 352},
	}
	assert(t, len(glyphData.Points) == len(exp))

	const overlapSimple = 0x40

	for i, v := range glyphData.Points {
		e := exp[i]
		assert(t, v.X == e.x && v.Y == e.y)
		assert(t, (v.Flag&overlapSimple != 0) == e.overlap)
	}
}

func TestGlyphCoordinates2(t *testing.T) {
	filepath := "common/SourceSans-VF.ttf"
	fp := readFontFile(t, filepath)
	ng := numGlyphs(t, fp)
	head, _, err := ParseHead(readTable(t, fp, "head"))
	assertNoErr(t, err)
	loca, err := ParseLoca(readTable(t, fp, "loca"), ng, head.IndexToLocFormat == 1)
	assertNoErr(t, err)
	glyf, err := ParseGlyf(readTable(t, fp, "glyf"), loca)
	assertNoErr(t, err)

	exp := Glyf{
		Glyph{
			XMin: 96, XMax: 528, YMin: 0, YMax: 660,
			Data: SimpleGlyph{
				EndPtsOfContours: []uint16{3, 9, 15, 18, 21},
				Points: []GlyphContourPoint{
					{X: 96, Y: 0},
					{X: 96, Y: 660},
					{X: 528, Y: 660},
					{X: 528, Y: 0},
					{X: 144, Y: 32},
					{X: 476, Y: 32},
					{X: 376, Y: 208},
					{X: 314, Y: 314},
					{X: 310, Y: 314},
					{X: 246, Y: 208},
					{X: 310, Y: 366},
					{X: 314, Y: 366},
					{X: 368, Y: 458},
					{X: 462, Y: 626},
					{X: 160, Y: 626},
					{X: 254, Y: 458},
					{X: 134, Y: 74},
					{X: 288, Y: 340},
					{X: 134, Y: 610},
					{X: 488, Y: 74},
					{X: 488, Y: 610},
					{X: 336, Y: 340},
				},
			},
		},
		{
			XMin: 10, XMax: 510, YMin: 0, YMax: 660,
			Data: SimpleGlyph{
				EndPtsOfContours: []uint16{13, 17},
				Points: []GlyphContourPoint{
					{X: 10, Y: 0},
					{X: 246, Y: 660},
					{X: 274, Y: 660},
					{X: 510, Y: 0},
					{X: 476, Y: 0},
					{X: 338, Y: 396},
					{X: 317, Y: 456},
					{X: 280, Y: 565},
					{X: 262, Y: 626},
					{X: 258, Y: 626},
					{X: 240, Y: 565},
					{X: 203, Y: 456},
					{X: 182, Y: 396},
					{X: 42, Y: 0},
					{X: 112, Y: 236},
					{X: 112, Y: 264},
					{X: 405, Y: 264},
					{X: 405, Y: 236},
				},
			},
		},
		{
			XMin: 10, XMax: 510, YMin: 0, YMax: 846,
			Data: CompositeGlyph{
				Glyphs: []CompositeGlyphPart{
					{GlyphIndex: 1, arg1: 0, arg2: 0, Scale: [4]float32{1, 0, 0, 1}},
					{GlyphIndex: 3, arg1: 260, arg2: 0, Scale: [4]float32{1, 0, 0, 1}},
				},
			},
		},
		{
			XMin: -36, XMax: 104, YMin: 710, YMax: 846,
			Data: SimpleGlyph{
				EndPtsOfContours: []uint16{3},
				Points: []GlyphContourPoint{
					{X: -22, Y: 710},
					{X: -36, Y: 726},
					{X: 82, Y: 846},
					{X: 104, Y: 822},
				},
			},
		},
	}

	assert(t, len(glyf) == len(exp))
	for i, e := range exp {
		g := glyf[i]
		assertGlyphSizeEqual(t, e, g)
		switch data := g.Data.(type) {
		case SimpleGlyph:
			eData := e.Data.(SimpleGlyph)
			assert(t, reflect.DeepEqual(eData.EndPtsOfContours, data.EndPtsOfContours))
			assert(t, len(eData.Points) == len(data.Points))
			for i, got := range data.Points {
				assertPointEqual(t, eData.Points[i], got)
			}
		case CompositeGlyph:
			eData := e.Data.(CompositeGlyph)
			assert(t, len(eData.Glyphs) == len(data.Glyphs))
			for i, got := range data.Glyphs {
				assertCompositeEqual(t, eData.Glyphs[i], got)
			}
		}
	}
}

func TestParseSbix(t *testing.T) {
	for _, file := range td.WithSbix {
		fp := readFontFile(t, file.Path)

		maxp, _, err := ParseMaxp(readTable(t, fp, "maxp"))
		assertNoErr(t, err)

		sbix, _, err := ParseSbix(readTable(t, fp, "sbix"), int(maxp.numGlyphs))
		assertNoErr(t, err)

		assertNoErr(t, err)
		assert(t, len(sbix.Strikes) == file.StrikesNumber)
	}
}

func TestParseCBLC(t *testing.T) {
	for _, file := range td.WithCBLC {
		fp := readFontFile(t, file.Path)

		cblc, _, err := ParseCBLC(readTable(t, fp, "CBLC"))
		assertNoErr(t, err)

		assertNoErr(t, err)
		assert(t, len(cblc.BitmapSizes) == file.StrikesNumber)
	}

	// The following sample has subtable format 3, and is copied from
	// https://github.com/fonttools/fonttools/blob/main/Tests/ttLib/tables/C_B_L_C_test.py
	cblcSample, err := base64.StdEncoding.DecodeString("AAMAAAAAAAEAAAA4AAAALAAAAAIAAAAAZeWIAAAAAAAAAAAAZeWIAAAAAAAAAAAAAAEAA" +
		"21tIAEAAQACAAAAEAADAAMAAAAgAAMAEQAAAAQAAAOmEQ0AAAADABEAABERAAAIUg==")
	assertNoErr(t, err)
	cblc, _, err := ParseCBLC(cblcSample)
	assertNoErr(t, err)
	assert(t, len(cblc.BitmapSizes) == 1)
}

func TestParseEBLC(t *testing.T) {
	for _, file := range td.WithEBLC {
		fp := readFontFile(t, file.Path)

		cblc, _, err := ParseCBLC(readTable(t, fp, "EBLC"))
		assertNoErr(t, err)

		assertNoErr(t, err)
		assert(t, len(cblc.BitmapSizes) == file.StrikesNumber)
	}
}

func TestParseVORG(t *testing.T) {
	filename := "collections/NotoSansCJK-Bold.ttc"

	file, err := td.Files.ReadFile(filename)
	assertNoErr(t, err)

	fonts, err := opentype.NewLoaders(bytes.NewReader(file))
	assertNoErr(t, err)

	for _, fp := range fonts {
		vorg, _, err := ParseVORG(readTable(t, fp, "VORG"))
		assertNoErr(t, err)
		assert(t, len(vorg.VertOriginYMetrics) == 233)
		assert(t, vorg.DefaultVertOriginY == 880)
	}
}
