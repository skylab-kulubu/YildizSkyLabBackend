-- name: SaveImage :one
INSERT INTO images (type, name, data, url, created_by)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetImageByUrl :one
SELECT *
FROM images
WHERE url = $1;
