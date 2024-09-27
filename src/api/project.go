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

	project, err := s.query.GetProject(c, req.ID)

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

	if ok := s.checkIfUserIsProjectLead(c, project.ID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}
	var leads []returnUserResponse
	err = json.Unmarshal(project.Leads, &leads)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	var members []returnUserResponse
	err = json.Unmarshal(project.Members, &members)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}

	var teams []sqlc.Team
	err = json.Unmarshal(project.Teams, &teams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
	}
	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project got successfully",
		Data: getProjectResponse{
			Id:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			Leads:       leads,
			Members:     members,
			Teams:       teams,
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
