package api

import (
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
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

	//team
	router.POST("/teams", server.RequireAuth, server.createTeam)
	router.GET("/teams/:id", server.RequireAuth, server.getTeam)
	router.GET("/teams", server.RequireAuth, server.getAllTeams)
	router.PUT("/teams/:id", server.RequireAuth, server.updateTeam)
	router.DELETE("/teams/:id", server.RequireAuth, server.deleteTeam)
	router.POST("/teams/lead", server.RequireAuth, server.addTeamLead)
	router.DELETE("/teams/lead", server.RequireAuth, server.removeTeamLead)
	router.POST("/teams/project", server.RequireAuth, server.addTeamProject)
	router.DELETE("/teams/project", server.RequireAuth, server.removeTeamProject)
	router.POST("/teams/member", server.RequireAuth, server.addTeamMember)
	router.DELETE("/teams/member", server.RequireAuth, server.removeTeamMember)

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
	router.POST("/projects/lead", server.RequireAuth, server.addProjectLead)
	router.DELETE("/projects/lead", server.RequireAuth, server.removeProjectLead)
	router.POST("/projects/member", server.RequireAuth, server.addProjectMember)
	router.DELETE("/projects/member", server.RequireAuth, server.removeProjectMember)

	server.router = router

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
