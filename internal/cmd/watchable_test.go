package cmd

import (
	"context"
	"reflect"
	"testing"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

func TestGetAnimeFiles(t *testing.T) {
	db, err := database.OpenMem(context.Background())
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	e := []anidb.Episode{
		{
			EpNo:   "1",
			Length: 25,
			Titles: []anidb.EpTitle{{Title: "使徒, 襲来", Lang: "ja"}},
		},
	}
	aid := 22
	a := &anidb.Anime{
		AID:          aid,
		Type:         "TV Series",
		EpisodeCount: 26,
		Titles:       []anidb.Title{{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"}},
		Episodes:     e,
	}
	if err := query.InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	if err := query.InsertEpisodeFile(db, 1, "/foobar"); err != nil {
		t.Fatalf("Error inserting episode file: %s", err)
	}
	got, err := getAnimeFiles(db, aid)
	if err != nil {
		t.Fatalf("getAnimeFiles returned error: %s", err)
	}
	want := animeFiles{
		Episodes: []query.Episode{
			{
				ID:     1,
				AID:    aid,
				Type:   query.EpRegular,
				Number: 1,
				Title:  "使徒, 襲来",
				Length: 25,
			},
		},
		Files: [][]query.EpisodeFile{
			{{EpisodeID: 1, Path: "/foobar"}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("getAnimeFiles() = %#v; want %#v", got, want)
	}
}
