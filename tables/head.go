package tables

//go:generate ../../binarygen/cmd/generator head.go

// TableHead contains critical information about the rest of the font.
// https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6head.html
// https://learn.microsoft.com/en-us/typography/opentype/spec/head
type Head struct {
	majorVersion       uint16
	minorVersion       uint16
	fontRevision       uint32
	checksumAdjustment uint32
	magicNumber        uint32
	flags              uint16
	UnitsPerEm         uint16
	created            longdatetime
	modified           longdatetime
	XMin               int16
	YMin               int16
	XMax               int16
	YMax               int16
	macStyle           uint16
	lowestRecPPEM      uint16
	fontDirectionHint  int16
	IndexToLocFormat   int16
	glyphDataFormat    int16
}
