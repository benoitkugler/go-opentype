package tables

import (
	"reflect"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

func TestParseTrak(t *testing.T) {
	fp := readFontFile(t, "toys/Trak.ttf")
	trak, _, err := ParseTrak(readTable(t, fp, "trak"))
	assertNoErr(t, err)
	assert(t, len(trak.Horiz.SizeTable) == 4)
	assert(t, len(trak.Vert.SizeTable) == 0)

	assert(t, reflect.DeepEqual(trak.Horiz.SizeTable, []float32{1, 2, 12, 96}))
	assert(t, reflect.DeepEqual(trak.Horiz.TrackTable[0].PerSizeTracking, []int16{200, 200, 0, -100}))
}

func TestParseFeat(t *testing.T) {
	fp := readFontFile(t, "toys/Feat.ttf")
	feat, _, err := ParseFeat(readTable(t, fp, "feat"))
	assertNoErr(t, err)

	expectedSettings := [...][]FeatureSettingName{
		{{2, 260}, {4, 259}, {10, 304}},
		{{0, 309}, {1, 263}, {3, 264}},
		{{0, 266}, {1, 267}},
		{{0, 271}, {2, 272}, {8, 273}},
		{{0, 309}, {1, 275}, {2, 277}, {3, 278}},
		{{0, 309}, {2, 280}},
		{{0, 283}},
		{{8, 308}},
		{{0, 309}, {3, 289}},
		{{0, 294}, {1, 295}, {2, 296}, {3, 297}},
		{{0, 309}, {1, 301}},
	}
	assert(t, len(feat.Names) == len(expectedSettings))
	for i, name := range feat.Names {
		exp := expectedSettings[i]
		got := name.SettingTable
		assert(t, reflect.DeepEqual(exp, got))
	}
}

func TestParseAnkr(t *testing.T) {
	table, err := td.Files.ReadFile("toys/tables/ankr.bin")
	assertNoErr(t, err)

	ankr, _, err := ParseAnkr(table, 1409)
	assertNoErr(t, err)

	_, isFormat4 := ankr.LookupTable.(AATLoopkup4)
	assert(t, isFormat4)
}
