package testutils

import (
	"embed"
	"path/filepath"
	"testing"

	"github.com/benoitkugler/go-opentype-testdata/opentype"
)

func Assert(t testing.TB, b bool) {
	t.Helper()
	AssertC(t, b, "assertion error")
}

func AssertNoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func AssertC(t testing.TB, b bool, context string) {
	t.Helper()
	if !b {
		t.Fatal(context)
	}
}

// Filenames return the "absolute" file names of the given directory
// excluding directories, and not recursing.
// It uses the opentype embed file system.
func Filenames(t testing.TB, dir string) []string {
	return FilenamesFS(t, &opentype.Files, dir)
}

func FilenamesFS(t testing.TB, fs *embed.FS, dir string) []string {
	t.Helper()

	files, err := fs.ReadDir(dir)
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
