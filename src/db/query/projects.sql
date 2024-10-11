-- name: GetAllProjects :many
SELECT * FROM projects WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;

-- name: GetProject :one 
SELECT * FROM projects WHERE deleted_at IS NULL AND id = $1; 

-- name: GetProjectWithDetails :many
SELECT
    p.id,
    p.name,
    p.description,
    l.id AS lead_id,
    l.name AS lead_name,
    l.last_name AS lead_last_name,
    l.email AS lead_email,
    l.telephone_number AS lead_telephone_number,
    l.university AS lead_university,
    l.department AS lead_department,
    l.date_of_birth AS lead_date_of_birth,
    t.id AS team_id,
    t.name AS team_name,
    t.description AS team_description,
    u.id AS member_id,
    u.name AS member_name,
    u.last_name AS member_last_name,
    u.email AS member_email,
    u.telephone_number AS member_telephone_number,
    u.university AS member_university,
    u.department AS member_department,
    u.date_of_birth AS member_date_of_birth
FROM projects p
LEFT JOIN project_users pu ON p.id = pu.project_id AND pu.role = 'lead' AND pu.deleted_at IS NULL
LEFT JOIN users l ON pu.user_id = l.id AND l.deleted_at IS NULL 
LEFT JOIN team_projects tp ON p.id = tp.project_id AND tp.deleted_at IS NULL 
LEFT JOIN teams t ON tp.team_id = t.id AND t.deleted_at IS NULL 
LEFT JOIN project_users pm ON p.id = pm.project_id AND pm.role = 'member' AND pm.deleted_at IS NULL 
LEFT JOIN users u on pm.user_id = u.id AND u.deleted_at IS NULL 
WHERE p.id = $1 
GROUP BY t.id, u.id, p.id, l.id;

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
