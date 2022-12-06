package tables

import (
	"bytes"
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/loader"
)

func assert(t *testing.T, b bool) {
	t.Helper()
	assertC(t, b, "assertion error")
}

func assertNoErr(t testing.TB, err error) {
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
func filenames(t testing.TB, dir string) []string {
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

// wrap td.Files.ReadFile
func readFontFile(t testing.TB, filepath string) *loader.Loader {
	t.Helper()

	file, err := td.Files.ReadFile(filepath)
	assertNoErr(t, err)

	fp, err := loader.NewLoader(bytes.NewReader(file))
	assertNoErr(t, err)

	return fp
}

func readTable(t testing.TB, fl *loader.Loader, tag string) []byte {
	t.Helper()

	table, err := fl.RawTable(loader.MustNewTag(tag))
	assertNoErr(t, err)

	return table
}

func numGlyphs(t *testing.T, fp *loader.Loader) int {
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

func TestParseCmap(t *testing.T) {
	// general parsing
	for _, filename := range filenames(t, "common") {
		fp := readFontFile(t, filename)
		_, _, err := ParseCmap(readTable(t, fp, "cmap"))
		assertNoErr(t, err)
	}

	// specialized tests for each format
	for _, filename := range filenames(t, "cmap") {
		fp := readFontFile(t, filename)
		_, _, err := ParseCmap(readTable(t, fp, "cmap"))
		assertNoErr(t, err)
	}

	// tests through a single table
	for _, filepath := range filenames(t, "cmap/table") {
		table, err := td.Files.ReadFile(filepath)
		assertNoErr(t, err)

		cmap, _, err := ParseCmap(table)
		assertNoErr(t, err)
		assert(t, len(cmap.Records) > 0)
	}
}
