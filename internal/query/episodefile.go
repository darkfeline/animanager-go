package query

import (
	"database/sql"
	"fmt"
)

type EpisodeFile struct {
	EpisodeID int
	Path      string
}

// InsertEpisodeFile inserts a file for an episode into the database.
func InsertEpisodeFile(db *sql.DB, id int, path string) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`INSERT INTO episode_file (episode_id, path) VALUES (?, ?)`, id, path)
	if err != nil {
		return fmt.Errorf("insert episode %d file for %s: %s", id, path, err)
	}
	return t.Commit()
}

// GetEpisodeFiles returns the EpisodeFiles for the episode.
func GetEpisodeFiles(db *sql.DB, episodeID int) (es []EpisodeFile, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("get episode %d files: %s", episodeID, err)
		}
	}()
	r, err := db.Query(`
SELECT episode_id, path
FROM episode_file WHERE episode_id=?`, episodeID)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for r.Next() {
		var e EpisodeFile
		if err := r.Scan(&e.EpisodeID, &e.Path); err != nil {
			return nil, err
		}
		es = append(es, e)
	}
	if r.Err() != nil {
		return nil, r.Err()
	}
	return es, nil
}

// DeleteEpisodeFiles deletes all episode files.
func DeleteEpisodeFiles(db *sql.DB) error {
	t, err := db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`DELETE FROM episode_file`)
	if err != nil {
		return err
	}
	return t.Commit()
}
