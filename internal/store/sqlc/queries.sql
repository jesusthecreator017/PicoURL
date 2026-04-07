-- name: SaveURL :exec
INSERT INTO urls (short_url, original_url) VALUES ($1, $2);

-- name: GetOriginalURL :one
SELECT original_url FROM urls WHERE short_url = $1;

-- name: IncrementCount :exec
UPDATE urls SET click_count = click_count + 1 WHERE short_url = $1;

-- name: GetCount :one
SELECT click_count FROM urls WHERE short_url = $1;

-- name: DeleteURL :exec
DELETE FROM urls WHERE short_url = $1;