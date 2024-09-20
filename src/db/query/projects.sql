-- name: GetAllProjects :many
SELECT * FROM projects WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateProject :one
INSERT INTO projects (
    name,
    description,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    NOW(),
    NOW()
) RETURNING *;

-- name: UpdateProject :one
UPDATE projects set
    name = $1,
    description = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING *;

-- name: DeleteProject :exec
UPDATE projects set
    deleted_at = NOW()
WHERE
    id = $1;
