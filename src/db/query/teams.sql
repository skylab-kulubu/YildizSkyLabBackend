-- name: GetAllTeams :many
SELECT * FROM teams WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;

-- name: GetTeam :one
SELECT 
        t.id,
        t.name,
        t.description,
        json_agg(
                json_build_object(
                        'user_id', u.id,
                        'user_name', u.name,
                        'user_last_name', u.last_name,
                        'user_email', u.email,
                        'user_password', u.password,
                        'user_telephone_number', u.telephone_number,
                        'user_university', u.university,
                        'user_department', u.department,
                        'user_date_of_birth', u.date_of_birth,
                        'user_role', u.role
                )
        ) as leads,
        json_agg(
                json_build_object(
                        'user_id', u.id,
                        'user_name', u.name,
                        'user_last_name', u.last_name,
                        'user_email', u.email,
                        'user_password', u.password,
                        'user_telephone_number', u.telephone_number,
                        'user_university', u.university,
                        'user_department', u.department,
                        'user_date_of_birth', u.date_of_birth,
                        'user_role', u.role
                )
        ) as members,
        json_agg(
                json_build_object(
                        'project_id', p.id,
                        'projet_name', p.name,
                        'project_description', p.description
                )
        ) as projects
FROM teams t
LEFT JOIN team_users tu ON t.id = tu.team_id AND tu.role = 'lead'
LEFT JOIN users u on tu.user_id = u.id
LEFT JOIN team_projects tp ON t.id = tp.team_id
LEFT JOIN projects p on tp.project_id = p.id
WHERE t.id = $1
GROUP BY t.id;

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
