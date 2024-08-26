-- name: CreateTeamMember :one
INSERT INTO team_users (team_id,user_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING *;

-- name: DeleteTeamMember :exec
UPDATE team_users SET deleted_at = NOW() WHERE team_id = $1 AND user_id = $2;
