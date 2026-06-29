-- name: CreateFeed :one
WITH feed AS (
  INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
  VALUES ($1, $2, $3, $4, $5, $6)
  RETURNING *
), feed_follow AS (
  INSERT INTO feed_follows (user_id, feed_id, created_at, updated_at)
  SELECT feed.user_id, feed.id, feed.created_at, feed.updated_at FROM feed
  RETURNING *
) SELECT feed.* FROM feed;

-- name: GetFeeds :many
SELECT
  feeds.*,
  users.name AS created_by
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at ASC NULLS FIRST LIMIT 1;