package tables

import (
	"bytes"
	"encoding/base64"
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

		loca, err := parseLoca(readTable(t, fp, "loca"), int(maxp.numGlyphs), head.IndexToLocFormat == 1)
		assertNoErr(t, err)

		glyphs, err := parseGlyf(readTable(t, fp, "glyf"), loca)
		assertNoErr(t, err)
		assert(t, len(glyphs) == len(loca)-1)
		assert(t, len(glyphs) == file.GlyphNumber)
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
