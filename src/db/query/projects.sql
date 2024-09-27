-- name: GetAllProjects :many
SELECT * FROM projects WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;

-- name: GetProject :one
SELECT
    p.id,
    p.name,
    p.description,
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
                    'team_id', t.id,
                    'team_name', t.name,
                    'team_description', t.description
            )
    ) as teams
FROM projects p
LEFT JOIN project_users pu ON p.id = pu.project_id AND pu.role = 'lead'
LEFT JOIN users u on pu.user_id = u.id
LEFT JOIN team_projects tp ON p.id = tp.project_id
LEFT JOIN teams t on tp.team_id = t.id
WHERE p.id = $1
GROUP BY p.id;

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
