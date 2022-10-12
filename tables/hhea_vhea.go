package tables

//go:generate ../../binarygen/cmd/generator hhea_vhea.go

// https://learn.microsoft.com/en-us/typography/opentype/spec/hhea
type Hhea struct {
	majorVersion         uint16
	minorVersion         uint16
	ascender             int16
	descender            int16
	lineGap              int16
	AdvanceMax           uint16
	MinFirstSideBearing  int16
	MinSecondSideBearing int16
	MaxExtent            int16
	CaretSlopeRise       int16
	CaretSlopeRun        int16
	CaretOffset          int16
	reserved             [4]uint16
	metricDataformat     int16
	NumOfLongMetrics     uint16
}

// https://learn.microsoft.com/en-us/typography/opentype/spec/vhea
type Vhea = Hhea
