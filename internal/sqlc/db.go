// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.deleteAllEpisodeFilesStmt, err = db.PrepareContext(ctx, deleteAllEpisodeFiles); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAllEpisodeFiles: %w", err)
	}
	if q.deleteAnimeFilesStmt, err = db.PrepareContext(ctx, deleteAnimeFiles); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteAnimeFiles: %w", err)
	}
	if q.deleteEpisodeStmt, err = db.PrepareContext(ctx, deleteEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteEpisode: %w", err)
	}
	if q.deleteWatchingStmt, err = db.PrepareContext(ctx, deleteWatching); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteWatching: %w", err)
	}
	if q.getAIDsStmt, err = db.PrepareContext(ctx, getAIDs); err != nil {
		return nil, fmt.Errorf("error preparing query GetAIDs: %w", err)
	}
	if q.getAllAnimeStmt, err = db.PrepareContext(ctx, getAllAnime); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllAnime: %w", err)
	}
	if q.getAllEpisodesStmt, err = db.PrepareContext(ctx, getAllEpisodes); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllEpisodes: %w", err)
	}
	if q.getAllWatchingStmt, err = db.PrepareContext(ctx, getAllWatching); err != nil {
		return nil, fmt.Errorf("error preparing query GetAllWatching: %w", err)
	}
	if q.getAnimeStmt, err = db.PrepareContext(ctx, getAnime); err != nil {
		return nil, fmt.Errorf("error preparing query GetAnime: %w", err)
	}
	if q.getAnimeCountStmt, err = db.PrepareContext(ctx, getAnimeCount); err != nil {
		return nil, fmt.Errorf("error preparing query GetAnimeCount: %w", err)
	}
	if q.getEpisodeStmt, err = db.PrepareContext(ctx, getEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisode: %w", err)
	}
	if q.getEpisodeCountStmt, err = db.PrepareContext(ctx, getEpisodeCount); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodeCount: %w", err)
	}
	if q.getEpisodeFilesStmt, err = db.PrepareContext(ctx, getEpisodeFiles); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodeFiles: %w", err)
	}
	if q.getEpisodesStmt, err = db.PrepareContext(ctx, getEpisodes); err != nil {
		return nil, fmt.Errorf("error preparing query GetEpisodes: %w", err)
	}
	if q.getFileHashStmt, err = db.PrepareContext(ctx, getFileHash); err != nil {
		return nil, fmt.Errorf("error preparing query GetFileHash: %w", err)
	}
	if q.getFileHashBySizeStmt, err = db.PrepareContext(ctx, getFileHashBySize); err != nil {
		return nil, fmt.Errorf("error preparing query GetFileHashBySize: %w", err)
	}
	if q.getWatchedEpisodeCountStmt, err = db.PrepareContext(ctx, getWatchedEpisodeCount); err != nil {
		return nil, fmt.Errorf("error preparing query GetWatchedEpisodeCount: %w", err)
	}
	if q.getWatchedMinutesStmt, err = db.PrepareContext(ctx, getWatchedMinutes); err != nil {
		return nil, fmt.Errorf("error preparing query GetWatchedMinutes: %w", err)
	}
	if q.getWatchingStmt, err = db.PrepareContext(ctx, getWatching); err != nil {
		return nil, fmt.Errorf("error preparing query GetWatching: %w", err)
	}
	if q.getWatchingCountStmt, err = db.PrepareContext(ctx, getWatchingCount); err != nil {
		return nil, fmt.Errorf("error preparing query GetWatchingCount: %w", err)
	}
	if q.insertAnimeStmt, err = db.PrepareContext(ctx, insertAnime); err != nil {
		return nil, fmt.Errorf("error preparing query InsertAnime: %w", err)
	}
	if q.insertEpisodeStmt, err = db.PrepareContext(ctx, insertEpisode); err != nil {
		return nil, fmt.Errorf("error preparing query InsertEpisode: %w", err)
	}
	if q.insertEpisodeFileStmt, err = db.PrepareContext(ctx, insertEpisodeFile); err != nil {
		return nil, fmt.Errorf("error preparing query InsertEpisodeFile: %w", err)
	}
	if q.insertFileHashStmt, err = db.PrepareContext(ctx, insertFileHash); err != nil {
		return nil, fmt.Errorf("error preparing query InsertFileHash: %w", err)
	}
	if q.insertWatchingStmt, err = db.PrepareContext(ctx, insertWatching); err != nil {
		return nil, fmt.Errorf("error preparing query InsertWatching: %w", err)
	}
	if q.updateEpisodeDoneStmt, err = db.PrepareContext(ctx, updateEpisodeDone); err != nil {
		return nil, fmt.Errorf("error preparing query UpdateEpisodeDone: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.deleteAllEpisodeFilesStmt != nil {
		if cerr := q.deleteAllEpisodeFilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAllEpisodeFilesStmt: %w", cerr)
		}
	}
	if q.deleteAnimeFilesStmt != nil {
		if cerr := q.deleteAnimeFilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteAnimeFilesStmt: %w", cerr)
		}
	}
	if q.deleteEpisodeStmt != nil {
		if cerr := q.deleteEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteEpisodeStmt: %w", cerr)
		}
	}
	if q.deleteWatchingStmt != nil {
		if cerr := q.deleteWatchingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteWatchingStmt: %w", cerr)
		}
	}
	if q.getAIDsStmt != nil {
		if cerr := q.getAIDsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAIDsStmt: %w", cerr)
		}
	}
	if q.getAllAnimeStmt != nil {
		if cerr := q.getAllAnimeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllAnimeStmt: %w", cerr)
		}
	}
	if q.getAllEpisodesStmt != nil {
		if cerr := q.getAllEpisodesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllEpisodesStmt: %w", cerr)
		}
	}
	if q.getAllWatchingStmt != nil {
		if cerr := q.getAllWatchingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAllWatchingStmt: %w", cerr)
		}
	}
	if q.getAnimeStmt != nil {
		if cerr := q.getAnimeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAnimeStmt: %w", cerr)
		}
	}
	if q.getAnimeCountStmt != nil {
		if cerr := q.getAnimeCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getAnimeCountStmt: %w", cerr)
		}
	}
	if q.getEpisodeStmt != nil {
		if cerr := q.getEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeStmt: %w", cerr)
		}
	}
	if q.getEpisodeCountStmt != nil {
		if cerr := q.getEpisodeCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeCountStmt: %w", cerr)
		}
	}
	if q.getEpisodeFilesStmt != nil {
		if cerr := q.getEpisodeFilesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodeFilesStmt: %w", cerr)
		}
	}
	if q.getEpisodesStmt != nil {
		if cerr := q.getEpisodesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getEpisodesStmt: %w", cerr)
		}
	}
	if q.getFileHashStmt != nil {
		if cerr := q.getFileHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFileHashStmt: %w", cerr)
		}
	}
	if q.getFileHashBySizeStmt != nil {
		if cerr := q.getFileHashBySizeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getFileHashBySizeStmt: %w", cerr)
		}
	}
	if q.getWatchedEpisodeCountStmt != nil {
		if cerr := q.getWatchedEpisodeCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWatchedEpisodeCountStmt: %w", cerr)
		}
	}
	if q.getWatchedMinutesStmt != nil {
		if cerr := q.getWatchedMinutesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWatchedMinutesStmt: %w", cerr)
		}
	}
	if q.getWatchingStmt != nil {
		if cerr := q.getWatchingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWatchingStmt: %w", cerr)
		}
	}
	if q.getWatchingCountStmt != nil {
		if cerr := q.getWatchingCountStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getWatchingCountStmt: %w", cerr)
		}
	}
	if q.insertAnimeStmt != nil {
		if cerr := q.insertAnimeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertAnimeStmt: %w", cerr)
		}
	}
	if q.insertEpisodeStmt != nil {
		if cerr := q.insertEpisodeStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertEpisodeStmt: %w", cerr)
		}
	}
	if q.insertEpisodeFileStmt != nil {
		if cerr := q.insertEpisodeFileStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertEpisodeFileStmt: %w", cerr)
		}
	}
	if q.insertFileHashStmt != nil {
		if cerr := q.insertFileHashStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertFileHashStmt: %w", cerr)
		}
	}
	if q.insertWatchingStmt != nil {
		if cerr := q.insertWatchingStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertWatchingStmt: %w", cerr)
		}
	}
	if q.updateEpisodeDoneStmt != nil {
		if cerr := q.updateEpisodeDoneStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing updateEpisodeDoneStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                         DBTX
	tx                         *sql.Tx
	deleteAllEpisodeFilesStmt  *sql.Stmt
	deleteAnimeFilesStmt       *sql.Stmt
	deleteEpisodeStmt          *sql.Stmt
	deleteWatchingStmt         *sql.Stmt
	getAIDsStmt                *sql.Stmt
	getAllAnimeStmt            *sql.Stmt
	getAllEpisodesStmt         *sql.Stmt
	getAllWatchingStmt         *sql.Stmt
	getAnimeStmt               *sql.Stmt
	getAnimeCountStmt          *sql.Stmt
	getEpisodeStmt             *sql.Stmt
	getEpisodeCountStmt        *sql.Stmt
	getEpisodeFilesStmt        *sql.Stmt
	getEpisodesStmt            *sql.Stmt
	getFileHashStmt            *sql.Stmt
	getFileHashBySizeStmt      *sql.Stmt
	getWatchedEpisodeCountStmt *sql.Stmt
	getWatchedMinutesStmt      *sql.Stmt
	getWatchingStmt            *sql.Stmt
	getWatchingCountStmt       *sql.Stmt
	insertAnimeStmt            *sql.Stmt
	insertEpisodeStmt          *sql.Stmt
	insertEpisodeFileStmt      *sql.Stmt
	insertFileHashStmt         *sql.Stmt
	insertWatchingStmt         *sql.Stmt
	updateEpisodeDoneStmt      *sql.Stmt
}

