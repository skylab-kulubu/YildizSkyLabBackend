package api

import (
	"context"
	"time"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	admin = "admin"
	lead  = "lead"
)

type Server struct {
	query  *sqlc.Queries
	router *gin.Engine
	secret string
}

func NewServer(query *sqlc.Queries, secret string) *Server {

	server := &Server{
		query:  query,
		secret: secret,
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/deneme", server.RequireAuth, server.RequireRole([]string{"mod"}, server.getAllTeams))

	//team
	router.POST("/teams", server.RequireAuth, server.RequireRole([]string{admin}, server.createTeam))
	router.GET("/teams/:id", server.RequireAuth, server.RequireRole([]string{admin, lead}, server.getTeam))
	router.GET("/teams", server.RequireAuth, server.getAllTeams)
	router.PUT("/teams/:id", server.RequireAuth, server.RequireRole([]string{admin, lead}, server.updateTeam))
	router.DELETE("/teams/:id", server.RequireAuth, server.RequireRole([]string{admin}, server.deleteTeam))
	router.POST("/teams/project", server.RequireAuth, server.RequireRole([]string{admin, lead}, server.addTeamProject))
	router.DELETE("/teams/project", server.RequireAuth, server.RequireRole([]string{admin, lead}, server.removeTeamProject))
	router.POST("/teams/member", server.RequireAuth, server.RequireRole([]string{admin, lead}, server.addTeamMember))
	router.DELETE("/teams/member", server.RequireAuth, server.RequireRole([]string{admin, lead}, server.removeTeamMember))

	//user
	router.POST("/users/signup", server.signup)
	router.POST("/users/login", server.login)
	router.GET("/users/:id", server.RequireAuth, server.getUser)
	router.GET("/users", server.RequireAuth, server.getAllUsers)
	router.PUT("/users/:id", server.RequireAuth, server.updateUser)
	router.DELETE("/users/:id", server.RequireAuth, server.deleteUser)

	//project
	router.POST("/projects", server.RequireAuth, server.createProject)
	router.GET("/projects/:id", server.RequireAuth, server.getProject)
	router.GET("/projects", server.RequireAuth, server.getAllProjects)
	router.PUT("/projects/:id", server.RequireAuth, server.updateProject)
	router.DELETE("/projects/:id", server.RequireAuth, server.deleteProject)
	router.POST("/projects/member", server.RequireAuth, server.addProjectMember)
	router.DELETE("/projects/member", server.RequireAuth, server.removeProjectMember)

	//image
	router.POST("/images", server.RequireAuth, server.createImage)
	router.GET("/images/:url", server.getImage)

	//news
	router.POST("/news", server.RequireAuth, server.createNews)
	router.GET("/news", server.getAllNews)

	server.router = router

	return server
}

func (s *Server) Start(address string) error {

	s.query.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:            "admin",
		LastName:        "admin",
		Email:           "admin@admin.com",
		Password:        "$2a$10$3QYYykR1IWPX.KG9ne2mN..A/jZjynJd4qK.o0lRDxR/KAxBQTCXi",
		TelephoneNumber: "123123123",
		Role:            "admin",
		University:      "ytu",
		Department:      "mtm",
		DateOfBirth:     time.Now(),
	})

	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type Response struct {
	IsSuccess bool        `json:"isSuccess"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
}
