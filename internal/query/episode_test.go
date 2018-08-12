package query

import (
	"testing"
)

func TestParseEpNo(t *testing.T) {
	t.Parallel()
	cases := []struct {
		EpNo   string
		Type   EpisodeType
		Number int
	}{
		{"S1", EpSpecial, 1},
		{"T2", EpTrailer, 2},
		{"15", EpRegular, 15},
		{"Clarion", EpInvalid, 0},
	}
	for _, c := range cases {
		t.Run(c.EpNo, func(t *testing.T) {
			t.Parallel()
			eptype, n := parseEpNo(c.EpNo)
			if eptype != c.Type || n != c.Number {
				t.Errorf("ParseEpNo(%s) = %s, %d (expected %s, %d)",
					c.EpNo, eptype, n, c.Type, c.Number)
			}
		})
	}
}
