-- name: CreatePost :exec
INSERT INTO posts(title,url,description,published_at,feed_id)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT (url) DO NOTHING;

-- name: GetPostsForUser :many
SELECT title,description,published_at FROM posts
WHERE feed_id = $1
ORDER BY published_at DESC NULLS LAST
LIMIT 2;

-- name: Browse :many
SELECT * FROM posts
ORDER BY published_at DESC NULLS LAST
LIMIT $1;