-- name: CreateTeamLead :one
INSERT INTO team_leads (team_id, user_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING *;

-- name: DeleteTeamLead :exec
UPDATE team_leads SET deleted_at = NOW() WHERE team_id = $1 AND user_id = $2;

-- name: GetTeamLead :one
SELECT * FROM team_leads WHERE team_id = $1 AND user_id = $2;
