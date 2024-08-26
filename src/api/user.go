package api

import (
	"database/sql"
	"net/http"
	"time"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Name            string    `json:"name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	TelephoneNumber string    `json:"telephone_number"`
	University      string    `json:"university"`
	Department      string    `json:"department"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Role            string    `json:"role"`
	Active          bool      `json:"active"`
}

func (s *Server) signup(c *gin.Context) {
	var req signupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err := s.query.CreateUser(c, sqlc.CreateUserParams{
		Name:            req.Name,
		LastName:        req.LastName,
		Email:           req.Email,
		Password:        string(hash),
		TelephoneNumber: req.TelephoneNumber,
		University:      req.University,
		Department:      req.Department,
		DateOfBirth:     req.DateOfBirth,
		Role:            req.Role,
		Active:          req.Active,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)

}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Server) login(c *gin.Context) {
	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.query.GetUserByEmail(c, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// generate a jwt toke

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.secret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", tokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, user)

}

type getUserRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (s *Server) getUser(c *gin.Context) {
	var req getUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.query.GetUser(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)

}

type getAllUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) getAllUsers(c *gin.Context) {
	var req getAllUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := s.query.GetAllUsers(c, sqlc.GetAllUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, users)

}

type updateUserRequest struct {
	ID              int32     `json:"id" binding:"required"`
	Name            string    `json:"name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	TelephoneNumber string    `json:"telephone_number"`
	University      string    `json:"university"`
	Department      string    `json:"department"`
	DateOfBirth     time.Time `json:"date_of_birth"`
	Role            string    `json:"role"`
	Active          bool      `json:"active"`
}

func (s *Server) updateUser(c *gin.Context) {
	var req updateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedUser, err := s.query.UpdateUser(c, sqlc.UpdateUserParams{
		ID:              req.ID,
		Name:            req.Name,
		LastName:        req.LastName,
		Email:           req.Email,
		Password:        req.Password,
		TelephoneNumber: req.TelephoneNumber,
		University:      req.University,
		Department:      req.Department,
		DateOfBirth:     req.DateOfBirth,
		Role:            req.Role,
		Active:          req.Active,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	c.JSON(http.StatusOK, updatedUser)
}

type deleteUserRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteUser(c *gin.Context) {
	var req deleteUserRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteUser(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
