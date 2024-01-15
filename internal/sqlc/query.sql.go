// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
)

const deleteAllEpisodeFiles = `-- name: DeleteAllEpisodeFiles :exec
DELETE FROM episode_file
`

func (q *Queries) DeleteAllEpisodeFiles(ctx context.Context) error {
	_, err := q.exec(ctx, q.deleteAllEpisodeFilesStmt, deleteAllEpisodeFiles)
	return err
}

const deleteAnimeFiles = `-- name: DeleteAnimeFiles :exec
DELETE FROM episode_file WHERE ROWID IN (
    SELECT episode_file.ROWID FROM episode_file
    JOIN episode ON (episode_file.eid = episode.eid)
    WHERE episode.aid=?
)
`

func (q *Queries) DeleteAnimeFiles(ctx context.Context, aid AID) error {
	_, err := q.exec(ctx, q.deleteAnimeFilesStmt, deleteAnimeFiles, aid)
	return err
}

const deleteEpisode = `-- name: DeleteEpisode :exec
DELETE FROM episode WHERE eid = ?
`

func (q *Queries) DeleteEpisode(ctx context.Context, eid EID) error {
	_, err := q.exec(ctx, q.deleteEpisodeStmt, deleteEpisode, eid)
	return err
}

const deleteWatching = `-- name: DeleteWatching :exec
DELETE FROM watching WHERE aid = ?
`

func (q *Queries) DeleteWatching(ctx context.Context, aid AID) error {
	_, err := q.exec(ctx, q.deleteWatchingStmt, deleteWatching, aid)
	return err
}

const getAIDs = `-- name: GetAIDs :many
SELECT aid FROM anime
`

