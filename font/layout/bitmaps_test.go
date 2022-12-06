package layout

import (
	"bytes"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
	"github.com/benoitkugler/go-opentype/opentype"
	"github.com/benoitkugler/go-opentype/tables"
)

func assert(t *testing.T, b bool) {
	t.Helper()
	if !b {
		t.Fatal("assertion error")
	}
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// wrap td.Files.ReadFile
func readFontFile(t testing.TB, filepath string) *opentype.Loader {
	t.Helper()

	file, err := td.Files.ReadFile(filepath)
	assertNoErr(t, err)

	fp, err := opentype.NewLoader(bytes.NewReader(file))
	assertNoErr(t, err)

	return fp
}

func readTable(t testing.TB, fl *opentype.Loader, tag string) []byte {
	t.Helper()

	table, err := fl.RawTable(opentype.MustNewTag(tag))
	assertNoErr(t, err)

	return table
}

func TestBloc(t *testing.T) {
	blocT, err := td.Files.ReadFile("toys/tables/bloc.bin")
	assertNoErr(t, err)
	bloc, _, err := tables.ParseCBLC(blocT)
	assertNoErr(t, err)

	bdatT, err := td.Files.ReadFile("toys/tables/bdat.bin")
	assertNoErr(t, err)

	bt, err := newBitmap(bloc, bdatT)
	assertNoErr(t, err)
	assert(t, len(bt) == 1)
	assert(t, len(bt[0].subTables) == 4)
}

func TestCBLC(t *testing.T) {
	for _, file := range td.WithCBLC {
		fp := readFontFile(t, file.Path)

		cblc, _, err := tables.ParseCBLC(readTable(t, fp, "CBLC"))
		assertNoErr(t, err)
		cbdt := readTable(t, fp, "CBDT")

		_, err = newBitmap(cblc, cbdt)
		assertNoErr(t, err)
	}
}

func TestEBLC(t *testing.T) {
	for _, file := range td.WithEBLC {
		fp := readFontFile(t, file.Path)

		eblc, _, err := tables.ParseCBLC(readTable(t, fp, "EBLC"))
		assertNoErr(t, err)
		ebdt := readTable(t, fp, "EBDT")

		_, err = newBitmap(eblc, ebdt)
		assertNoErr(t, err)
	}
}
