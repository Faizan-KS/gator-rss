-- name: CreateFeedFollow :one
WITH inserted_follow AS(
    INSERT INTO feed_follows(user_id,feed_id)
    VALUES($1, $2)
    ON CONFLICT (user_id, feed_id)
    DO UPDATE SET user_id = feed_follows.user_id
    RETURNING *
)
SELECT ff.id, ff.user_id, u.name AS user_name, ff.feed_id, f.name AS feed_name, f.url AS feed_url  
FROM inserted_follow ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT u.name AS feed_of, f.name AS feed_name, f.url AS feed_url, feed_id 
FROM feed_follows ff
JOIN feeds f ON f.id = ff.feed_id
JOIN users u ON u.id = f.user_id
WHERE ff.user_id = $1;

-- name: DeleteFeedByURL :exec
DELETE FROM feed_follows ff
USING feeds f
WHERE ff.feed_id = f.id
    AND ff.user_id = $1
    AND f.url = $2;

