package testutils

import (
	"path/filepath"
	"testing"

	td "github.com/benoitkugler/go-opentype-testdata/data"
)

func Assert(t *testing.T, b bool) {
	t.Helper()
	AssertC(t, b, "assertion error")
}

func AssertNoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func AssertC(t *testing.T, b bool, context string) {
	t.Helper()
	if !b {
		t.Fatal(context)
	}
}

// Filenames return the "absolute" file names of the given directory
// excluding directories, and not recursing
func Filenames(t testing.TB, dir string) []string {
	t.Helper()

	files, err := td.Files.ReadDir(dir)
	AssertNoErr(t, err)

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
