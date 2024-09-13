-- name: GetAllTeams :many
SELECT * FROM teams WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;

-- name: GetTeam :one
SELECT * FROM teams WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateTeam :one
INSERT INTO teams (
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

-- name: UpdateTeam :one
UPDATE teams set
    name = $1,
    description = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING *;

-- name: DeleteTeam :exec
UPDATE teams set
    deleted_at = NOW()
WHERE
    id = $1;