func (q *Queries) GetAIDs(ctx context.Context) ([]AID, error) {
	rows, err := q.query(ctx, q.getAIDsStmt, getAIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AID
	for rows.Next() {
		var aid AID
		if err := rows.Scan(&aid); err != nil {
			return nil, err
		}
		items = append(items, aid)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllAnime = `-- name: GetAllAnime :many
SELECT aid, title, type, episodecount, startdate, enddate FROM anime
`

func (q *Queries) GetAllAnime(ctx context.Context) ([]Anime, error) {
	rows, err := q.query(ctx, q.getAllAnimeStmt, getAllAnime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Anime
	for rows.Next() {
		var i Anime
		if err := rows.Scan(
			&i.Aid,
			&i.Title,
			&i.Type,
			&i.Episodecount,
			&i.Startdate,
			&i.Enddate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllEpisodes = `-- name: GetAllEpisodes :many
SELECT eid, aid, type, number, title, length, user_watched FROM episode
`

func (q *Queries) GetAllEpisodes(ctx context.Context) ([]Episode, error) {
	rows, err := q.query(ctx, q.getAllEpisodesStmt, getAllEpisodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Episode
	for rows.Next() {
		var i Episode
		if err := rows.Scan(
			&i.Eid,
			&i.Aid,
			&i.Type,
			&i.Number,
			&i.Title,
			&i.Length,
			&i.UserWatched,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllWatching = `-- name: GetAllWatching :many
SELECT aid, "regexp", "offset" FROM watching
`

func (q *Queries) GetAllWatching(ctx context.Context) ([]Watching, error) {
	rows, err := q.query(ctx, q.getAllWatchingStmt, getAllWatching)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Watching
	for rows.Next() {
		var i Watching
		if err := rows.Scan(&i.Aid, &i.Regexp, &i.Offset); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAnime = `-- name: GetAnime :one
SELECT aid, title, type, episodecount, startdate, enddate FROM anime WHERE aid = ?
`

func (q *Queries) GetAnime(ctx context.Context, aid AID) (Anime, error) {
	row := q.queryRow(ctx, q.getAnimeStmt, getAnime, aid)
	var i Anime
	err := row.Scan(
		&i.Aid,
		&i.Title,
		&i.Type,
		&i.Episodecount,
		&i.Startdate,
		&i.Enddate,
	)
	return i, err
}

const getAnimeCount = `-- name: GetAnimeCount :one
SELECT COUNT(*) FROM anime
`

func (q *Queries) GetAnimeCount(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.getAnimeCountStmt, getAnimeCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getEpisode = `-- name: GetEpisode :one
SELECT eid, aid, type, number, title, length, user_watched FROM episode WHERE eid = ? LIMIT 1
`

func (q *Queries) GetEpisode(ctx context.Context, eid EID) (Episode, error) {
	row := q.queryRow(ctx, q.getEpisodeStmt, getEpisode, eid)
	var i Episode
	err := row.Scan(
		&i.Eid,
		&i.Aid,
		&i.Type,
		&i.Number,
		&i.Title,
		&i.Length,
		&i.UserWatched,
	)
	return i, err
}

const getEpisodeCount = `-- name: GetEpisodeCount :one
SELECT COUNT(*) FROM episode
`

func (q *Queries) GetEpisodeCount(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.getEpisodeCountStmt, getEpisodeCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getEpisodeFiles = `-- name: GetEpisodeFiles :many
SELECT id, eid, path FROM episode_file WHERE eid=?
`

func (q *Queries) GetEpisodeFiles(ctx context.Context, eid EID) ([]EpisodeFile, error) {
	rows, err := q.query(ctx, q.getEpisodeFilesStmt, getEpisodeFiles, eid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []EpisodeFile
	for rows.Next() {
		var i EpisodeFile
		if err := rows.Scan(&i.ID, &i.Eid, &i.Path); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEpisodes = `-- name: GetEpisodes :many
SELECT eid, aid, type, number, title, length, user_watched FROM episode WHERE aid = ? ORDER BY type, number
`

func (q *Queries) GetEpisodes(ctx context.Context, aid AID) ([]Episode, error) {
	rows, err := q.query(ctx, q.getEpisodesStmt, getEpisodes, aid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Episode
	for rows.Next() {
		var i Episode
		if err := rows.Scan(
			&i.Eid,
			&i.Aid,
			&i.Type,
			&i.Number,
			&i.Title,
			&i.Length,
			&i.UserWatched,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFileHash = `-- name: GetFileHash :one
SELECT size, hash, eid, aid, filename FROM filehash WHERE size=? AND hash=?
`

type GetFileHashParams struct {
	Size int64
	Hash Hash
}

func (q *Queries) GetFileHash(ctx context.Context, arg GetFileHashParams) (Filehash, error) {
	row := q.queryRow(ctx, q.getFileHashStmt, getFileHash, arg.Size, arg.Hash)
	var i Filehash
	err := row.Scan(
		&i.Size,
		&i.Hash,
		&i.Eid,
		&i.Aid,
		&i.Filename,
	)
	return i, err
}

const getWatchedEpisodeCount = `-- name: GetWatchedEpisodeCount :one
SELECT COUNT(*) FROM episode WHERE user_watched=1
`

func (q *Queries) GetWatchedEpisodeCount(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.getWatchedEpisodeCountStmt, getWatchedEpisodeCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getWatchedMinutes = `-- name: GetWatchedMinutes :one
SELECT CAST(SUM(length) AS INT) FROM episode WHERE user_watched=1
`

func (q *Queries) GetWatchedMinutes(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.getWatchedMinutesStmt, getWatchedMinutes)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getWatching = `-- name: GetWatching :one
SELECT aid, "regexp", "offset" FROM watching WHERE aid = ?
`

func (q *Queries) GetWatching(ctx context.Context, aid AID) (Watching, error) {
	row := q.queryRow(ctx, q.getWatchingStmt, getWatching, aid)
	var i Watching
	err := row.Scan(&i.Aid, &i.Regexp, &i.Offset)
	return i, err
}

const getWatchingCount = `-- name: GetWatchingCount :one
SELECT COUNT(*) FROM watching
`

func (q *Queries) GetWatchingCount(ctx context.Context) (int64, error) {
	row := q.queryRow(ctx, q.getWatchingCountStmt, getWatchingCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const insertAnime = `-- name: InsertAnime :exec
INSERT INTO anime (aid, title, type, episodecount, startdate, enddate)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (aid) DO UPDATE SET
title=excluded.title, type=excluded.type, episodecount=excluded.episodecount,
startdate=excluded.startdate, enddate=excluded.enddate
WHERE aid=excluded.aid
`

type InsertAnimeParams struct {
	Aid          AID
	Title        string
	Type         string
	Episodecount int
	Startdate    sql.NullInt64
	Enddate      sql.NullInt64
}

func (q *Queries) InsertAnime(ctx context.Context, arg InsertAnimeParams) error {
	_, err := q.exec(ctx, q.insertAnimeStmt, insertAnime,
		arg.Aid,
		arg.Title,
		arg.Type,
		arg.Episodecount,
		arg.Startdate,
		arg.Enddate,
	)
	return err
}

const insertEpisode = `-- name: InsertEpisode :exec
INSERT INTO episode (eid, aid, type, number, title, length)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (eid) DO UPDATE SET
aid=excluded.aid, type=excluded.type, number=excluded.number,
title=excluded.title, length=excluded.length
WHERE eid=excluded.eid
`

type InsertEpisodeParams struct {
	Eid    EID
	Aid    AID
	Type   EpisodeType
	Number int
	Title  string
	Length int
}

func (q *Queries) InsertEpisode(ctx context.Context, arg InsertEpisodeParams) error {
	_, err := q.exec(ctx, q.insertEpisodeStmt, insertEpisode,
		arg.Eid,
		arg.Aid,
		arg.Type,
		arg.Number,
		arg.Title,
		arg.Length,
	)
	return err
}

const insertEpisodeFile = `-- name: InsertEpisodeFile :exec
INSERT INTO episode_file (eid, path) VALUES (?, ?)
`

type InsertEpisodeFileParams struct {
	Eid  EID
	Path string
}

func (q *Queries) InsertEpisodeFile(ctx context.Context, arg InsertEpisodeFileParams) error {
	_, err := q.exec(ctx, q.insertEpisodeFileStmt, insertEpisodeFile, arg.Eid, arg.Path)
	return err
}

const insertFileHash = `-- name: InsertFileHash :exec
INSERT INTO filehash (size, hash, eid, aid, filename)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (size, hash) DO UPDATE SET
eid=excluded.eid, aid=excluded.aid, filename=excluded.filename
WHERE size=excluded.size AND hash=excluded.hash
`

type InsertFileHashParams struct {
	Size     int64
	Hash     Hash
	Eid      EID
	Aid      AID
	Filename string
}

func (q *Queries) InsertFileHash(ctx context.Context, arg InsertFileHashParams) error {
	_, err := q.exec(ctx, q.insertFileHashStmt, insertFileHash,
		arg.Size,
		arg.Hash,
		arg.Eid,
		arg.Aid,
		arg.Filename,
	)
	return err
}

const insertWatching = `-- name: InsertWatching :exec
INSERT INTO watching (aid, regexp, offset) VALUES (?, ?, ?)
ON CONFLICT (aid) DO UPDATE
SET regexp=excluded.regexp, offset=excluded.offset
WHERE aid=excluded.aid
`

type InsertWatchingParams struct {
	Aid    AID
	Regexp string
	Offset int
}

func (q *Queries) InsertWatching(ctx context.Context, arg InsertWatchingParams) error {
	_, err := q.exec(ctx, q.insertWatchingStmt, insertWatching, arg.Aid, arg.Regexp, arg.Offset)
	return err
}

const updateEpisodeDone = `-- name: UpdateEpisodeDone :exec
UPDATE episode SET user_watched = ? WHERE eid = ?
`

type UpdateEpisodeDoneParams struct {
	UserWatched bool
	Eid         EID
}

func (q *Queries) UpdateEpisodeDone(ctx context.Context, arg UpdateEpisodeDoneParams) error {
	_, err := q.exec(ctx, q.updateEpisodeDoneStmt, updateEpisodeDone, arg.UserWatched, arg.Eid)
	return err
}
