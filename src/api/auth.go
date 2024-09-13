package api

import (
	"fmt"
	"net/http"
	"time"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (s *Server) RequireAuth(c *gin.Context) {
	/*
		tokenString, err := c.Cookie("Auth")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	*/

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.JSON(http.StatusUnauthorized, Response{
				IsSuccess: false,
				Message:   err.Error(),
			})
			return
		}

		id := claims["sub"]

		intID := int32(id.(float64))

		userRow, err := s.query.GetUser(c, intID)

		user := sqlc.User{
			ID:              userRow.UserID,
			Name:            userRow.Name,
			LastName:        userRow.LastName,
			Email:           userRow.Email,
			Password:        userRow.Password,
			TelephoneNumber: userRow.TelephoneNumber,
			University:      userRow.University,
			Department:      userRow.Department,
			DateOfBirth:     userRow.DateOfBirth,
			Role:            userRow.Role,
		}

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}

func (s *Server) RequireRole(roles []string, f gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		anyUser, ok := c.Get("user")
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := anyUser.(sqlc.User)

		for _, v := range roles {
			if user.Role == v {
				f(c)
				return
			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
