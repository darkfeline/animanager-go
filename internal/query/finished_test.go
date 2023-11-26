// Copyright (C) 2020  Allen Li
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
	"reflect"
	"testing"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/date"
)

func TestGetAnimeFinished(t *testing.T) {
	db := database.OpenMem(t)
	e := []anidb.Episode{
		{
			EID:    113,
			EpNo:   "1",
			Length: 25,
			Titles: []anidb.EpTitle{{Title: "使徒, 襲来", Lang: "ja"}},
		},
	}
	const aid = 22
	a := &anidb.Anime{
		AID:          aid,
		Type:         "TV Series",
		EpisodeCount: 1, // Modified for testing
		Titles:       []anidb.Title{{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"}},
		Episodes:     e,
		EndDate:      "1996-03-27",
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatal(err)
	}
	if err := UpdateEpisodeDone(db, 113, true); err != nil {
		t.Fatal(err)
	}
	got, err := GetAnimeFinished(db)
	if err != nil {
		t.Fatal(err)
	}
	want := []AnimeBool{
		{
			Anime: &Anime{
				AID:          22,
				Title:        "Shinseiki Evangelion",
				Type:         "TV Series",
				EpisodeCount: 1,
				NullEndDate:  date.New(1996, 3, 27).NullInt64(),
			},
			Value: true,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAnimeFinished() = %#v; want %#v", got, want)
	}
}

func TestIsAnimeFinished(t *testing.T) {
	t.Parallel()
	t.Run("finished", func(t *testing.T) {
		t.Parallel()
		a := Anime{
			AID:          22,
			Title:        "Shinseiki Evangelion",
			Type:         "TV Series",
			EpisodeCount: 1,
			NullEndDate:  date.New(1996, 3, 27).NullInt64(),
		}
		eps := []Episode{
			{
				Type:        EpRegular,
				Number:      1,
				UserWatched: true,
			},
		}
		got := isAnimeFinished(&a, eps)
		want := true
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetAnimeFinished() = %#v; want %#v", got, want)
		}
	})
	t.Run("not finished", func(t *testing.T) {
		t.Parallel()
		a := Anime{
			AID:          22,
			Title:        "Shinseiki Evangelion",
			Type:         "TV Series",
			EpisodeCount: 1,
			NullEndDate:  date.New(1996, 3, 27).NullInt64(),
		}
		eps := []Episode{
			{
				Type:        EpRegular,
				Number:      1,
				UserWatched: false,
			},
		}
		got := isAnimeFinished(&a, eps)
		want := false
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetAnimeFinished() = %#v; want %#v", got, want)
		}
	})
}
