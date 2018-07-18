package migrate

import (
	"database/sql"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
)

var Logger = log.New(ioutil.Discard, "", log.LstdFlags)

func Migrate(d *sql.DB) error {
	for {
		v, err := getUserVersion(d)
		if err != nil {
			return errors.Wrap(err, "get user version")
		}
		m, ok := migrations[v]
		if !ok {
			return nil
		}
		Logger.Printf("Migrating from %d", v)
		if err := m(d); err != nil {
			return errors.Wrapf(err, "migrate from %d", v)
		}
	}
}

var migrations = map[int]migrateFunc{
	0: migrate3,
}

type migrateFunc func(*sql.DB) error

func migrate3(d *sql.DB) error {
	t, err := d.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	_, err = t.Exec(`
CREATE TABLE "anime" (
    aid INTEGER,
    title TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    episodecount INTEGER NOT NULL,
    startdate INTEGER,
    enddate INTEGER,
    PRIMARY KEY (aid)
)`)
	if err != nil {
		return err
	}
	_, err = t.Exec("CREATE INDEX anime_titles ON anime (title)")
	if err != nil {
		return err
	}
	_, err = t.Exec(`
CREATE TABLE IF NOT EXISTS "episode" (
    id INTEGER,
    aid INTEGER NOT NULL,
    type INTEGER NOT NULL,
    number INTEGER NOT NULL,
    title TEXT NOT NULL,
    length INTEGER NOT NULL,
    user_watched INTEGER NOT NULL CHECK (user_watched IN (0, 1))
	DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE (aid, type, number),
    FOREIGN KEY (aid) REFERENCES anime (aid)
	ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (type) REFERENCES episode_type (id)
	ON DELETE RESTRICT ON UPDATE CASCADE
)`)
	if err != nil {
		return err
	}
	_, err = t.Exec(`
CREATE TABLE watching (
    aid INTEGER,
    regexp TEXT NOT NULL,
    PRIMARY KEY (aid),
    FOREIGN KEY (aid) REFERENCES anime (aid)
	ON DELETE CASCADE ON UPDATE CASCADE
)`)
	if err != nil {
		return err
	}
	if err := t.Commit(); err != nil {
		return err
	}
	if err := setUserVersion(d, 3); err != nil {
		return err
	}
	return nil
}