func (q *Queries) PrepareDeleteAllEpisodeFiles(ctx context.Context) error {
	var err error
	if q.deleteAllEpisodeFilesStmt, err = q.db.PrepareContext(ctx, deleteAllEpisodeFiles); err != nil {
		return fmt.Errorf("error preparing query DeleteAllEpisodeFiles: %w", err)
	}
	return nil
}

func (q *Queries) PrepareDeleteAnimeFiles(ctx context.Context) error {
	var err error
	if q.deleteAnimeFilesStmt, err = q.db.PrepareContext(ctx, deleteAnimeFiles); err != nil {
		return fmt.Errorf("error preparing query DeleteAnimeFiles: %w", err)
	}
	return nil
}

func (q *Queries) PrepareDeleteEpisode(ctx context.Context) error {
	var err error
	if q.deleteEpisodeStmt, err = q.db.PrepareContext(ctx, deleteEpisode); err != nil {
		return fmt.Errorf("error preparing query DeleteEpisode: %w", err)
	}
	return nil
}

func (q *Queries) PrepareDeleteWatching(ctx context.Context) error {
	var err error
	if q.deleteWatchingStmt, err = q.db.PrepareContext(ctx, deleteWatching); err != nil {
		return fmt.Errorf("error preparing query DeleteWatching: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetAIDs(ctx context.Context) error {
	var err error
	if q.getAIDsStmt, err = q.db.PrepareContext(ctx, getAIDs); err != nil {
		return fmt.Errorf("error preparing query GetAIDs: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetAllAnime(ctx context.Context) error {
	var err error
	if q.getAllAnimeStmt, err = q.db.PrepareContext(ctx, getAllAnime); err != nil {
		return fmt.Errorf("error preparing query GetAllAnime: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetAllEpisodes(ctx context.Context) error {
	var err error
	if q.getAllEpisodesStmt, err = q.db.PrepareContext(ctx, getAllEpisodes); err != nil {
		return fmt.Errorf("error preparing query GetAllEpisodes: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetAllWatching(ctx context.Context) error {
	var err error
	if q.getAllWatchingStmt, err = q.db.PrepareContext(ctx, getAllWatching); err != nil {
		return fmt.Errorf("error preparing query GetAllWatching: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetAnime(ctx context.Context) error {
	var err error
	if q.getAnimeStmt, err = q.db.PrepareContext(ctx, getAnime); err != nil {
		return fmt.Errorf("error preparing query GetAnime: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetAnimeCount(ctx context.Context) error {
	var err error
	if q.getAnimeCountStmt, err = q.db.PrepareContext(ctx, getAnimeCount); err != nil {
		return fmt.Errorf("error preparing query GetAnimeCount: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetEpisode(ctx context.Context) error {
	var err error
	if q.getEpisodeStmt, err = q.db.PrepareContext(ctx, getEpisode); err != nil {
		return fmt.Errorf("error preparing query GetEpisode: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetEpisodeCount(ctx context.Context) error {
	var err error
	if q.getEpisodeCountStmt, err = q.db.PrepareContext(ctx, getEpisodeCount); err != nil {
		return fmt.Errorf("error preparing query GetEpisodeCount: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetEpisodeFiles(ctx context.Context) error {
	var err error
	if q.getEpisodeFilesStmt, err = q.db.PrepareContext(ctx, getEpisodeFiles); err != nil {
		return fmt.Errorf("error preparing query GetEpisodeFiles: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetEpisodes(ctx context.Context) error {
	var err error
	if q.getEpisodesStmt, err = q.db.PrepareContext(ctx, getEpisodes); err != nil {
		return fmt.Errorf("error preparing query GetEpisodes: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetFileHash(ctx context.Context) error {
	var err error
	if q.getFileHashStmt, err = q.db.PrepareContext(ctx, getFileHash); err != nil {
		return fmt.Errorf("error preparing query GetFileHash: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetFileHashBySize(ctx context.Context) error {
	var err error
	if q.getFileHashBySizeStmt, err = q.db.PrepareContext(ctx, getFileHashBySize); err != nil {
		return fmt.Errorf("error preparing query GetFileHashBySize: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetWatchedEpisodeCount(ctx context.Context) error {
	var err error
	if q.getWatchedEpisodeCountStmt, err = q.db.PrepareContext(ctx, getWatchedEpisodeCount); err != nil {
		return fmt.Errorf("error preparing query GetWatchedEpisodeCount: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetWatchedMinutes(ctx context.Context) error {
	var err error
	if q.getWatchedMinutesStmt, err = q.db.PrepareContext(ctx, getWatchedMinutes); err != nil {
		return fmt.Errorf("error preparing query GetWatchedMinutes: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetWatching(ctx context.Context) error {
	var err error
	if q.getWatchingStmt, err = q.db.PrepareContext(ctx, getWatching); err != nil {
		return fmt.Errorf("error preparing query GetWatching: %w", err)
	}
	return nil
}

func (q *Queries) PrepareGetWatchingCount(ctx context.Context) error {
	var err error
	if q.getWatchingCountStmt, err = q.db.PrepareContext(ctx, getWatchingCount); err != nil {
		return fmt.Errorf("error preparing query GetWatchingCount: %w", err)
	}
	return nil
}

func (q *Queries) PrepareInsertAnime(ctx context.Context) error {
	var err error
	if q.insertAnimeStmt, err = q.db.PrepareContext(ctx, insertAnime); err != nil {
		return fmt.Errorf("error preparing query InsertAnime: %w", err)
	}
	return nil
}

func (q *Queries) PrepareInsertEpisode(ctx context.Context) error {
	var err error
	if q.insertEpisodeStmt, err = q.db.PrepareContext(ctx, insertEpisode); err != nil {
		return fmt.Errorf("error preparing query InsertEpisode: %w", err)
	}
	return nil
}

func (q *Queries) PrepareInsertEpisodeFile(ctx context.Context) error {
	var err error
	if q.insertEpisodeFileStmt, err = q.db.PrepareContext(ctx, insertEpisodeFile); err != nil {
		return fmt.Errorf("error preparing query InsertEpisodeFile: %w", err)
	}
	return nil
}

func (q *Queries) PrepareInsertFileHash(ctx context.Context) error {
	var err error
	if q.insertFileHashStmt, err = q.db.PrepareContext(ctx, insertFileHash); err != nil {
		return fmt.Errorf("error preparing query InsertFileHash: %w", err)
	}
	return nil
}

func (q *Queries) PrepareInsertWatching(ctx context.Context) error {
	var err error
	if q.insertWatchingStmt, err = q.db.PrepareContext(ctx, insertWatching); err != nil {
		return fmt.Errorf("error preparing query InsertWatching: %w", err)
	}
	return nil
}

func (q *Queries) PrepareUpdateEpisodeDone(ctx context.Context) error {
	var err error
	if q.updateEpisodeDoneStmt, err = q.db.PrepareContext(ctx, updateEpisodeDone); err != nil {
		return fmt.Errorf("error preparing query UpdateEpisodeDone: %w", err)
	}
	return nil
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                         tx,
		tx:                         tx,
		deleteAllEpisodeFilesStmt:  q.deleteAllEpisodeFilesStmt,
		deleteAnimeFilesStmt:       q.deleteAnimeFilesStmt,
		deleteEpisodeStmt:          q.deleteEpisodeStmt,
		deleteWatchingStmt:         q.deleteWatchingStmt,
		getAIDsStmt:                q.getAIDsStmt,
		getAllAnimeStmt:            q.getAllAnimeStmt,
		getAllEpisodesStmt:         q.getAllEpisodesStmt,
		getAllWatchingStmt:         q.getAllWatchingStmt,
		getAnimeStmt:               q.getAnimeStmt,
		getAnimeCountStmt:          q.getAnimeCountStmt,
		getEpisodeStmt:             q.getEpisodeStmt,
		getEpisodeCountStmt:        q.getEpisodeCountStmt,
		getEpisodeFilesStmt:        q.getEpisodeFilesStmt,
		getEpisodesStmt:            q.getEpisodesStmt,
		getFileHashStmt:            q.getFileHashStmt,
		getFileHashBySizeStmt:      q.getFileHashBySizeStmt,
		getWatchedEpisodeCountStmt: q.getWatchedEpisodeCountStmt,
		getWatchedMinutesStmt:      q.getWatchedMinutesStmt,
		getWatchingStmt:            q.getWatchingStmt,
		getWatchingCountStmt:       q.getWatchingCountStmt,
		insertAnimeStmt:            q.insertAnimeStmt,
		insertEpisodeStmt:          q.insertEpisodeStmt,
		insertEpisodeFileStmt:      q.insertEpisodeFileStmt,
		insertFileHashStmt:         q.insertFileHashStmt,
		insertWatchingStmt:         q.insertWatchingStmt,
		updateEpisodeDoneStmt:      q.updateEpisodeDoneStmt,
	}
}
