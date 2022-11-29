package tables

import (
	"reflect"
	"testing"
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
