package tables

import (
	"bytes"
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/opentype"
)

func assert(t *testing.T, b bool) {
	t.Helper()
	assertC(t, b, "assertion error")
}

func assertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertC(t *testing.T, b bool, context string) {
	t.Helper()
	if !b {
		t.Fatal(context)
	}
}

// filenames return the "absolute" file names of the given directory
// excluding directories, and not recursing
func filenames(t *testing.T, dir string) []string {
	t.Helper()

	files, err := td.Files.ReadDir(dir)
	assertNoErr(t, err)

	var out []string
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		filename := filepath.Join(dir, entry.Name())
		out = append(out, filename)
	}
	return out
}

func readFontFile(t *testing.T, filepath string) *opentype.Loader {
	t.Helper()

	file, err := td.Files.ReadFile(filepath)
	assertC(t, err == nil, "reading "+filepath)

	fp, err := opentype.NewLoader(bytes.NewReader(file))
	assertC(t, err == nil, "reading "+filepath)

	return fp
}

func readTable(t *testing.T, fl *opentype.Loader, tag string) []byte {
	t.Helper()

	table, err := fl.RawTable(opentype.MustNewTag(tag))
	assertNoErr(t, err)

	return table
}

func numGlyphs(t *testing.T, fp *opentype.Loader) int {
	t.Helper()

	table := readTable(t, fp, "maxp")
	maxp, _, err := ParseMaxp(table)
	assertNoErr(t, err)

	return int(maxp.numGlyphs)
}

func TestParseBasicTables(t *testing.T) {
	for _, filename := range append(filenames(t, "morx"), filenames(t, "common")...) {
		fp := readFontFile(t, filename)
		_, _, err := ParseOs2(readTable(t, fp, "OS/2"))
		assertNoErr(t, err)

		_, _, err = ParseHead(readTable(t, fp, "head"))
		assertNoErr(t, err)

		_, _, err = ParseMaxp(readTable(t, fp, "maxp"))
		assertNoErr(t, err)

		_, _, err = ParseName(readTable(t, fp, "name"))
		assertNoErr(t, err)
	}
}

func TestParseAATMorx(t *testing.T) {
	files := filenames(t, "morx")
	for _, filename := range files {
		fp := readFontFile(t, filename)
		table, _, err := ParseMorx(readTable(t, fp, "morx"))
		assertNoErr(t, err)
		assert(t, int(table.nChains) == len(table.chains))

		ng := numGlyphs(t, fp)
		for _, chain := range table.chains {
			subtables, err := chain.process(ng)
			assertNoErr(t, err)
			assert(t, len(subtables) == int(chain.nSubtable))
		}
	}
}

func TestParseCmap(t *testing.T) {
	// general parsing
	for _, filename := range filenames(t, "common") {
		fp := readFontFile(t, filename)
		cmap, _, err := ParseCmap(readTable(t, fp, "cmap"))
		assertNoErr(t, err)
		assert(t, len(cmap.records) == len(cmap.subtables))
	}

	// specialized tests for each format
	for _, filename := range filenames(t, "cmap") {
		fp := readFontFile(t, filename)
		cmap, _, err := ParseCmap(readTable(t, fp, "cmap"))
		assertNoErr(t, err)
		assert(t, len(cmap.records) == len(cmap.subtables))
	}

	// test format 2 through a single table
	file, err := td.Files.ReadFile("cmap/table/CMAP2.bin")
	assertNoErr(t, err)
	cmap, _, err := ParseCmapSubtable2(file)
	assertNoErr(t, err)
	assert(t, cmap.format == 2)
}

func TestParseOTLayout(t *testing.T) {
	for _, filename := range td.WithOTLayout {
		fp := readFontFile(t, filename)
		gsub, _, err := ParseLayout(readTable(t, fp, "GSUB"))
		assertNoErr(t, err)
		assert(t, len(gsub.lookupList.lookups) == len(gsub.lookupList.lookups))
		assert(t, len(gsub.lookupList.lookups) > 0)

		gpos, _, err := ParseLayout(readTable(t, fp, "GPOS"))
		assertNoErr(t, err)
		assert(t, len(gpos.lookupList.lookups) == len(gpos.lookupList.lookups))
		assert(t, len(gpos.lookupList.lookups) > 0)
	}

	for _, filename := range filenames(t, "toys/gsub") {
		fp := readFontFile(t, filename)
		gsub, _, err := ParseLayout(readTable(t, fp, "GSUB"))
		assertNoErr(t, err)
		assert(t, len(gsub.lookupList.lookups) == len(gsub.lookupList.lookups))
		assert(t, len(gsub.lookupList.lookups) > 0)

		for _, lookup := range gsub.lookupList.lookups {
			assert(t, lookup.lookupType > 0)
			_, err = lookup.AsGSUBLookups()
			assertNoErr(t, err)
		}
	}
}
