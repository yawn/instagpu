package database

import (
	"fmt"
	"strings"

	"github.com/yawn/spottty/detect"
)

type Result struct {
	Index    int            `json:"index"`
	Prices   *detect.Prices `json:"prices"`
	Relative float64        `json:"score_relative_to_best"`
	Score    float64        `json:"score"`
}

func (s *Result) String() string {

	var b strings.Builder

	fmt.Fprintf(&b, "🏅 %2d", s.Index)
	fmt.Fprintf(&b, "\t🔢 %3.2f", s.Score)
	fmt.Fprintf(&b, "\t🔝 %3.2f%%", s.Relative * 100)
	fmt.Fprintf(&b, "\t%s", s.Prices.String())

	return b.String()

}
