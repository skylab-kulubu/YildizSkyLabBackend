// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: teams.sql

package sqlc

import (
	"context"
	"database/sql"
)

const createTeam = `-- name: CreateTeam :one
INSERT INTO teams (
    name,
    description,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    NOW(),
    NOW()
) RETURNING id, name, description, created_at, updated_at, deleted_at
`

type CreateTeamParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateTeam(ctx context.Context, arg CreateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, createTeam, arg.Name, arg.Description)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteTeam = `-- name: DeleteTeam :exec
UPDATE teams set
    deleted_at = NOW()
WHERE
    id = $1
`

func (q *Queries) DeleteTeam(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteTeam, id)
	return err
}

const getAllTeams = `-- name: GetAllTeams :many
SELECT id, name, description, created_at, updated_at, deleted_at FROM teams WHERE deleted_at IS NULL LIMIT $1 OFFSET $2
`

type GetAllTeamsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllTeams(ctx context.Context, arg GetAllTeamsParams) ([]Team, error) {
	rows, err := q.db.QueryContext(ctx, getAllTeams, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Team{}
	for rows.Next() {
		var i Team
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamWithDetails = `-- name: GetTeamWithDetails :many
SELECT
    t.id,
    t.name,
    t.description,
    l.id AS lead_id,
    l.name AS lead_name,
    l.last_name AS lead_last_name,
    l.email AS lead_email,
    l.telephone_number AS lead_telephone_number,
    l.university AS lead_university,
    l.department AS lead_department,
    l.date_of_birth AS lead_date_of_birth,
    p.id AS project_id,
    p.name AS project_name,
    p.description AS project_description,
    u.id AS member_id,
    u.name AS member_name,
    u.last_name AS member_last_name,
    u.email AS member_email,
    u.telephone_number AS member_telephone_number,
    u.university AS member_university,
    u.department AS member_department,
    u.date_of_birth AS member_date_of_birth
FROM teams t
LEFT JOIN team_users tu ON t.id = tu.team_id AND tu.role = 'lead' AND tu.deleted_at IS NULL
LEFT JOIN users l ON tu.user_id = l.id AND l.deleted_at IS NULL
LEFT JOIN team_projects tp ON t.id = tp.team_id AND tp.deleted_at IS NULL
LEFT JOIN projects p ON tp.project_id = p.id AND p.deleted_at IS NULL
LEFT JOIN team_users tm ON t.id = tm.team_id AND tm.role = 'member' AND tm.deleted_at IS NULL
LEFT JOIN users u ON tm.user_id = u.id AND u.deleted_at IS NULL
WHERE t.id = $1
GROUP BY t.id, u.id, p.id, l.id
`

type GetTeamWithDetailsRow struct {
	ID                    int32          `json:"id"`
	Name                  string         `json:"name"`
	Description           string         `json:"description"`
	LeadID                sql.NullInt32  `json:"lead_id"`
	LeadName              sql.NullString `json:"lead_name"`
	LeadLastName          sql.NullString `json:"lead_last_name"`
	LeadEmail             sql.NullString `json:"lead_email"`
	LeadTelephoneNumber   sql.NullString `json:"lead_telephone_number"`
	LeadUniversity        sql.NullString `json:"lead_university"`
	LeadDepartment        sql.NullString `json:"lead_department"`
	LeadDateOfBirth       sql.NullTime   `json:"lead_date_of_birth"`
	ProjectID             sql.NullInt32  `json:"project_id"`
	ProjectName           sql.NullString `json:"project_name"`
	ProjectDescription    sql.NullString `json:"project_description"`
	MemberID              sql.NullInt32  `json:"member_id"`
	MemberName            sql.NullString `json:"member_name"`
	MemberLastName        sql.NullString `json:"member_last_name"`
	MemberEmail           sql.NullString `json:"member_email"`
	MemberTelephoneNumber sql.NullString `json:"member_telephone_number"`
	MemberUniversity      sql.NullString `json:"member_university"`
	MemberDepartment      sql.NullString `json:"member_department"`
	MemberDateOfBirth     sql.NullTime   `json:"member_date_of_birth"`
}

func (q *Queries) GetTeamWithDetails(ctx context.Context, id int32) ([]GetTeamWithDetailsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTeamWithDetails, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTeamWithDetailsRow{}
	for rows.Next() {
		var i GetTeamWithDetailsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.LeadID,
			&i.LeadName,
			&i.LeadLastName,
			&i.LeadEmail,
			&i.LeadTelephoneNumber,
			&i.LeadUniversity,
			&i.LeadDepartment,
			&i.LeadDateOfBirth,
			&i.ProjectID,
			&i.ProjectName,
			&i.ProjectDescription,
			&i.MemberID,
			&i.MemberName,
			&i.MemberLastName,
			&i.MemberEmail,
			&i.MemberTelephoneNumber,
			&i.MemberUniversity,
			&i.MemberDepartment,
			&i.MemberDateOfBirth,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTeam = `-- name: UpdateTeam :one
UPDATE teams set
    name = $1,
    description = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING id, name, description, created_at, updated_at, deleted_at
`

type UpdateTeamParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          int32  `json:"id"`
}

func (q *Queries) UpdateTeam(ctx context.Context, arg UpdateTeamParams) (Team, error) {
	row := q.db.QueryRowContext(ctx, updateTeam, arg.Name, arg.Description, arg.ID)
	var i Team
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
