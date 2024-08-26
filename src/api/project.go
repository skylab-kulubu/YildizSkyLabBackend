package api

import (
	"database/sql"
	"net/http"
	"yildizskylab/src/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) createProject(c *gin.Context) {

	var req createProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	project, err := s.query.CreateProject(c, sqlc.CreateProjectParams{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, project)
}

type getProjectRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getProject(c *gin.Context) {
	var req getProjectRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	project, err := s.query.GetProject(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, project)
}

type getAllProjectsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) getAllProjects(c *gin.Context) {
	var req getAllProjectsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.GetAllProjectsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	projects, err := s.query.GetAllProjects(c, arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, projects)
}

type updateProjectRequest struct {
	ID          int32  `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) updateProject(c *gin.Context) {
	var req updateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updatedProject, err := s.query.UpdateProject(c, sqlc.UpdateProjectParams{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
	}

	c.JSON(http.StatusOK, updatedProject)
}

type deleteProjectRequest struct {
	ID int32 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteProject(c *gin.Context) {
	var req deleteProjectRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteProject(c, req.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type addProjectLeadRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) addProjectLead(c *gin.Context) {
	var req addProjectLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.query.GetProject(c, req.ProjectID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.query.GetUser(c, req.UserID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	projectLead, err := s.query.CreateProjectLead(c, sqlc.CreateProjectLeadParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, projectLead)
}

type removeProjectLeadRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeProjectLead(c *gin.Context) {
	var req removeProjectLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteProjectLead(c, sqlc.DeleteProjectLeadParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type addProjectMemberRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) addProjectMember(c *gin.Context) {
	var req addProjectMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.query.GetProject(c, req.ProjectID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = s.query.GetUser(c, req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	projectMember, err := s.query.CreateProjectMember(c, sqlc.CreateProjectMemberParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, projectMember)
}

type removeProjectMemberRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeProjectMember(c *gin.Context) {
	var req removeProjectMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.query.DeleteProjectMember(c, sqlc.DeleteProjectMemberParams{
		ProjectID: req.ProjectID,
		UserID:    req.UserID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
