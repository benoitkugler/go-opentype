package tables

//go:generate ../../binarygen/cmd/generator maxp.go

// https://learn.microsoft.com/en-us/typography/opentype/spec/maxp
type maxp struct {
	version   maxpVersionKind
	numGlyphs uint16
	data      maxpVersion `kind-field:"version"`
}

type maxpVersionKind uint32

const (
	maxpVersionKind05 maxpVersionKind = 0x00005000
	maxpVersionKind1  maxpVersionKind = 0x00010000
)

type maxpVersion interface {
	isMaxpVersion()
}

func (maxpVersion05) isMaxpVersion() {}
func (maxpVersion1) isMaxpVersion()  {}

type maxpVersion05 struct{}

type maxpVersion1 struct {
	data [13]uint16
}
