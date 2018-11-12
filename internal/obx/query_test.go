// Copyright (C) 2018  Allen Li
//
// This file is part of Animanager.
//
// Animanager is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Animanager is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Animanager.  If not, see <http://www.gnu.org/licenses/>.

package obx

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
	got, err := GetAnimeFiles(db, aid)
	if err != nil {
		t.Fatalf("GetAnimeFiles returned error: %s", err)
	}
	want := []EpisodeFiles{
		{
			Episode: query.Episode{
				ID:     1,
				AID:    aid,
				Type:   query.EpRegular,
				Number: 1,
				Title:  "使徒, 襲来",
				Length: 25,
			},
			Files: []query.EpisodeFile{
				{EpisodeID: 1, Path: "/foobar"},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAnimeFiles() = %#v; want %#v", got, want)
	}
}
