// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: projects.sql

package sqlc

import (
	"context"
)

const createProject = `-- name: CreateProject :one
INSERT INTO projects (
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

type CreateProjectParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (q *Queries) CreateProject(ctx context.Context, arg CreateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, createProject, arg.Name, arg.Description)
	var i Project
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

const deleteProject = `-- name: DeleteProject :exec
UPDATE projects set
    deleted_at = NOW()
WHERE
    id = $1
`

func (q *Queries) DeleteProject(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteProject, id)
	return err
}

const getAllProjects = `-- name: GetAllProjects :many
SELECT
    p.id AS project_id,
    p.name,
    p.description,
    COALESCE(STRING_AGG(DISTINCT u.name, ',')) AS user_names,
    COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names
FROM
    projects p
LEFT JOIN
    project_users up on p.id = up.project_id
LEFT JOIN
    users u on up.user_id = u.id
LEFT JOIN
    team_projects tp on p.id = tp.project_id
LEFT JOIN
    teams t on tp.team_id = t.id
WHERE
    p.deleted_at IS NULL
GROUP BY
    p.id, p.name
LIMIT $1 OFFSET $2
`

type GetAllProjectsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAllProjectsRow struct {
	ProjectID   int32       `json:"project_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	UserNames   interface{} `json:"user_names"`
	TeamNames   interface{} `json:"team_names"`
}

func (q *Queries) GetAllProjects(ctx context.Context, arg GetAllProjectsParams) ([]GetAllProjectsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllProjects, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllProjectsRow{}
	for rows.Next() {
		var i GetAllProjectsRow
		if err := rows.Scan(
			&i.ProjectID,
			&i.Name,
			&i.Description,
			&i.UserNames,
			&i.TeamNames,
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

const getProject = `-- name: GetProject :one
SELECT
    p.id AS project_id,
    p.name,
    p.description,
    COALESCE(STRING_AGG(DISTINCT u.name, ',')) AS user_names,
    COALESCE(STRING_AGG(DISTINCT t.name, ',')) AS team_names
FROM
    projects p
LEFT JOIN
    project_users up on p.id = up.project_id
LEFT JOIN
    users u on up.user_id = u.id
LEFT JOIN
    team_projects tp on p.id = tp.project_id
LEFT JOIN
    teams t on tp.team_id = t.id
WHERE
    p.deleted_at IS NULL AND p.id = $1
GROUP BY
    p.id, p.name
`

type GetProjectRow struct {
	ProjectID   int32       `json:"project_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	UserNames   interface{} `json:"user_names"`
	TeamNames   interface{} `json:"team_names"`
}

func (q *Queries) GetProject(ctx context.Context, id int32) (GetProjectRow, error) {
	row := q.db.QueryRowContext(ctx, getProject, id)
	var i GetProjectRow
	err := row.Scan(
		&i.ProjectID,
		&i.Name,
		&i.Description,
		&i.UserNames,
		&i.TeamNames,
	)
	return i, err
}

const updateProject = `-- name: UpdateProject :one
UPDATE projects set
    name = $1,
    description = $2,
    updated_at = NOW()
WHERE
    id = $3
RETURNING id, name, description, created_at, updated_at, deleted_at
`

type UpdateProjectParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ID          int32  `json:"id"`
}

func (q *Queries) UpdateProject(ctx context.Context, arg UpdateProjectParams) (Project, error) {
	row := q.db.QueryRowContext(ctx, updateProject, arg.Name, arg.Description, arg.ID)
	var i Project
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
