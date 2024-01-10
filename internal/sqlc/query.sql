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
SELECT SUM(length) FROM episode WHERE user_watched=1;

-- name: GetEpisodeFiles :many
SELECT * FROM episode_file WHERE eid=?;
