-- name: CreateTeamProject :one
INSERT INTO team_projects(team_id, project_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING *;

-- name: DeleteTeamProject :exec
UPDATE team_projects SET deleted_at = NOW() WHERE team_id = $1 AND project_id = $2;
