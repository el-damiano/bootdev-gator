-- name: CreateFeedFollow :one
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    inserted_feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follows
INNER JOIN users
    ON inserted_feed_follows.user_id = users.id
INNER JOIN feeds
    ON inserted_feed_follows.feed_id = feeds.id
;

-- name: GetFeedFollowsForUser :many
SELECT
    *,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN users
    ON feed_follows.user_id = users.id
INNER JOIN feeds
    ON feed_follows.feed_id = feeds.id
WHERE users.ID = $1;

-- name: DeleteFeedFollow :one
DELETE FROM feed_follows
    USING feeds
WHERE feed_follows.feed_id = feeds.id
    AND feed_follows.user_id = $1
    AND feeds.url = $2
RETURNING 1;
