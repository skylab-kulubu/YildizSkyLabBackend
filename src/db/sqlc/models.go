// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql"
	"time"
)

type Announcement struct {
	ID        int32        `json:"id"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	AuthorID  int32        `json:"author_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Image struct {
	ID        int32        `json:"id"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Data      []byte       `json:"data"`
	Url       string       `json:"url"`
	CreatedBy int32        `json:"created_by"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type News struct {
	ID           int32         `json:"id"`
	Title        string        `json:"title"`
	PublishDate  time.Time     `json:"publish_date"`
	Description  string        `json:"description"`
	CoverImageID sql.NullInt32 `json:"cover_image_id"`
	CreatedByID  sql.NullInt32 `json:"created_by_id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	DeletedAt    sql.NullTime  `json:"deleted_at"`
}

type Project struct {
	ID          int32        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type ProjectUser struct {
	ID        int32        `json:"id"`
	ProjectID int32        `json:"project_id"`
	UserID    int32        `json:"user_id"`
	Role      string       `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Team struct {
	ID          int32        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at"`
}

type TeamProject struct {
	ID        int32        `json:"id"`
	TeamID    int32        `json:"team_id"`
	ProjectID int32        `json:"project_id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type TeamUser struct {
	ID        int32        `json:"id"`
	TeamID    int32        `json:"team_id"`
	UserID    int32        `json:"user_id"`
	Role      string       `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type User struct {
	ID              int32        `json:"id"`
	Name            string       `json:"name"`
	LastName        string       `json:"last_name"`
	Email           string       `json:"email"`
	Password        string       `json:"password"`
	TelephoneNumber string       `json:"telephone_number"`
	University      string       `json:"university"`
	Department      string       `json:"department"`
	DateOfBirth     time.Time    `json:"date_of_birth"`
	Role            string       `json:"role"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
}
