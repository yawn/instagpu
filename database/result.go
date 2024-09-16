package database

import (
	"fmt"

	"github.com/yawn/spottty/detect"
)

type Result struct {
	Index    int
	Prices   *detect.Prices
	Relative float64
	Score    float64
}

func (s *Result) String() string {
	return fmt.Sprintf("%f", s.Score)
}
