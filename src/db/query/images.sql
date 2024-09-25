-- name: SaveImage :one
INSERT INTO images (type, name, data, url, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, type, name, url, created_by, created_at;

-- name: GetImageByUrl :one
SELECT *
FROM images
WHERE url = $1;