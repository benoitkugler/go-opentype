package tables

import (
	"encoding/binary"
	"fmt"
)

// Code generated by binarygen from hmtx_vmtx.go. DO NOT EDIT

func ParseHmtx(src []byte, hmetricsNum int, leftsidebearingsNum int) (Hmtx, int, error) {
	var item Hmtx
	n := 0
	{
		subSlice := src[n:]
		if L := len(subSlice); L < +hmetricsNum*4 {
			return Hmtx{}, 0, fmt.Errorf("EOF: expected length: %d, got %d", +hmetricsNum*4, L)
		}

		item.hMetrics = make([]longHorMetric, hmetricsNum) // allocation guarded by the previous check
		for i := range item.hMetrics {
			item.hMetrics[i].mustParse(subSlice[+i*4:])
		}

		n += hmetricsNum * 4
	}
	{
		subSlice := src[n:]
		if L := len(subSlice); L < +leftsidebearingsNum*2 {
			return Hmtx{}, 0, fmt.Errorf("EOF: expected length: %d, got %d", +leftsidebearingsNum*2, L)
		}

		item.leftSideBearings = make([]int16, leftsidebearingsNum) // allocation guarded by the previous check
		for i := range item.leftSideBearings {
			item.leftSideBearings[i] = int16(binary.BigEndian.Uint16(subSlice[+i*2:]))
		}

		n += leftsidebearingsNum * 2
	}
	return item, n, nil
}
func (item *longHorMetric) mustParse(src []byte) {
	_ = src[3] // early bound checking
	item.advanceWidth = int16(binary.BigEndian.Uint16(src[0:]))
	item.lsb = int16(binary.BigEndian.Uint16(src[2:]))
}
