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
	"sync"

	"github.com/google/subcommands"
	"go.felesatra.moe/go2/errors"

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

// refreshFiles updates episode files using the given video file
// paths.
func refreshFiles(db *sql.DB, files []string) error {
	ws, err := query.GetAllWatching(db)
	if err != nil {
		return errors.Wrap(err, "refresh files")
	}
	if err := query.DeleteEpisodeFiles(db); err != nil {
		return errors.Wrap(err, "refresh files")
	}

	var wg sync.WaitGroup
	epChan := make(chan epFile)
	errChan := make(chan error)
	for _, w := range ws {
		Logger.Printf("Finding registered files for %d", w.AID)
		eps, err := query.GetEpisodes(db, w.AID)
		if err != nil {
			return errors.Wrap(err, "refresh files")
		}
		wg.Add(1)
		go func(w query.Watching) {
			findRegisteredFiles(w, eps, files, epChan, errChan)
			Logger.Printf("Finished finding registered files for %d", w.AID)
			wg.Done()
		}(w)
	}
	go func() {
		wg.Wait()
		close(epChan)
	}()
	go insertEpisodeFiles(db, epChan, errChan)
	for err := range errChan {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	return nil
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

type epFile struct {
	ID   int
	Path string
}

// findRegisteredFiles sends files that match the watching entry as
// epFiles to the output channel.  Errors are sent to the error
// channel.
func findRegisteredFiles(w query.Watching, eps []query.Episode, files []string, c chan<- epFile, e chan<- error) {
	r, err := regexp.Compile("(?i)" + w.Regexp)
	if err != nil {
		e <- err
		return
	}
	regEps := make([]query.Episode, maxEpisodeNumber(eps)+1)
	for _, e := range eps {
		if e.Type == query.EpRegular {
			regEps[e.Number] = e
		}
	}
	for _, f := range files {
		ms := r.FindStringSubmatch(filepath.Base(f))
		if ms == nil {
			continue
		}
		if len(ms) < 2 {
			e <- fmt.Errorf("regexp %#v has no submatch", w.Regexp)
			return
		}
		n, err := strconv.Atoi(ms[1])
		if err != nil {
			e <- fmt.Errorf("regexp %#v submatch not a number", w.Regexp)
			return
		}
		n += w.Offset
		if n >= len(regEps) {
			continue
		}
		c <- epFile{
			ID:   regEps[n].ID,
			Path: f,
		}
	}
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

// insertEpisodeFiles inserts received episodes files from the channel
// and inserts them into the database.  Errors are sent to the error
// channel.  This function returns when the input channel is closed.
// The error channel is closed at the end.
func insertEpisodeFiles(db *sql.DB, c <-chan epFile, e chan<- error) {
	for ep := range c {
		if err := query.InsertEpisodeFile(db, ep.ID, ep.Path); err != nil {
			e <- err
		}
	}
	close(e)
}
