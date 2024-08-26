package api

import (
	"database/sql"
	"net/http"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) createTeam(c *gin.Context) {
	var req createTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	team, err := s.query.CreateTeam(c, sqlc.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, team)
}

type getTeamRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getTeam(c *gin.Context) {
	var req getTeamRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	team, err := s.query.GetTeam(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, team)
}

type getAllTeamsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) getAllTeams(c *gin.Context) {
	var req getAllTeamsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.GetAllTeamsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	teams, err := s.query.GetAllTeams(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, teams)
}

type updateTeamRequest struct {
	ID          int32  `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) updateTeam(c *gin.Context) {
	var req updateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedTeam, err := s.query.UpdateTeam(c, sqlc.UpdateTeamParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	c.JSON(http.StatusOK, updatedTeam)
}

type deleteTeamRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteTeam(c *gin.Context) {
	var req deleteTeamRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteTeam(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type addTeamLeadRequest struct {
	TeamID int32 `json:"team_id" binding:"required,min=1"`
	UserId int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) addTeamLead(c *gin.Context) {
	var req addTeamLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.query.GetTeam(c, req.TeamID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.query.GetUser(c, req.UserId) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	teamLead, err := s.query.CreateTeamLead(c, sqlc.CreateTeamLeadParams{
		TeamID: req.TeamID,
		UserID: req.UserId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, teamLead)
}

type removeTeamLeadRequest struct {
	TeamID int32 `json:"team_id" binding:"required,min=1"`
	UserId int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeTeamLead(c *gin.Context) {
	var req removeTeamLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteTeamLead(c, sqlc.DeleteTeamLeadParams{
		TeamID: req.TeamID,
		UserID: req.UserId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type addTeamProjectRequest struct {
	TeamID    int32 `json:"team_id" binding:"required,min=1"`
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
}

func (s *Server) addTeamProject(c *gin.Context) {
	var req addTeamProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.query.GetTeam(c, req.TeamID)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.query.GetProject(c, req.ProjectID)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	teamProject, err := s.query.CreateTeamProject(c, sqlc.CreateTeamProjectParams{
		TeamID:    req.TeamID,
		ProjectID: req.ProjectID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, teamProject)
}

type removeTeamProjectRequest struct {
	TeamID    int32 `json:"team_id" binding:"required,min=1"`
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
}

func (s *Server) removeTeamProject(c *gin.Context) {
	var req removeTeamProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteTeamProject(c, sqlc.DeleteTeamProjectParams{
		TeamID:    req.TeamID,
		ProjectID: req.ProjectID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type addTeamMemberRequest struct {
	TeamID int32 `json:"team_id" binding:"required,min=1"`
	UserId int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) addTeamMember(c *gin.Context) {
	var req addTeamMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.query.GetTeam(c, req.TeamID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.query.GetUser(c, req.UserId) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	teamMember, err := s.query.CreateTeamMember(c, sqlc.CreateTeamMemberParams{
		TeamID: req.TeamID,
		UserID: req.UserId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, teamMember)
}

type removeTeamMemberRequest struct {
	TeamID int32 `json:"team_id" binding:"required,min=1"`
	UserID int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeTeamMember(c *gin.Context) {
	var req removeTeamMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteTeamMember(c, sqlc.DeleteTeamMemberParams{
		TeamID: req.TeamID,
		UserID: req.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
