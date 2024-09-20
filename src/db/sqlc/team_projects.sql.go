// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: team_projects.sql

package sqlc

import (
	"context"
)

const createTeamProject = `-- name: CreateTeamProject :one
INSERT INTO team_projects(team_id, project_id,created_at,updated_at) values ($1,$2,NOW(),NOW()) RETURNING id, team_id, project_id, created_at, updated_at, deleted_at
`

type CreateTeamProjectParams struct {
	TeamID    int32 `json:"team_id"`
	ProjectID int32 `json:"project_id"`
}

func (q *Queries) CreateTeamProject(ctx context.Context, arg CreateTeamProjectParams) (TeamProject, error) {
	row := q.db.QueryRowContext(ctx, createTeamProject, arg.TeamID, arg.ProjectID)
	var i TeamProject
	err := row.Scan(
		&i.ID,
		&i.TeamID,
		&i.ProjectID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteTeamProject = `-- name: DeleteTeamProject :exec
UPDATE team_projects SET deleted_at = NOW() WHERE team_id = $1 AND project_id = $2
`

type DeleteTeamProjectParams struct {
	TeamID    int32 `json:"team_id"`
	ProjectID int32 `json:"project_id"`
}

func (q *Queries) DeleteTeamProject(ctx context.Context, arg DeleteTeamProjectParams) error {
	_, err := q.db.ExecContext(ctx, deleteTeamProject, arg.TeamID, arg.ProjectID)
	return err
}

const getProjectTeamByProjectId = `-- name: GetProjectTeamByProjectId :many
SELECT team_id FROM team_projects WHERE project_id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetProjectTeamByProjectId(ctx context.Context, projectID int32) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getProjectTeamByProjectId, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var team_id int32
		if err := rows.Scan(&team_id); err != nil {
			return nil, err
		}
		items = append(items, team_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamProjectByTeamId = `-- name: GetTeamProjectByTeamId :many
SELECT project_id FROM team_projects WHERE team_id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetTeamProjectByTeamId(ctx context.Context, teamID int32) ([]int32, error) {
	rows, err := q.db.QueryContext(ctx, getTeamProjectByTeamId, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int32{}
	for rows.Next() {
		var project_id int32
		if err := rows.Scan(&project_id); err != nil {
			return nil, err
		}
		items = append(items, project_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
