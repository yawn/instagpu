package detect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPRTGIndex(t *testing.T) {

	assert := assert.New(t)

	var (
		perf     float64 = 100
		instance         = &Instance{GPU: &GPU{FP32: &perf}}
		prices           = &Prices{Avg: 1, Instance: instance}
	)

	assert.EqualValues(100, prices.PTGPIndex())

	prices.Avg = 2

	assert.EqualValues(50, prices.PTGPIndex())

}
