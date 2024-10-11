package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
)

// CREATE PROJECT
type createProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) createProject(c *gin.Context) {

	var req createProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	project, err := s.query.CreateProject(c, sqlc.CreateProjectParams{
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
		Message:   "Project created successfully",
		Data:      project,
	})
}

////////////////////////

// GET PROJECT
type getProjectRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

type getProjectResponse struct {
	Id          int32                `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Leads       []returnUserResponse `json:"leads"`
	Members     []returnUserResponse `json:"members"`
	Teams       []sqlc.Team          `json:"teams"`
}

func (s *Server) getProject(c *gin.Context) {
	var req getProjectRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	projectWithDetails, err := s.query.GetProjectWithDetails(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{
				IsSuccess: false,
				Message:   "Project not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, projectWithDetails[0].ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	var project getProjectResponse

	project.Id = projectWithDetails[0].ID
	project.Name = projectWithDetails[0].Name
	project.Description = projectWithDetails[0].Description

	leads := []returnUserResponse{}

	if !projectWithDetails[0].LeadID.Valid {
		project.Leads = nil
	} else {
		for _, pl := range projectWithDetails {
			lead := returnUserResponse{
				Id:              pl.LeadID.Int32,
				Name:            pl.LeadName.String,
				LastName:        pl.LeadLastName.String,
				Email:           pl.LeadEmail.String,
				TelephoneNumber: pl.LeadTelephoneNumber.String,
				University:      pl.LeadUniversity.String,
				Department:      pl.LeadDepartment.String,
				DateOfBirth:     pl.LeadDateOfBirth.Time,
			}
			leads = append(leads, lead)
		}

		seen := make(map[int32]bool)

		for _, lead := range leads {
			if !seen[lead.Id] {
				project.Leads = append(project.Leads, lead)
				seen[lead.Id] = true
			}
		}
	}

	members := []returnUserResponse{}

	if !projectWithDetails[0].MemberID.Valid {
		project.Members = nil
	} else {
		for _, pm := range projectWithDetails {
			member := returnUserResponse{
				Id:              pm.LeadID.Int32,
				Name:            pm.LeadName.String,
				LastName:        pm.LeadLastName.String,
				Email:           pm.LeadEmail.String,
				TelephoneNumber: pm.LeadTelephoneNumber.String,
				University:      pm.LeadUniversity.String,
				Department:      pm.LeadDepartment.String,
				DateOfBirth:     pm.LeadDateOfBirth.Time,
			}
			members = append(members, member)
		}

		seen := make(map[int32]bool)

		for _, member := range members {
			if !seen[member.Id] {
				project.Members = append(project.Members, member)
				seen[member.Id] = true
			}
		}
	}

	teams := []sqlc.Team{}

	if !projectWithDetails[0].TeamID.Valid {
		project.Teams = nil
	} else {
		for _, pt := range projectWithDetails {
			team := sqlc.Team{
				ID:          pt.TeamID.Int32,
				Name:        pt.TeamName.String,
				Description: pt.TeamDescription.String,
			}

			teams = append(teams, team)
		}

		seen := make(map[int32]bool)

		for _, team := range teams {
			if !seen[team.ID] {
				project.Teams = append(project.Teams, team)
				seen[team.ID] = true
			}
		}
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project got successfully",
		Data: getProjectResponse{
			Id:          project.Id,
			Name:        project.Name,
			Description: project.Description,
			Leads:       project.Leads,
			Members:     project.Members,
			Teams:       project.Teams,
		},
	})
}

////////////////////////

// GET ALL PROJECTS
type getAllProjectsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) getAllProjects(c *gin.Context) {
	var req getAllProjectsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	projects, err := s.query.GetAllProjects(c, sqlc.GetAllProjectsParams{
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
		Message:   "Projects got successfully",
		Data:      projects,
	})
}

////////////////////////

// UPDATE PROJECT
type updateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description" `
}

func (s *Server) updateProject(c *gin.Context) {

	var id int32

	idParam := c.Param("id")

	i, err := strconv.ParseInt(idParam, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	id = int32(i)

	var req updateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, id); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	updatedProject, err := s.query.GetProject(c, id)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	if req.Name != nil {
		updatedProject.Name = *req.Name

	}

	if req.Description != nil {
		updatedProject.Description = *req.Description
	}

	project, err := s.query.UpdateProject(c, sqlc.UpdateProjectParams{
		ID:          id,
		Name:        updatedProject.Name,
		Description: updatedProject.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project updated successfully",
		Data:      project,
	})
}

////////////////////////

// DELETE PROJECT
type deleteProjectRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteProject(c *gin.Context) {
	var req deleteProjectRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, req.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	err := s.query.DeleteProject(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project deleted successfully",
	})
}

////////////////////////

// ADD PROJECT MEMBER
type addProjectMemberRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) addProjectMember(c *gin.Context) {
	var req addProjectMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	project, err := s.query.GetProject(c, req.ProjectID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Project not found",
		})
		return
	}

	_, err = s.query.GetUserWithDetails(c, req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, project.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	projectMember, err := s.query.CreateProjectMember(c, sqlc.CreateProjectMemberParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		Role:      "member",
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
		Message:   "Project member added successfully",
		Data:      projectMember,
	})
}

////////////////////////

func (s *Server) addProjectLead(c *gin.Context) {
	var req addProjectMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	project, err := s.query.GetProject(c, req.ProjectID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Project not found",
		})
		return
	}

	_, err = s.query.GetUserWithDetails(c, req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, project.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	projectMember, err := s.query.CreateProjectMember(c, sqlc.CreateProjectMemberParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
		Role:      "lead",
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
		Message:   "Project member added successfully",
		Data:      projectMember,
	})
}

// REMOVE PROJECT MEMBER
type removeProjectMemberRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeProjectMember(c *gin.Context) {
	var req removeProjectMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, req.ProjectID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	err := s.query.DeleteProjectMember(c, sqlc.DeleteProjectMemberParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
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
		Message:   "Project member removed successfully",
	})
}

// //////////////////////
// checkIfUserIsTeamLead checks if the user is team lead
func (s *Server) checkIfUserIsProjectLead(c *gin.Context, projectID int32) bool {

	anyUser, ok := c.Get("user")
	if !ok {
		return false
	}

	user := anyUser.(sqlc.User)

	if user.Role == "admin" {
		return true
	}

	member, err := s.query.GetProjectMember(c, sqlc.GetProjectMemberParams{
		ProjectID: projectID,
		UserID:    user.ID,
	})

	if err != nil {
		return false
	}

	if member.Role == "lead" {
		return true
	}

	return false
}
