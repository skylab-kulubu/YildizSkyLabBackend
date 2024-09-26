package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	notExists = 0
	exists    = 1
	deleted   = 2
)

type returnUserResponse struct {
	Id              int32     `json:"id"`
	Name            string    `json:"name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	TelephoneNumber string    `json:"telephone_number"`
	University      string    `json:"university"`
	Department      string    `json:"department"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Role            string    `json:"role"`
}

// SIGNUP
type signupRequest struct {
	Name            string    `json:"name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	TelephoneNumber string    `json:"telephone_number"`
	University      string    `json:"university"`
	Department      string    `json:"department"`
	DateOfBirth     time.Time `json:"date_of_birth"`
}

func (s *Server) signup(c *gin.Context) {
	var req signupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	var user sqlc.User

	isExists, existedUser := s.checkUserIfNotExistByEmail(c, req.Email)

	switch isExists {
	case exists:
		c.JSON(http.StatusConflict, Response{
			IsSuccess: false,
			Message:   "User already exists",
		})
		return
	case notExists:
		user, err = s.query.CreateUser(c, sqlc.CreateUserParams{
			Name:            req.Name,
			LastName:        req.LastName,
			Email:           req.Email,
			Password:        string(hash),
			TelephoneNumber: req.TelephoneNumber,
			University:      req.University,
			Department:      req.Department,
			DateOfBirth:     req.DateOfBirth,
			Role:            "member",
		})

	case deleted:
		user, err = s.overwriteUser(c, existedUser)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "User created successfully",
		Data: returnUserResponse{
			Id:              user.ID,
			Name:            user.Name,
			LastName:        user.LastName,
			Email:           user.Email,
			TelephoneNumber: user.TelephoneNumber,
			University:      user.University,
			Department:      user.Department,
			DateOfBirth:     user.DateOfBirth,
			Role:            user.Role,
		},
	})
}

////////////////////////

// LOGIN USER

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	user, err := s.query.GetUserByEmail(c, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{
				IsSuccess: false,
				Message:   "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Invalid password",
		})
		return
	}

	// generate a jwt toke

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "User loggedin successfully",
		Data:      tokenString,
	})
}

////////////////////////

// GET USER
type getUserRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

type getUserResponseWithDetails struct {
	Id              int32          `json:"id"`
	Name            string         `json:"name"`
	LastName        string         `json:"last_name"`
	Email           string         `json:"email"`
	TelephoneNumber string         `json:"telephone_number"`
	University      string         `json:"university"`
	Department      string         `json:"department"`
	DateOfBirth     time.Time      `json:"date_of_birth"`
	Role            string         `json:"role"`
	Teams           []sqlc.Team    `json:"teams"`
	Projects        []sqlc.Project `json:"projects"`
}

func (s *Server) getUser(c *gin.Context) {
	var req getUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkUserPermission(c, req.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this user",
		})
		return
	}

	user, err := s.query.GetUserWithDetails(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{
				IsSuccess: false,
				Message:   "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	var teams []sqlc.Team
	err = json.Unmarshal(user.Teams, &teams)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	var projects []sqlc.Project
	err = json.Unmarshal(user.Projects, &projects)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "User got successfully",
		Data: getUserResponseWithDetails{
			Id:              user.ID,
			Name:            user.Name,
			LastName:        user.LastName,
			Email:           user.Email,
			TelephoneNumber: user.TelephoneNumber,
			University:      user.University,
			Department:      user.Department,
			DateOfBirth:     user.DateOfBirth,
			Role:            user.Role,
			Teams:           teams,
			Projects:        projects,
		},
	})
}

////////////////////////

// GET ALL USERS
type getAllUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) getAllUsers(c *gin.Context) {
	var req getAllUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	users, err := s.query.GetAllUsers(c, sqlc.GetAllUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	returnUsers := make([]returnUserResponse, len(users))

	for i, user := range users {
		returnUsers[i] = returnUserResponse{
			Id:              user.ID,
			Name:            user.Name,
			LastName:        user.LastName,
			Email:           user.Email,
			TelephoneNumber: user.TelephoneNumber,
			University:      user.University,
			Department:      user.Department,
			DateOfBirth:     user.DateOfBirth,
			Role:            user.Role,
		}
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Users got successfully",
		Data:      returnUsers,
	})
}

////////////////////////

// UPDATE USER
type updateUserRequest struct {
	Name            *string   `json:"name"`
	LastName        *string   `json:"last_name"`
	Email           *string   `json:"email"`
	Password        *string   `json:"password"`
	TelephoneNumber *string   `json:"telephone_number"`
	University      *string   `json:"university"`
	Department      *string   `json:"department"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Role            *string   `json:"role"`
}

