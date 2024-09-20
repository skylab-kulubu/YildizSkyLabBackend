-- name: CreateTeamMember :one
INSERT INTO team_users (team_id,user_id,role,created_at,updated_at) values ($1,$2,$3,NOW(),NOW()) RETURNING *;

-- name: DeleteTeamMember :exec
UPDATE team_users SET deleted_at = NOW() WHERE team_id = $1 AND user_id = $2;

-- name: DeleteTeamMemberByTeamId :exec
UPDATE team_users SET deleted_at = NOW() WHERE team_id = $1;

-- name: DeleteTeamMemberByUserId :exec
UPDATE team_users SET deleted_at = NOW() WHERE user_id = $1;

-- name: GetTeamMember :one
SELECT * FROM team_users WHERE team_id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: GetTeamLeadByTeamId :many
SELECT user_id FROM team_users WHERE team_id = $1 AND role = 'lead' AND  deleted_at IS NULL;

-- name: GetTeamsByUserId :many
SELECT team_id FROM team_users where user_id = $1 and deleted_at is null;
