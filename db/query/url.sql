-- name: CreateURL :one
INSERT INTO urls (
  origin_url,
  short_url
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetURL :one
SELECT * FROM urls
WHERE short_url = $1 LIMIT 1;

-- name: UpdateURL :one
UPDATE urls
SET short_url = $2
WHERE id = $1
RETURNING *;

-- name: DeleteURL :exec
DELETE FROM urls
WHERE id = $1;