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
	Leads       []returnUserResponse `json:"leads"`
	Members     []returnUserResponse `json:"members"`
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

	teamsWithDetails, err := s.query.GetTeamWithDetails(c, req.ID)

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

	if ok := s.checkIfUserIsTeamLead(c, teamsWithDetails[0].ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	var team getTeamResponse

	team.Id = teamsWithDetails[0].ID
	team.Name = teamsWithDetails[0].Name
	team.Description = teamsWithDetails[0].Description

	leads := []returnUserResponse{}

	if !teamsWithDetails[0].LeadID.Valid {
		team.Leads = nil
	} else {
		for _, tl := range teamsWithDetails {
			lead := returnUserResponse{
				Id:              tl.LeadID.Int32,
				Name:            tl.LeadName.String,
				LastName:        tl.LeadLastName.String,
				Email:           tl.LeadEmail.String,
				TelephoneNumber: tl.LeadTelephoneNumber.String,
				University:      tl.LeadUniversity.String,
				Department:      tl.LeadDepartment.String,
				DateOfBirth:     tl.LeadDateOfBirth.Time,
			}
			leads = append(leads, lead)
		}

		seen := make(map[int32]bool)

		for _, lead := range leads {
			if !seen[lead.Id] {
				team.Leads = append(team.Leads, lead)
				seen[lead.Id] = true
			}
		}
	}

	members := []returnUserResponse{}

	if !teamsWithDetails[0].MemberID.Valid {
		team.Members = nil
	} else {
		for _, tm := range teamsWithDetails {
			member := returnUserResponse{
				Id:              tm.MemberID.Int32,
				Name:            tm.MemberName.String,
				LastName:        tm.MemberLastName.String,
				Email:           tm.MemberEmail.String,
				TelephoneNumber: tm.MemberTelephoneNumber.String,
				University:      tm.MemberUniversity.String,
				Department:      tm.MemberDepartment.String,
				DateOfBirth:     tm.MemberDateOfBirth.Time,
			}
			members = append(members, member)
		}

		seen := make(map[int32]bool)

		for _, member := range members {
			if !seen[member.Id] {
				team.Members = append(team.Members, member)
				seen[member.Id] = true
			}
		}
	}
	projects := []sqlc.Project{}

	if !teamsWithDetails[0].ProjectID.Valid {
		team.Projects = nil
	} else {
		for _, tp := range teamsWithDetails {
			project := sqlc.Project{
				ID:          tp.ProjectID.Int32,
				Name:        tp.ProjectName.String,
				Description: tp.ProjectDescription.String,
			}
			projects = append(projects, project)
		}

		seen := make(map[int32]bool)

		for _, project := range projects {
			if !seen[project.ID] {
				team.Projects = append(team.Projects, project)
				seen[project.ID] = true
			}
		}
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Team got successfully",
		Data: getTeamResponse{
			Id:          team.Id,
			Name:        team.Name,
			Description: team.Description,
			Leads:       team.Leads,
			Members:     team.Members,
			Projects:    team.Projects,
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
	Name        *string `json:"name"`
	Description *string `json:"description" `
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

	updatedTeam, err := s.query.GetTeamWithDetails(c, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	if req.Name != nil {
		updatedTeam[0].Name = *req.Name

	}

	if req.Description != nil {
		updatedTeam[0].Description = *req.Description
	}

	team, err := s.query.UpdateTeam(c, sqlc.UpdateTeamParams{
		ID:          id,
		Name:        updatedTeam[0].Name,
		Description: updatedTeam[0].Description,
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
		Data:      team,
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

	_, err := s.query.GetTeamWithDetails(c, req.TeamID)

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
	TeamID int32 `json:"team_id" binding:"required,min=1"`
	UserId int32 `json:"user_id" binding:"required,min=1"`
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

	_, err := s.query.GetTeamWithDetails(c, req.TeamID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Team not found",
		})
		return
	}

	_, err = s.query.GetUserWithDetails(c, req.UserId) // need a better solition but work for now

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
		Role:   "member",
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

func (s *Server) addTeamLead(c *gin.Context) {
	var req addTeamMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	_, err := s.query.GetTeamWithDetails(c, req.TeamID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Team not found",
		})
		return
	}

	_, err = s.query.GetUserWithDetails(c, req.UserId) // need a better solition but work for now

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
		Role:   "lead",
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
