package cmd

import (
	"regexp"
	"testing"

	"go.felesatra.moe/animanager/internal/query"
)

func TestAnimeDefaultRegexp(t *testing.T) {
	t.Parallel()
	a := &query.Anime{
		Title: "Keit-ai 2",
	}
	got := animeDefaultRegexp(a)
	p, err := regexp.Compile(got)
	if err != nil {
		t.Fatalf("Could not compile regexp %#v: %s", got, err)
	}
	cases := []struct {
		String  string
		Version string
	}{
		{"[Meme-raws] Keit-ai 2nd Season - 13 END [720p].mkv", "13"},
	}
	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			t.Parallel()
			got := p.FindStringSubmatch(c.String)
			if got[1] != c.Version {
				t.Errorf("p.FindStringSubmatch(%#v) = %#v (expected ANY, %#v)",
					c.String, got, c.Version)
			}
		})
	}
}
