-- name: CreateFeedFollow :one
WITH feed_follow AS (
  INSERT INTO feed_follows (user_id, feed_id, created_at, updated_at)
  VALUES ($1, $2, $3, $4)
  RETURNING *
) SELECT
  feed_follow.*,
  users.name AS user_name,
  feeds.name AS feed_name
FROM feed_follow
INNER JOIN users ON feed_follow.user_id = users.id
INNER JOIN feeds ON feed_follow.feed_id = feeds.id;

-- name: GetFeedFollows :many
SELECT
  feeds.name as feed_name
FROM feeds
INNER JOIN feed_follows ON feed_follows.feed_id = feeds.id AND feed_follows.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;