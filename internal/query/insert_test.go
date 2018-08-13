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

package query

import (
	"context"
	"reflect"
	"testing"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/database"
)

func TestInsertAndGetAnime(t *testing.T) {
	db, err := database.OpenMem(context.Background())
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	e := []anidb.Episode{
		{
			EpNo:   "1",
			Length: 25,
			Titles: []anidb.EpTitle{
				{Title: "使徒, 襲来", Lang: "ja"},
				{Title: "Angel Attack!", Lang: "en"},
				{Title: "Shito, Shuurai", Lang: "x-jat"},
			},
		},
		{
			EpNo:   "S1",
			Length: 75,
			Titles: []anidb.EpTitle{
				{Title: "Revival of Evangelion Extras Disc", Lang: "en"},
			},
		},
	}
	a := &anidb.Anime{
		AID:          22,
		Type:         "TV Series",
		EpisodeCount: 26,
		StartDate:    "1995-10-04",
		EndDate:      "1996-03-27",
		Titles: []anidb.Title{
			{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"},
		},
		Episodes: e,
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	t.Run("get anime", func(t *testing.T) {
		got, err := GetAnime(db, 22)
		if err != nil {
			t.Fatalf("Error getting anime: %s", err)
		}
		want := &Anime{
			AID:          22,
			Title:        "Shinseiki Evangelion",
			Type:         "TV Series",
			EpisodeCount: 26,
			NStartDate:   int64(812764800),
			NEndDate:     int64(827884800),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetAnime(db, 22) = %#v; want %#v", got, want)
		}
	})
	t.Run("get episodes", func(t *testing.T) {
		got, err := GetEpisodes(db, 22)
		if err != nil {
			t.Fatalf("Error getting episodes: %s", err)
		}
		want := []Episode{
			{
				ID:     1,
				AID:    22,
				Type:   EpRegular,
				Number: 1,
				Title:  "使徒, 襲来",
				Length: 25,
			},
			{
				ID:     2,
				AID:    22,
				Type:   EpSpecial,
				Number: 1,
				Title:  "Revival of Evangelion Extras Disc",
				Length: 75,
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetEpisodes(db, 22) = %#v; want %#v", got, want)
		}
	})
}

func TestInsertAndGetWatching(t *testing.T) {
	db, err := database.OpenMem(context.Background())
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	aid := 22
	a := &anidb.Anime{
		AID:          aid,
		Type:         "TV Series",
		EpisodeCount: 26,
		StartDate:    "1995-10-04",
		EndDate:      "1996-03-27",
		Titles: []anidb.Title{
			{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"},
		},
		Episodes: []anidb.Episode{},
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	p := "foobar"
	if err := InsertWatching(db, aid, p); err != nil {
		t.Fatalf("Error inserting watching: %s", err)
	}
	got, err := GetWatching(db, aid)
	if err != nil {
		t.Fatalf("Error getting anime: %s", err)
	}
	want := Watching{AID: aid, Regexp: p}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetWatching(db, %d) = %#v; want %#v", aid, got, want)
	}
}

func TestInsertAndGetAllWatching(t *testing.T) {
	db, err := database.OpenMem(context.Background())
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	aid := 22
	a := &anidb.Anime{
		AID:          aid,
		Type:         "TV Series",
		EpisodeCount: 26,
		StartDate:    "1995-10-04",
		EndDate:      "1996-03-27",
		Titles: []anidb.Title{
			{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"},
		},
		Episodes: []anidb.Episode{},
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	p := "foobar"
	if err := InsertWatching(db, aid, p); err != nil {
		t.Fatalf("Error inserting watching: %s", err)
	}
	got, err := GetAllWatching(db)
	if err != nil {
		t.Fatalf("Error getting anime: %s", err)
	}
	want := []Watching{{AID: aid, Regexp: p}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAllWatching(db) = %#v; want %#v", got, want)
	}
}

func TestInsertAndGetAnime_nullFields(t *testing.T) {
	db, err := database.OpenMem(context.Background())
	if err != nil {
		t.Fatalf("Error opening database: %s", err)
	}
	defer db.Close()
	a := &anidb.Anime{
		AID:          22,
		Type:         "TV Series",
		EpisodeCount: 26,
		Titles: []anidb.Title{
			{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"},
		},
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	got, err := GetAnime(db, 22)
	if err != nil {
		t.Fatalf("Error getting anime: %s", err)
	}
	want := &Anime{
		AID:          22,
		Title:        "Shinseiki Evangelion",
		Type:         "TV Series",
		EpisodeCount: 26,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAnime(db, 22) = %#v; want %#v", got, want)
	}
}
