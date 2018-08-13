package cmd

import (
	"reflect"
	"testing"

	"go.felesatra.moe/animanager/internal/query"
)

func TestMaxEpisodeNumber(t *testing.T) {
	eps := []query.Episode{
		{Type: query.EpRegular, Number: 1},
		{Type: query.EpRegular, Number: 3},
		{Type: query.EpRegular, Number: 5},
		{Type: query.EpOther, Number: 13},
	}
	got := maxEpisodeNumber(eps)
	want := 5
	if got != want {
		t.Errorf("maxEpisodeNumber(%#v) = %#v; want %#v", eps, got, want)
	}
}

func TestFindRegisteredFiles(t *testing.T) {
	w := query.Watching{
		Regexp: "lacia([0-9]+)",
	}
	eps := []query.Episode{
		{ID: 1, Type: query.EpRegular, Number: 1},
		{ID: 2, Type: query.EpRegular, Number: 3},
		{ID: 3, Type: query.EpRegular, Number: 5},
		{ID: 4, Type: query.EpOther, Number: 13},
	}
	files := []string{
		"/foo/lacia1",
		"/foo/lacia1v2",
		"/foo/lacia2",
		"/foo/foobar",
		"/foo/lacia5",
		"/foo/lacia13",
	}
	c := make(chan epFile, 10)
	e := make(chan error, 2)
	findRegisteredFiles(w, eps, files, c, e)
	close(c)
	close(e)
	var got []epFile
	for e := range c {
		got = append(got, e)
	}
	var errs []error
	for e := range e {
		errs = append(errs, e)
	}
	want := []epFile{
		{ID: 1, Path: "/foo/lacia1"},
		{ID: 1, Path: "/foo/lacia1v2"},
		{ID: 0, Path: "/foo/lacia2"},
		{ID: 3, Path: "/foo/lacia5"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("findRegisteredFiles() = %#v; want %#v", got, want)
	}
	if len(errs) > 0 {
		t.Errorf("findRegisteredFiles returned errors %#v", errs)
	}
}
