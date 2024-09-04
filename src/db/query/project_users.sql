-- name: CreateProjectMember :one
INSERT INTO project_users(user_id, project_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING *;

-- name: DeleteProjectMember :exec
UPDATE project_users SET deleted_at = NOW() WHERE user_id = $1 AND project_id = $2;

-- name: DeleteProjectMemberByProjectId :exec
UPDATE project_users SET deleted_at = NOW() WHERE project_id = $1;

-- name: DeleteProjectMemberByUserId :exec 
UPDATE project_users SET deleted_at = NOW() WHERE user_id= $1;
