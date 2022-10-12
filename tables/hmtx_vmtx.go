package tables

//go:generate ../../binarygen/cmd/generator hmtx_vmtx.go

// https://learn.microsoft.com/en-us/typography/opentype/spec/hmtx
type Hmtx struct {
	hMetrics []longHorMetric `len:""`
	// avances are padded with the last value
	// and side bearings are given
	leftSideBearings []int16 `len:""`
}

type longHorMetric struct {
	advanceWidth, lsb int16
}

// https://learn.microsoft.com/en-us/typography/opentype/spec/vmtx
type Vmtx = Hmtx
