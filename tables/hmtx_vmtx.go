package tables

//go:generate ../../binarygen/cmd/generator hmtx_vmtx.go

// https://learn.microsoft.com/en-us/typography/opentype/spec/hmtx
type htmx struct {
	hMetrics         []longHorMetric `len-size:"-"`
	leftSideBearings []int16         `len-size:"-"`
}

type longHorMetric struct {
	advanceWidth, lsb int16
}
