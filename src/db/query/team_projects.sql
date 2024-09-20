-- name: CreateTeamProject :one
INSERT INTO team_projects(team_id, project_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING *;

-- name: DeleteTeamProject :exec
UPDATE team_projects SET deleted_at = NOW() WHERE team_id = $1 AND project_id = $2;

-- name: GetTeamProjectByTeamId :many
SELECT project_id FROM team_projects WHERE team_id = $1 AND deleted_at IS NULL;

-- name: GetProjectTeamByProjectId :many
SELECT team_id FROM team_projects WHERE project_id = $1 AND deleted_at IS NULL;
