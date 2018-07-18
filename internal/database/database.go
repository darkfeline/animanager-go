package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"go.felesatra.moe/animanager/internal/migrate"
)

func Open(p string) (d *sql.DB, err error) {
	d, err = sql.Open("sqlite3", p)
	if err != nil {
		return nil, err
	}
	defer func(d *sql.DB) {
		if err != nil {
			d.Close()
		}
	}(d)
	if err := migrate.Migrate(d); err != nil {
		return nil, errors.Wrap(err, "migrate database")
	}
	return d, nil
}
