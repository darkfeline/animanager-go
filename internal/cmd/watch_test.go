package cmd

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s    string
		want string
	}{
		{"lacia", "lacia"},
		{"lacia\n", "lacia"},
	}
	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			r := bufio.NewReader(strings.NewReader(c.s))
			got, err := readLine(r)
			if err != nil {
				t.Fatalf("readLine(%#v) returned error %s", c.s, err)
			}
			if got != c.want {
				t.Errorf("readLine(%#v) = %#v; want %#v", c.s, got, c.want)
			}
		})
	}
}
