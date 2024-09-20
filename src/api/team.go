package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
)

// CREATE TEAM
type createTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) createTeam(c *gin.Context) {
	var req createTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	team, err := s.query.CreateTeam(c, sqlc.CreateTeamParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team created successfully",
		Data:      team,
	})
}

///////////////////////////////

// GET TEAM
type getTeamRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type getTeamResponse struct {
	Id          int32                `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	TeamLeads   []returnUserResponse `json:"team_leads"`
	Projects    []sqlc.Project       `json:"projects"`
}

func (s *Server) getTeam(c *gin.Context) {
	var req getTeamRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	team, err := s.query.GetTeam(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{
				IsSuccess: false,
				Message:   "Team not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	leadIds, err := s.query.GetTeamLeadByTeamId(c, team.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	var leads []returnUserResponse

	for _, leadId := range leadIds {
		lead, err := s.query.GetUserWithNoDetails(c, leadId)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				IsSuccess: false,
				Message:   err.Error(),
			})
			return
		}

		leadToReturn := returnUserResponse{
			Id:              lead.ID,
			Name:            lead.Name,
			LastName:        lead.LastName,
			Email:           lead.Email,
			TelephoneNumber: lead.TelephoneNumber,
			University:      lead.University,
			Department:      lead.Department,
			Role:            lead.Role,
			DateOfBirth:     lead.DateOfBirth,
		}

		leads = append(leads, leadToReturn)
	}

	projectIds, err := s.query.GetTeamProjectByTeamId(c, team.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	var projects []sqlc.Project

	for _, projectId := range projectIds {
		project, err := s.query.GetProject(c, projectId)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				IsSuccess: false,
				Message:   err.Error(),
			})
		}
		projects = append(projects, project)
	}

	if ok := s.checkIfUserIsTeamLead(c, team.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team got successfully",
		Data: getTeamResponse{
			Id:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			TeamLeads:   leads,
			Projects:    projects,
		},
	})
}

///////////////////////////////

// GET ALL TEAMS
type getAllTeamsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) getAllTeams(c *gin.Context) {
	var req getAllTeamsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	teams, err := s.query.GetAllTeams(c, sqlc.GetAllTeamsParams{
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

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Teams got successfully",
		Data:      teams,
	})
}

///////////////////////////////

// UPDATE TEAM
type updateTeamRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) updateTeam(c *gin.Context) {

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

	var req updateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsTeamLead(c, id); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to update this team",
		})
		return
	}

	updatedTeam, err := s.query.UpdateTeam(c, sqlc.UpdateTeamParams{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team updated successfully",
		Data:      updatedTeam,
	})
}

///////////////////////////////

// DELETE TEAM
type deleteTeamRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteTeam(c *gin.Context) {
	var req deleteTeamRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsTeamLead(c, req.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to delete this team",
		})
		return
	}

	err := s.query.DeleteTeam(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team deleted successfully",
	})
}

// ADD TEAM PROJECT
type addTeamProjectRequest struct {
	TeamID    int32 `json:"team_id" binding:"required,min=1"`
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
}

func (s *Server) addTeamProject(c *gin.Context) {
	var req addTeamProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	_, err := s.query.GetTeam(c, req.TeamID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Team not found",
		})
		return
	}

	_, err = s.query.GetProject(c, req.ProjectID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Project not found",
		})
		return
	}

	if ok := s.checkIfUserIsTeamLead(c, req.TeamID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to add project to this team",
		})
		return
	}

	teamProject, err := s.query.CreateTeamProject(c, sqlc.CreateTeamProjectParams{
		TeamID:    req.TeamID,
		ProjectID: req.ProjectID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project added to team successfully",
		Data:      teamProject,
	})
}

///////////////////////////////

// REMOVE TEAM PROJECT
type removeTeamProjectRequest struct {
	TeamID    int32 `json:"team_id" binding:"required,min=1"`
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
}

func (s *Server) removeTeamProject(c *gin.Context) {
	var req removeTeamProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsTeamLead(c, req.TeamID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to remove project from this team",
		})
		return
	}

	err := s.query.DeleteTeamProject(c, sqlc.DeleteTeamProjectParams{
		TeamID:    req.TeamID,
		ProjectID: req.ProjectID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project removed from team successfully",
	})
}

///////////////////////////////

// ADD TEAM MEMBER
type addTeamMemberRequest struct {
	TeamID int32  `json:"team_id" binding:"required,min=1"`
	UserId int32  `json:"user_id" binding:"required,min=1"`
	Role   string `json:"role" binding:"required"`
}

func (s *Server) addTeamMember(c *gin.Context) {
	var req addTeamMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	_, err := s.query.GetTeam(c, req.TeamID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Team not found",
		})
		return
	}

	_, err = s.query.GetUser(c, req.UserId) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}

	if ok := s.checkIfUserIsTeamLead(c, req.TeamID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to add member to this team",
		})
		return
	}

	teamMember, err := s.query.CreateTeamMember(c, sqlc.CreateTeamMemberParams{
		TeamID: req.TeamID,
		UserID: req.UserId,
		Role:   req.Role,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team member added successfully",
		Data:      teamMember,
	})
}

///////////////////////////////

// REMOVE TEAM MEMBER
type removeTeamMemberRequest struct {
	TeamID int32 `json:"team_id" binding:"required,min=1"`
	UserID int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeTeamMember(c *gin.Context) {
	var req removeTeamMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsTeamLead(c, req.TeamID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to remove member from this team",
		})
		return
	}

	err := s.query.DeleteTeamMember(c, sqlc.DeleteTeamMemberParams{
		TeamID: req.TeamID,
		UserID: req.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team member removed successfully",
	})
}

///////////////////////////////

// UTILS
// checkIfUserIsTeamLead checks if the user is team lead
func (s *Server) checkIfUserIsTeamLead(c *gin.Context, teamID int32) bool {

	anyUser, ok := c.Get("user")
	if !ok {
		return false
	}

	user := anyUser.(sqlc.User)

	if user.Role == "admin" {
		return true
	}

	member, err := s.query.GetTeamMember(c, sqlc.GetTeamMemberParams{
		TeamID: teamID,
		UserID: user.ID,
	})

	if err != nil {
		return false
	}

	if member.Role == "lead" {
		return true
	}

	return false
}
