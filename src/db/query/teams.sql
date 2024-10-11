-- name: GetAllTeams :many
SELECT * FROM teams WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;

-- name: GetTeamWithDetails :many
SELECT
    t.id,
    t.name,
    t.description,
    l.id AS lead_id,
    l.name AS lead_name,
    l.last_name AS lead_last_name,
    l.email AS lead_email,
    l.telephone_number AS lead_telephone_number,
    l.university AS lead_university,
    l.department AS lead_department,
    l.date_of_birth AS lead_date_of_birth,
    p.id AS project_id,
    p.name AS project_name,
    p.description AS project_description,
    u.id AS member_id,
    u.name AS member_name,
    u.last_name AS member_last_name,
    u.email AS member_email,
    u.telephone_number AS member_telephone_number,
    u.university AS member_university,
    u.department AS member_department,
    u.date_of_birth AS member_date_of_birth
FROM teams t
LEFT JOIN team_users tu ON t.id = tu.team_id AND tu.role = 'lead' AND tu.deleted_at IS NULL
LEFT JOIN users l ON tu.user_id = l.id AND l.deleted_at IS NULL
LEFT JOIN team_projects tp ON t.id = tp.team_id AND tp.deleted_at IS NULL
LEFT JOIN projects p ON tp.project_id = p.id AND p.deleted_at IS NULL
LEFT JOIN team_users tm ON t.id = tm.team_id AND tm.role = 'member' AND tm.deleted_at IS NULL
LEFT JOIN users u ON tm.user_id = u.id AND u.deleted_at IS NULL
WHERE t.id = $1
GROUP BY t.id, u.id, p.id, l.id;


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
