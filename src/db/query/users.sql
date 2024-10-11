-- name: GetAllUsers :many
Select * from users where deleted_at is null LIMIT $1 OFFSET $2;

-- name: GetUserWithNoDetails :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserWithDetails :many
SELECT
        u.id,
        u.name,
        u.last_name,
        u.email,
        u.telephone_number,
        u.university,
        u.department,
        u.date_of_birth,
        u.role,
        t.id as team_id,
        t.name as team_name,
        t.description as team_description,
        p.id as project_id,
        p.name as project_name,
        p.description as project_description
FROM users u
LEFT JOIN team_users tu ON u.id = tu.user_id
LEFT JOIN teams t ON tu.team_id = t.id
LEFT JOIN project_users pu ON u.id = pu.user_id
LEFT JOIN projects p ON pu.project_id = p.id
WHERE u.id = $1
GROUP BY u.id, t.id, p.id;


-- name: GetUserByEmail :one
SELECT
        u.id,
        u.name,
        u.last_name,
        u.email,
        u.password,
        u.telephone_number,
        u.university,
        u.department,
        u.date_of_birth,
        u.role,
        t.id as team_id,
        t.name as team_name,
        t.description as team_description,
        p.id as project_id,
        p.name as project_name,
        p.description as project_description
FROM users u
LEFT JOIN team_users tu ON u.id = tu.user_id
LEFT JOIN teams t ON tu.team_id = t.id
LEFT JOIN project_users pu ON u.id = pu.user_id
LEFT JOIN projects p ON pu.project_id = p.id
WHERE u.email = $1
GROUP BY u.id, t.id, p.id;

-- name: CheckUserIfExistByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: CreateUser :one
INSERT INTO users (
    name,
    last_name,
    email,
    password,
    telephone_number,
    role,
    university,
    department,
    date_of_birth,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    NOW(),
    NOW()
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    name = $2,
    last_name = $3,
    email = $4,
    password = $5,
    telephone_number = $6,
    role = $7,
    university = $8,
    department = $9,
    date_of_birth = $10,
    updated_at = NOW()
WHERE
    id = $1
returning *;

-- name: DeleteUser :exec
UPDATE users SET
    deleted_at = NOW()
WHERE
    id = $1;


-- name: OverwriteUser :one
UPDATE users SET
    name = $2,
    last_name = $3,
    email = $4,
    password = $5,
    telephone_number = $6,
    role = $7,
    university = $8,
    department = $9,
    date_of_birth = $10,
    created_at = NOW(),
    updated_at = NOW(),
    deleted_at = NULL
WHERE
    id = $1
returning *;
