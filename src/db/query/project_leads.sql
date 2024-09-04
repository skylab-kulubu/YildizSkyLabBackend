-- name: CreateProjectLead :one
INSERT INTO project_leads(project_id, user_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING *;

-- name: DeleteProjectLead :exec
UPDATE project_leads SET deleted_at = NOW() WHERE project_id = $1 AND user_id = $2;

-- name: GetProjectLead :one
SELECT * FROM project_leads WHERE project_id = $1 AND user_id = $2;
