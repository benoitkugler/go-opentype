package tables

// https://learn.microsoft.com/en-us/typography/opentype/spec/hmtx
type Hmtx struct {
	hMetrics []longHorMetric `arrayCount:""`
	// avances are padded with the last value
	// and side bearings are given
	leftSideBearings []int16 `arrayCount:""`
}

type longHorMetric struct {
	advanceWidth, lsb int16
}

// https://learn.microsoft.com/en-us/typography/opentype/spec/vmtx
type Vmtx = Hmtx
