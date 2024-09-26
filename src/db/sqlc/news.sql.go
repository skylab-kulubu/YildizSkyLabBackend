// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: news.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const createNews = `-- name: CreateNews :one
INSERT INTO news (title, publish_date, description, cover_image_id, created_by_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, publish_date, description, cover_image_id, created_by_id
`

type CreateNewsParams struct {
	Title        string        `json:"title"`
	PublishDate  time.Time     `json:"publish_date"`
	Description  string        `json:"description"`
	CoverImageID sql.NullInt32 `json:"cover_image_id"`
	CreatedByID  sql.NullInt32 `json:"created_by_id"`
}

type CreateNewsRow struct {
	ID           int32         `json:"id"`
	Title        string        `json:"title"`
	PublishDate  time.Time     `json:"publish_date"`
	Description  string        `json:"description"`
	CoverImageID sql.NullInt32 `json:"cover_image_id"`
	CreatedByID  sql.NullInt32 `json:"created_by_id"`
}

func (q *Queries) CreateNews(ctx context.Context, arg CreateNewsParams) (CreateNewsRow, error) {
	row := q.db.QueryRowContext(ctx, createNews,
		arg.Title,
		arg.PublishDate,
		arg.Description,
		arg.CoverImageID,
		arg.CreatedByID,
	)
	var i CreateNewsRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PublishDate,
		&i.Description,
		&i.CoverImageID,
		&i.CreatedByID,
	)
	return i, err
}

const deleteNews = `-- name: DeleteNews :exec
DELETE FROM news WHERE id = $1
`

func (q *Queries) DeleteNews(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteNews, id)
	return err
}

const getANewsWithDetails = `-- name: GetANewsWithDetails :one
SELECT 
    n.id,
    n.title,
    n.publish_date,
    n.description,
    i.id as image_id,
    i.url as image_url,
    i.type as image_type,
    u.id as user_id,
    u.name as user_name,
    u.last_name as user_last_name,
    u.email as user_email,
    u.university as user_university,
    u.department as user_department
FROM news n
JOIN images i ON i.id = n.cover_image_id
JOIN users u ON u.id = n.created_by_id
WHERE n.id = $1
`

type GetANewsWithDetailsRow struct {
	ID             int32     `json:"id"`
	Title          string    `json:"title"`
	PublishDate    time.Time `json:"publish_date"`
	Description    string    `json:"description"`
	ImageID        int32     `json:"image_id"`
	ImageUrl       string    `json:"image_url"`
	ImageType      string    `json:"image_type"`
	UserID         int32     `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserLastName   string    `json:"user_last_name"`
	UserEmail      string    `json:"user_email"`
	UserUniversity string    `json:"user_university"`
	UserDepartment string    `json:"user_department"`
}

func (q *Queries) GetANewsWithDetails(ctx context.Context, id int32) (GetANewsWithDetailsRow, error) {
	row := q.db.QueryRowContext(ctx, getANewsWithDetails, id)
	var i GetANewsWithDetailsRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PublishDate,
		&i.Description,
		&i.ImageID,
		&i.ImageUrl,
		&i.ImageType,
		&i.UserID,
		&i.UserName,
		&i.UserLastName,
		&i.UserEmail,
		&i.UserUniversity,
		&i.UserDepartment,
	)
	return i, err
}

const getAllNews = `-- name: GetAllNews :many
SELECT id, title, publish_date, description, cover_image_id, created_by_id
FROM news
ORDER BY id DESC
LIMIT $1 OFFSET $2
`

type GetAllNewsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAllNewsRow struct {
	ID           int32         `json:"id"`
	Title        string        `json:"title"`
	PublishDate  time.Time     `json:"publish_date"`
	Description  string        `json:"description"`
	CoverImageID sql.NullInt32 `json:"cover_image_id"`
	CreatedByID  sql.NullInt32 `json:"created_by_id"`
}

func (q *Queries) GetAllNews(ctx context.Context, arg GetAllNewsParams) ([]GetAllNewsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllNews, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllNewsRow{}
	for rows.Next() {
		var i GetAllNewsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.PublishDate,
			&i.Description,
			&i.CoverImageID,
			&i.CreatedByID,
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

const getNews = `-- name: GetNews :one
SELECT id, title, publish_date, description, cover_image_id, created_by_id
FROM news
WHERE id = $1
`

type GetNewsRow struct {
	ID           int32         `json:"id"`
	Title        string        `json:"title"`
	PublishDate  time.Time     `json:"publish_date"`
	Description  string        `json:"description"`
	CoverImageID sql.NullInt32 `json:"cover_image_id"`
	CreatedByID  sql.NullInt32 `json:"created_by_id"`
}

func (q *Queries) GetNews(ctx context.Context, id int32) (GetNewsRow, error) {
	row := q.db.QueryRowContext(ctx, getNews, id)
	var i GetNewsRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.PublishDate,
		&i.Description,
		&i.CoverImageID,
		&i.CreatedByID,
	)
	return i, err
}

const getNewsWithDetails = `-- name: GetNewsWithDetails :many
SELECT 
    n.id,
    n.title,
    n.publish_date,
    n.description,
    i.id as image_id,
    i.url as image_url,
    i.type as image_type,
    u.id as user_id,
    u.name as user_name,
    u.last_name as user_last_name,
    u.email as user_email,
    u.university as user_university,
    u.department as user_department
FROM news n
JOIN images i ON i.id = n.cover_image_id
JOIN users u ON u.id = n.created_by_id
ORDER BY n.id DESC
LIMIT $1 OFFSET $2
`

type GetNewsWithDetailsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetNewsWithDetailsRow struct {
	ID             int32     `json:"id"`
	Title          string    `json:"title"`
	PublishDate    time.Time `json:"publish_date"`
	Description    string    `json:"description"`
	ImageID        int32     `json:"image_id"`
	ImageUrl       string    `json:"image_url"`
	ImageType      string    `json:"image_type"`
	UserID         int32     `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserLastName   string    `json:"user_last_name"`
	UserEmail      string    `json:"user_email"`
	UserUniversity string    `json:"user_university"`
	UserDepartment string    `json:"user_department"`
}

func (q *Queries) GetNewsWithDetails(ctx context.Context, arg GetNewsWithDetailsParams) ([]GetNewsWithDetailsRow, error) {
	rows, err := q.db.QueryContext(ctx, getNewsWithDetails, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetNewsWithDetailsRow{}
	for rows.Next() {
		var i GetNewsWithDetailsRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.PublishDate,
			&i.Description,
			&i.ImageID,
			&i.ImageUrl,
			&i.ImageType,
			&i.UserID,
			&i.UserName,
			&i.UserLastName,
			&i.UserEmail,
			&i.UserUniversity,
			&i.UserDepartment,
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

const updateNews = `-- name: UpdateNews :exec
UPDATE news
SET title = $2, publish_date = $3, description = $4, cover_image_id = $5, created_by_id = $6
WHERE id = $1
`

type UpdateNewsParams struct {
	ID           int32         `json:"id"`
	Title        string        `json:"title"`
	PublishDate  time.Time     `json:"publish_date"`
	Description  string        `json:"description"`
	CoverImageID sql.NullInt32 `json:"cover_image_id"`
	CreatedByID  sql.NullInt32 `json:"created_by_id"`
}

func (q *Queries) UpdateNews(ctx context.Context, arg UpdateNewsParams) error {
	_, err := q.db.ExecContext(ctx, updateNews,
		arg.ID,
		arg.Title,
		arg.PublishDate,
		arg.Description,
		arg.CoverImageID,
		arg.CreatedByID,
	)
	return err
}
