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
	"log/slog"
	"reflect"
	"testing"

	"go.felesatra.moe/anidb"
	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/sqlc"
)

func TestGetAnimeFiles(t *testing.T) {
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
		EpisodeCount: 26,
		Titles:       []anidb.Title{{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"}},
		Episodes:     e,
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	efs := []EpisodeFile{{EID: 113, Path: "/foobar"}}
	ctx := context.Background()
	if err := InsertEpisodeFiles(ctx, sqlc.New(db), nullLogger(), efs); err != nil {
		t.Fatalf("Error inserting episode file: %s", err)
	}
	got, err := GetAnimeFiles(db, aid)
	if err != nil {
		t.Fatalf("GetAnimeFiles returned error: %s", err)
	}
	want := []EpisodeFiles{
		{
			Episode: Episode{
				EID:    113,
				AID:    aid,
				Type:   EpRegular,
				Number: 1,
				Title:  "使徒, 襲来",
				Length: 25,
			},
			Files: []EpisodeFile{
				{EID: 113, Path: "/foobar"},
			},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAnimeFiles() = %#v; want %#v", got, want)
	}
}

func TestDeleteAnimeFiles(t *testing.T) {
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
		EpisodeCount: 26,
		Titles:       []anidb.Title{{Name: "Shinseiki Evangelion", Type: "main", Lang: "x-jat"}},
		Episodes:     e,
	}
	if err := InsertAnime(db, a); err != nil {
		t.Fatalf("Error inserting anime: %s", err)
	}
	efs := []EpisodeFile{{EID: 113, Path: "/foobar"}}
	ctx := context.Background()
	if err := InsertEpisodeFiles(ctx, sqlc.New(db), nullLogger(), efs); err != nil {
		t.Fatalf("Error inserting episode file: %s", err)
	}
	if err := DeleteAnimeFiles(db, aid); err != nil {
		t.Fatalf("Error deleting episode files: %s", err)
	}
	got, err := GetAnimeFiles(db, aid)
	if err != nil {
		t.Fatalf("GetAnimeFiles returned error: %s", err)
	}
	want := []EpisodeFiles{
		{
			Episode: Episode{
				EID:    113,
				AID:    aid,
				Type:   EpRegular,
				Number: 1,
				Title:  "使徒, 襲来",
				Length: 25,
			},
			Files: nil,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAnimeFiles() = %#v; want %#v", got, want)
	}
}

func nullLogger() *slog.Logger {
	return slog.New(nullHandler{})
}

type nullHandler struct{}

var _ slog.Handler = nullHandler{}

func (nullHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (nullHandler) Handle(context.Context, slog.Record) error { return nil }
func (h nullHandler) WithAttrs([]slog.Attr) slog.Handler      { return h }
func (h nullHandler) WithGroup(string) slog.Handler           { return h }
