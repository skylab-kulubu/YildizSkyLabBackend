-- name: GetAllTeams :many
SELECT
    t.id AS team_id,
    t.name,
    t.description,
    COALESCE(STRING_AGG(DISTINCT u.name, ',')) AS lead_names,
    COALESCE(STRING_AGG(DISTINCT p.name, ',')) AS project_names
FROM
    teams t
LEFT JOIN
    team_users tu ON t.id = tu.team_id AND tu.role = 'lead'
LEFT JOIN
    users u ON tu.user_id = u.id
LEFT JOIN
    team_projects tp ON t.id = tp.team_id
LEFT JOIN
    projects p ON tp.project_id = p.id
WHERE
    t.deleted_at IS NULL
GROUP BY
    t.id, t.name
LIMIT $1 OFFSET $2;

-- name: GetTeam :one
SELECT
    t.id AS team_id,
    t.name,
    t.description,
    COALESCE(STRING_AGG(DISTINCT u.name, ',')) AS lead_names,
    COALESCE(STRING_AGG(DISTINCT p.name, ',')) AS project_names
FROM
    teams t
LEFT JOIN
    team_users tu ON t.id = tu.team_id AND tu.role = 'lead'
LEFT JOIN
    users u ON tu.user_id = u.id
LEFT JOIN
    team_projects tp ON t.id = tp.team_id
LEFT JOIN
    projects p ON tp.project_id = p.id
WHERE
    t.deleted_at IS NULL AND t.id = $1
GROUP BY
    t.id, t.name;

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