func (s *Server) updateUser(c *gin.Context) {
	var id int32

	idParam := c.Param("id")

	i, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	id = int32(i)

	var req updateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkUserPermission(c, id); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to update this user",
		})
		return
	}

	updatedUser, err := s.query.GetUserWithDetails(c, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if req.Name != nil {
		updatedUser.Name = *req.Name
	}

	if req.LastName != nil {
		updatedUser.LastName = *req.LastName
	}

	if req.Email != nil {
		updatedUser.Email = *req.Email
	}

	if req.Password != nil {
		updatedUser.Password = *req.Password
	}

	if req.TelephoneNumber != nil {
		updatedUser.TelephoneNumber = *req.TelephoneNumber
	}

	if req.University != nil {
		updatedUser.University = *req.University
	}

	if req.Department != nil {
		updatedUser.Department = *req.Department
	}
	if !req.DateOfBirth.IsZero() {
		updatedUser.DateOfBirth = req.DateOfBirth
	}
	if req.Role != nil {
		updatedUser.Role = *req.Role
	}
	if req.TelephoneNumber != nil {
		updatedUser.TelephoneNumber = *req.TelephoneNumber
	}
	user, err := s.query.UpdateUser(c, sqlc.UpdateUserParams{
		ID:              id,
		Name:            updatedUser.Name,
		LastName:        updatedUser.LastName,
		Email:           updatedUser.Email,
		Password:        updatedUser.Password,
		TelephoneNumber: updatedUser.TelephoneNumber,
		University:      updatedUser.University,
		Department:      updatedUser.Department,
		DateOfBirth:     updatedUser.DateOfBirth,
		Role:            updatedUser.Role,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "User updated successfully",
		Data: returnUserResponse{
			Id:              user.ID,
			Name:            user.Name,
			LastName:        user.LastName,
			Email:           user.Email,
			TelephoneNumber: user.TelephoneNumber,
			University:      user.University,
			Department:      user.Department,
			DateOfBirth:     user.DateOfBirth,
			Role:            user.Role,
		},
	})
}

////////////////////////

// DELETE USER
type deleteUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteUser(c *gin.Context) {
	var req deleteUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkUserPermission(c, req.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to delete this user",
		})
		return
	}

	err := s.query.DeleteProjectMemberByUserId(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	err = s.query.DeleteTeamMemberByUserId(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	err = s.query.DeleteUser(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "User deleted successfully",
	})
}

////////////////////////

func (s *Server) currentUser(c *gin.Context) {

	anyUser, ok := c.Get("user")

	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "user not found",
		})
	}

	user := anyUser.(sqlc.User)

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "User got successfully",
		Data: returnUserResponse{
			Id:              user.ID,
			Name:            user.Name,
			LastName:        user.LastName,
			Email:           user.Email,
			TelephoneNumber: user.TelephoneNumber,
			University:      user.University,
			Department:      user.Department,
			DateOfBirth:     user.DateOfBirth,
			Role:            user.Role,
		},
	})
}

// UTILS

func (s *Server) checkUserIfNotExistByEmail(c *gin.Context, email string) (int, sqlc.User) {

	user, err := s.query.CheckUserIfExistByEmail(c, email)

	if err == sql.ErrNoRows {
		return notExists, user
	}

	if user.DeletedAt.Valid {
		return deleted, user
	}

	return exists, user
}

func (s *Server) overwriteUser(c *gin.Context, user sqlc.User) (sqlc.User, error) {

	arg := sqlc.OverwriteUserParams{
		ID:              user.ID,
		Name:            user.Name,
		LastName:        user.LastName,
		Email:           user.Email,
		Password:        user.Password,
		TelephoneNumber: user.TelephoneNumber,
		University:      user.University,
		Department:      user.Department,
		DateOfBirth:     user.DateOfBirth,
		Role:            "member",
	}

	return s.query.OverwriteUser(c, arg)

}

func (s *Server) checkUserPermission(c *gin.Context, userId int32) bool {
	anyUser, ok := c.Get("user")
	if !ok {
		return false
	}

	user := anyUser.(sqlc.User)

	if user.Role == "admin" {
		return true
	}

	if user.ID == userId {
		return true
	}

	return false

}
