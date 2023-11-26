// Copyright (C) 2023  Allen Li
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

package fileid

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindVideoFiles returns a slice of paths of all video files found
// recursively under the given paths.  The returned paths are absolute.
func FindVideoFiles(dirs []string) ([]string, error) {
	var result []string
	for _, d := range dirs {
		r, err := findVideoFilesSingle(d)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}
	return result, nil
}

// findVideoFilesSingle returns a slice of paths of all video files found
// recursively under the given path.  The returned paths are absolute.
func findVideoFilesSingle(path string) (result []string, err error) {
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
