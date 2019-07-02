// Copyright (C) 2019  Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"context"
	"reflect"
	"testing"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/database"
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
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	efs := []EpisodeFile{{EpisodeID: 1, Path: "/foobar"}}
	if err := InsertEpisodeFiles(db, efs); err != nil {
		t.Fatalf("Error inserting episode file: %s", err)
	}
	got, err := GetAnimeFiles(db, aid)
	if err != nil {
		t.Fatalf("GetAnimeFiles returned error: %s", err)
	}
	want := []EpisodeFiles{
		{
			Episode: Episode{
				ID:     1,
				AID:    aid,
				Type:   EpRegular,
				Number: 1,
				Title:  "使徒, 襲来",
				Length: 25,
			},
			Files: []EpisodeFile{
				{EpisodeID: 1, Path: "/foobar"},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAnimeFiles() = %#v; want %#v", got, want)
	}
}
