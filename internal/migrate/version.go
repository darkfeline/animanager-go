package migrate

import (
	"database/sql"
	"fmt"
)

func getUserVersion(d *sql.DB) (int, error) {
	r, err := d.Query("PRAGMA user_version")
	if err != nil {
		return 0, err
	}
	defer r.Close()
	ok := r.Next()
	if !ok {
		return 0, r.Err()
	}
	var v int
	if err := r.Scan(&v); err != nil {
		return 0, err
	}
	r.Close()
	if err := r.Err(); err != nil {
		return 0, err
	}
	return v, nil
}

func setUserVersion(d *sql.DB, v int) error {
	_, err := d.Exec(fmt.Sprintf("PRAGMA user_version=%d", v))
	if err != nil {
		return err
	}
	return nil
}
