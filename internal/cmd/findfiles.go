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

package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/google/subcommands"
	"golang.org/x/xerrors"

	"go.felesatra.moe/animanager/internal/database"
	"go.felesatra.moe/animanager/internal/query"
)

type FindFiles struct {
}

func (*FindFiles) Name() string     { return "findfiles" }
func (*FindFiles) Synopsis() string { return "Find episode files." }
func (*FindFiles) Usage() string {
	return `Usage: findfiles
Find episode files.
`
}

func (*FindFiles) SetFlags(f *flag.FlagSet) {
}

func (ff *FindFiles) Execute(ctx context.Context, f *flag.FlagSet, x ...interface{}) subcommands.ExitStatus {
	if f.NArg() != 0 {
		fmt.Fprint(os.Stderr, ff.Usage())
		return subcommands.ExitUsageError
	}

	c := getConfig(x)
	db, err := database.Open(ctx, c.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		return subcommands.ExitFailure
	}
	defer db.Close()
	Logger.Printf("Finding video files")
	files, err := findVideoFilesMany(c.WatchDirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	Logger.Printf("Finished finding video files")
	if err := refreshFiles(db, files); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func findVideoFilesMany(dirs []string) ([]string, error) {
	var result []string
	for _, d := range dirs {
		r, err := findVideoFiles(d)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}
	return result, nil
}

// findVideoFiles returns a slice of paths of all video files found
// recursively under the given path.  The returned paths are absolute.
func findVideoFiles(path string) (result []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("find video files in %s: %s", path, err)
		}
	}()
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if isVideoFile(path, info) {
			result = append(result, path)
		}
		return nil
	})
	return result, nil
}

var videoExts = []string{".mkv", ".mp4", ".avi"}

func isVideoFile(path string, fi os.FileInfo) bool {
	if fi.IsDir() {
		return false
	}
	ext := filepath.Ext(path)
	for _, s := range videoExts {
		if ext == s {
			return true
		}
	}
	return false
}

// refreshFiles updates episode files using the given video file
// paths.
func refreshFiles(db *sql.DB, files []string) error {
	ws, err := query.GetAllWatching(db)
	if err != nil {
		return xerrors.Errorf("refresh files: %w", err)
	}
	if err := query.DeleteEpisodeFiles(db); err != nil {
		return xerrors.Errorf("refresh files: %w", err)
	}
	var efs []query.EpisodeFile
	Logger.Print("Matching files")
	for _, w := range ws {
		Logger.Printf("Matching files for %d", w.AID)
		eps, err := query.GetEpisodes(db, w.AID)
		if err != nil {
			return xerrors.Errorf("refresh files: %w", err)
		}
		efs2, err := filterFiles(w, eps, files)
		if err != nil {
			return xerrors.Errorf("refresh files: %w", err)
		}
		Logger.Printf("Found files for %d: %#v", w.AID, efs2)
		efs = append(efs, efs2...)
	}
	Logger.Print("Inserting files")
	err = query.InsertEpisodeFiles(db, efs)
	return xerrors.Errorf("refresh files: %w", err)
}

// filterFiles returns files that match the watching entry.
func filterFiles(w query.Watching, eps []query.Episode, files []string) ([]query.EpisodeFile, error) {
	var result []query.EpisodeFile
	r, err := regexp.Compile("(?i)" + w.Regexp)
	if err != nil {
		return nil, xerrors.Errorf("filter files for %d: %w", w.AID, err)
	}
	regEps := reindexEpisodeByNumber(eps)
	for _, f := range files {
		ms := r.FindStringSubmatch(filepath.Base(f))
		if ms == nil {
			continue
		}
		if len(ms) < 2 {
			return nil, fmt.Errorf("filter files for %d: regexp %#v has no submatch",
				w.AID, w.Regexp)

		}
		n, err := strconv.Atoi(ms[1])
		if err != nil {
			return nil, fmt.Errorf("filter files for %d: regexp %#v submatch not a number",
				w.AID, w.Regexp)
		}
		n += w.Offset
		if n >= len(regEps) {
			continue
		}
		result = append(result, query.EpisodeFile{
			EpisodeID: regEps[n].ID,
			Path:      f,
		})
	}
	return result, nil
}

// reindexEpisodeByNumber returns a slice where each index maps to the regular
// episode with the same number.  The zero index will be empty since
// episodes cannot be number zero.
func reindexEpisodeByNumber(eps []query.Episode) []query.Episode {
	m := make([]query.Episode, maxEpisodeNumber(eps)+1)
	for _, e := range eps {
		if e.Type == query.EpRegular {
			m[e.Number] = e
		}
	}
	return m
}

func maxEpisodeNumber(eps []query.Episode) int {
	maxEp := 0
	for _, e := range eps {
		if e.Type == query.EpRegular && e.Number > maxEp {
			maxEp = e.Number
		}
	}
	return maxEp
}
