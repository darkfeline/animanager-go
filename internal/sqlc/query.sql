-- name: InsertAnime :exec
INSERT INTO anime (aid, title, type, episodecount, startdate, enddate)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (aid) DO UPDATE SET
title=excluded.title, type=excluded.type, episodecount=excluded.episodecount,
startdate=excluded.startdate, enddate=excluded.enddate
WHERE aid=excluded.aid;

-- name: GetAnimeCount :one
SELECT COUNT(*) FROM anime;

-- name: GetAIDs :many
SELECT aid FROM anime;

-- name: GetAnime :one
SELECT * FROM anime WHERE aid = ?;

-- name: GetAllAnime :many
SELECT * FROM anime;

-- name: InsertEpisode :exec
INSERT INTO episode (eid, aid, type, number, title, length)
VALUES (?, ?, ?, ?, ?, ?)
ON CONFLICT (eid) DO UPDATE SET
aid=excluded.aid, type=excluded.type, number=excluded.number,
title=excluded.title, length=excluded.length
WHERE eid=excluded.eid;

-- name: DeleteAnimeFiles :exec
DELETE FROM episode_file WHERE ROWID IN (
    SELECT episode_file.ROWID FROM episode_file
    JOIN episode ON (episode_file.eid = episode.eid)
    WHERE episode.aid=?
);

-- name: GetEpisode :one
SELECT * FROM episode WHERE eid = ? LIMIT 1;

-- name: DeleteEpisode :exec
DELETE FROM episode WHERE eid = ?;

-- name: GetEpisodes :many
SELECT * FROM episode WHERE aid = ? ORDER BY type, number;

-- name: GetEpisodeCount :one
SELECT COUNT(*) FROM episode;

-- name: GetAllEpisodes :many
SELECT * FROM episode;

-- name: UpdateEpisodeDone :exec
UPDATE episode SET user_watched = ? WHERE eid = ?;

-- name: GetWatchedEpisodeCount :one
SELECT COUNT(*) FROM episode WHERE user_watched=1;

-- name: GetWatchedMinutes :one
SELECT CAST(SUM(length) AS INT) FROM episode WHERE user_watched=1;

-- name: GetEpisodeFiles :many
SELECT * FROM episode_file WHERE eid=?;

-- name: InsertEpisodeFile :exec
INSERT INTO episode_file (eid, path) VALUES (?, ?);

-- name: DeleteAllEpisodeFiles :exec
DELETE FROM episode_file;

-- name: InsertFileHash :exec
INSERT INTO filehash (size, hash, eid, aid, filename)
VALUES (?, ?, ?, ?, ?)
ON CONFLICT (size, hash) DO UPDATE SET
eid=excluded.eid, aid=excluded.aid, filename=excluded.filename
WHERE size=excluded.size AND hash=excluded.hash;

-- name: GetFileHash :one
SELECT * FROM filehash WHERE size=? AND hash=?;

-- name: InsertWatching :exec
INSERT INTO watching (aid, regexp, offset) VALUES (?, ?, ?)
ON CONFLICT (aid) DO UPDATE
SET regexp=excluded.regexp, offset=excluded.offset
WHERE aid=excluded.aid;

-- name: GetWatching :one
SELECT * FROM watching WHERE aid = ?;

-- name: GetWatchingCount :one
SELECT COUNT(*) FROM watching;

-- name: GetAllWatching :many
SELECT * FROM watching;

-- name: DeleteWatching :exec
DELETE FROM watching WHERE aid = ?;
