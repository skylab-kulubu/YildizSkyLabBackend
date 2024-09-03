-- name: GetAllUsers :many
SELECT
    u.id AS user_id,
	u.name,
	u.last_name,
	u.email,
	u.password,
	u.telephone_number,
	u.role,
	u.university,
	u.department,
	u.date_of_birth,
	u.active,
	COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names,
	COALESCE(STRING_AGG(DISTINCT p.name, ',')) AS project_names
FROM
	users u
LEFT JOIN
	team_users ut on u.id =ut.user_id
LEFT JOIN
	teams t on ut.team_id = t.id
LEFT JOIN
    project_users up on u.id = up.user_id
LEFT JOIN
    projects p on up.project_id = p.id
WHERE
    u.deleted_at IS NULL
GROUP BY
	u.id, u.email
LIMIT $1 OFFSET $2;

-- name: GetUser :one
SELECT
    u.id AS user_id,
	u.name,
	u.last_name,
	u.email,
	u.password,
	u.telephone_number,
	u.role,
	u.university,
	u.department,
	u.date_of_birth,
	u.active,
	COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names,
	COALESCE(STRING_AGG(DISTINCT p.name, ',')) AS project_names
FROM
	users u
LEFT JOIN
	team_users ut on u.id =ut.user_id
LEFT JOIN
	teams t on ut.team_id = t.id
LEFT JOIN
    project_users up on u.id = up.user_id
LEFT JOIN
    projects p on up.project_id = p.id
WHERE
    u.id = $1 AND u.deleted_at IS NULL
GROUP BY
	u.id, u.email;


-- name: GetUserByEmail :one
SELECT
    u.id AS user_id,
		u.name,
		u.last_name,
		u.email,
		u.password,
		u.telephone_number,
		u.role,
		u.university,
		u.department,
		u.date_of_birth,
		u.active,
		COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names,
		COALESCE(STRING_AGG(DISTINCT p.name, ',')) AS project_names
FROM
		users u
LEFT JOIN
		team_users ut on u.id =ut.user_id
LEFT JOIN
		teams t on ut.team_id = t.id
LEFT JOIN
    project_users up on u.id = up.user_id
LEFT JOIN
    projects p on up.project_id = p.id
WHERE
    u.email = $1 AND u.deleted_at IS NULL
GROUP BY
		u.id, u.email;


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
    active,
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
    $10,
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
    active = $11,
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
    active = $11,
    created_at = NOW(),
    updated_at = NOW(),
    deleted_at = NULL
WHERE
    id = $1
returning *;
