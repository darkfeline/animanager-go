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
	"reflect"
	"testing"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/database"
)

func TestInsertAndGetWatching(t *testing.T) {
	db := database.OpenMem(t)
	const aid = 22
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
	want := Watching{
		Aid:    aid,
		Regexp: "foo",
		Offset: 2,
	}
	if err := InsertWatching(db, want); err != nil {
		t.Fatalf("Error inserting watching: %s", err)
	}
	got, err := GetWatching(db, aid)
	if err != nil {
		t.Fatalf("Error getting anime: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetWatching(db, %d) = %#v; want %#v", aid, got, want)
	}
}

func TestInsertInvalidRegexp(t *testing.T) {
	db := database.OpenMem(t)
	const aid = 22
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
	want := Watching{
		Aid:    aid,
		Regexp: "blah[",
		Offset: 2,
	}
	if err := InsertWatching(db, want); err == nil {
		t.Errorf("Expected error")
	}
}

func TestInsertAndGetAllWatching(t *testing.T) {
	db := database.OpenMem(t)
	const aid = 22
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
	want := []Watching{
		{Aid: aid, Regexp: "foo", Offset: 2},
	}
	if err := InsertWatching(db, want[0]); err != nil {
		t.Fatalf("Error inserting watching: %s", err)
	}
	got, err := GetAllWatching(db)
	if err != nil {
		t.Fatalf("Error getting anime: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAllWatching(db) = %#v; want %#v", got, want)
	}
}
