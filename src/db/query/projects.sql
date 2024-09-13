-- name: GetAllProjects :many
SELECT
    p.id AS project_id,
    p.name,
    p.description,
    COALESCE(STRING_AGG(DISTINCT u.name, ',')) AS user_names,
    COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names
FROM
    projects p
LEFT JOIN
    project_users up on p.id = up.project_id
LEFT JOIN
    users u on up.user_id = u.id
LEFT JOIN
    team_projects tp on p.id = tp.project_id
LEFT JOIN
    teams t on tp.team_id = t.id
WHERE
    p.deleted_at IS NULL
GROUP BY
    p.id, p.name
LIMIT $1 OFFSET $2;

-- name: GetProject :one
SELECT
    p.id AS project_id,
    p.name,
    p.description,
    COALESCE(STRING_AGG(DISTINCT u.name, ',')) AS user_names,
    COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names
FROM
    projects p
LEFT JOIN
    project_users up on p.id = up.project_id
LEFT JOIN
    users u on up.user_id = u.id
LEFT JOIN
    team_projects tp on p.id = tp.project_id
LEFT JOIN
    teams t on tp.team_id = t.id
WHERE
    p.deleted_at IS NULL AND p.id = $1
GROUP BY
    p.id, p.name;

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
