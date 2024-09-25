-- name: CreateNews :one
INSERT INTO news (title, publish_date, description, cover_image_id, created_by_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, publish_date, description, cover_image_id, created_by_id;

-- name: GetNews :one
SELECT id, title, publish_date, description, cover_image_id, created_by_id
FROM news
WHERE id = $1;

-- name: GetAllNews :many
SELECT id, title, publish_date, description, cover_image_id, created_by_id
FROM news
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateNews :exec
UPDATE news
SET title = $2, publish_date = $3, description = $4, cover_image_id = $5, created_by_id = $6
WHERE id = $1;

-- name: DeleteNews :exec
DELETE FROM news WHERE id = $1;

-- name: GetNewsWithDetails :many
SELECT 
    n.id,
    n.title,
    n.publish_date,
    n.description,
    i.id as image_id,
    i.url as image_url,
    i.type as image_type,
    u.id as user_id,
    u.name as user_name,
    u.last_name as user_last_name,
    u.email as user_email,
    u.university as user_university,
    u.department as user_department
FROM news n
JOIN images i ON i.id = n.cover_image_id
JOIN users u ON u.id = n.created_by_id
ORDER BY n.id DESC
LIMIT $1 OFFSET $2;

-- name: GetANewsWithDetails :one
SELECT 
    n.id,
    n.title,
    n.publish_date,
    n.description,
    i.id as image_id,
    i.url as image_url,
    i.type as image_type,
    u.id as user_id,
    u.name as user_name,
    u.last_name as user_last_name,
    u.email as user_email,
    u.university as user_university,
    u.department as user_department
FROM news n
JOIN images i ON i.id = n.cover_image_id
JOIN users u ON u.id = n.created_by_id
WHERE n.id = $1;
