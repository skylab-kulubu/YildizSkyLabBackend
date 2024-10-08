-- name: CreateProjectMember :one
INSERT INTO project_users(user_id, project_id,role,created_at,updated_at) values ($1,$2,$3,NOW(),NOW()) RETURNING *;

-- name: DeleteProjectMember :exec
UPDATE project_users SET deleted_at = NOW() WHERE user_id = $1 AND project_id = $2;

-- name: DeleteProjectMemberByProjectId :exec
UPDATE project_users SET deleted_at = NOW() WHERE project_id = $1;

-- name: DeleteProjectMemberByUserId :exec
UPDATE project_users SET deleted_at = NOW() WHERE user_id= $1;

-- name: GetProjectMember :one
SELECT * FROM project_users WHERE user_id = $1 AND project_id = $2 AND deleted_at IS NULL;

-- name: GetProjectLeadByProjectId :many
SELECT user_id FROM project_users WHERE project_id = $1 AND role = 'lead' AND  deleted_at IS NULL;

-- name: GetProjectsByUserId :many
SELECT project_id FROM project_users where user_id = $1 AND deleted_at is NULL;
