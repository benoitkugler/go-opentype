package tables

import (
	"fmt"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

func TestParseKern(t *testing.T) {
	filepath := "common/FreeSerif.ttf"
	fp := readFontFile(t, filepath)
	_, _, err := ParseKern(readTable(t, fp, "kern"))
	assertNoErr(t, err)

	for _, filepath := range []string{
		"toys/tables/kern0Exp.bin",
		"toys/tables/kern1.bin",
		"toys/tables/kern02.bin",
		"toys/tables/kern3.bin",
	} {
		table, err := td.Files.ReadFile(filepath)
		assertNoErr(t, err)

		kern, _, err := ParseKern(table)
		assertNoErr(t, err)
		assert(t, len(kern.Tables) > 0)

		for _, subtable := range kern.Tables {
			assert(t, subtable.Data() != nil)
		}
	}
}

func TestParseKern0(t *testing.T) {
	table, err := td.Files.ReadFile("toys/tables/kern0Exp.bin")
	assertNoErr(t, err)

	kern, _, err := ParseKern(table)
	assertNoErr(t, err)
	assert(t, len(kern.Tables) == 1)
	data, ok := kern.Tables[0].Data().(KernData0)
	assert(t, ok)

	expecteds := []struct { // value extracted from harfbuzz run
		left, right GlyphID
		kerning     int16
	}{
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 104, 0},
		{104, 504, -42},
		{504, 1108, 0},
		{1108, 65535, 0},
		{65535, 552, 0},
		{552, 573, 0},
		{573, 1059, 0},
		{1059, 65535, 0},
	}

	for _, exp := range expecteds {
		got := data.KernPair(exp.left, exp.right)
		assert(t, got == exp.kerning)
	}
}

func TestKern2(t *testing.T) {
	filepath := "toys/Kern2.ttf"
	fp := readFontFile(t, filepath)
	kern, _, err := ParseKern(readTable(t, fp, "kern"))
	assertNoErr(t, err)

	assert(t, len(kern.Tables) == 3)

	k0, k1, k2 := kern.Tables[0].Data(), kern.Tables[1].Data(), kern.Tables[2].Data()
	_, is0 := k0.(KernData0)
	assert(t, is0)

	type2, ok := k1.(KernData2)
	assert(t, ok)
	expectedSubtable := map[[2]GlyphID]int16{
		{67, 68}: 0,
		{68, 69}: 0,
		{69, 70}: -30,
		{70, 71}: 0,
		{71, 72}: 0,
		{72, 73}: -20,
		{73, 74}: 0,
		{74, 75}: 0,
		{75, 76}: 0,
		{76, 77}: 0,
		{77, 78}: 0,
		{78, 79}: 0,
		{79, 80}: 0,
		{80, 81}: 0,
		{81, 82}: 0,
		{36, 57}: 0,
	}
	for k, exp := range expectedSubtable {
		got := type2.KernPair(k[0], k[1])
		assertC(t, exp == got, fmt.Sprintf("invalid kern subtable : for (%d, %d) expected %d, got %d", k[0], k[1], exp, got))
	}

	type2, ok = k2.(KernData2)
	assert(t, ok)
	expectedSubtable = map[[2]GlyphID]int16{
		{36, 57}: -80,
	}
	for k, exp := range expectedSubtable {
		got := type2.KernPair(k[0], k[1])
		assert(t, exp == got)
		assertC(t, exp == got, fmt.Sprintf("invalid kern subtable : for (%d, %d) expected %d, got %d", k[0], k[1], exp, got))
	}
}

func TestKern3(t *testing.T) {
	table, err := td.Files.ReadFile("toys/tables/kern3.bin")
	assertNoErr(t, err)

	kern, _, err := ParseKern(table)
	assertNoErr(t, err)
	assert(t, len(kern.Tables) == 5)

	expectedsLengths := [...][3]int{
		{570, 5688, 92},
		{570, 6557, 104},
		{570, 5832, 107},
		{570, 6083, 106},
		{570, 4828, 82},
	}
	for i := range kern.Tables {
		data, ok := kern.Tables[i].Data().(KernData3)
		assert(t, ok)
		exp := expectedsLengths[i]
		assert(t, len(data.leftClass) == exp[0])
		assert(t, len(data.kernIndex) == exp[1])
		assert(t, len(data.kernings) == exp[2])
	}
}
