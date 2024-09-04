package api

import (
	"fmt"
	"net/http"
	"time"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func (s *Server) RequireAuth(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Auth")

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.JSON(http.StatusUnauthorized, Response{
				IsSuccess: false,
				Message:   err.Error(),
			})
			return
		}

		id := claims["sub"]

		intID := int32(id.(float64))

		userRow, err := s.query.GetUser(ctx, intID)

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
			Active:          userRow.Active,
		}

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", user)

		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}

func (s *Server) RequireRole(roles []string, f gin.HandlerFunc) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		anyUser, ok := ctx.Get("user")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user := anyUser.(sqlc.User)

		for _, v := range roles {
			if user.Role == v {
				f(ctx)
				return
			}
		}

		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}
