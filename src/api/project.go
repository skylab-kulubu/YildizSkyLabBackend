package api

import (
	"database/sql"
	"net/http"
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

	if ok := s.checkIfUserIsProjectLead(c, project.ProjectID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Project got successfully",
		Data:      project,
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

	arg := sqlc.GetAllProjectsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	projects, err := s.query.GetAllProjects(c, arg)

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
	ID          int32  `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func (s *Server) updateProject(c *gin.Context) {
	var req updateProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
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

	updatedProject, err := s.query.UpdateProject(c, sqlc.UpdateProjectParams{
		ID:          req.ID,
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
		Message:   "Project updated successfully",
		Data:      updatedProject,
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

// ADD PROJECT LEAD
type addProjectLeadRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) addProjectLead(c *gin.Context) {
	var req addProjectLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	_, err := s.query.GetProject(c, req.ProjectID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Project not found",
		})
		return
	}

	_, err = s.query.GetUser(c, req.UserID) // need a better solition but work for now

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}

	projectLead, err := s.query.CreateProjectLead(c, sqlc.CreateProjectLeadParams{
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
		Message:   "Project lead added successfully",
		Data:      projectLead,
	})
}

////////////////////////

// REMOVE PROJECT LEAD
type removeProjectLeadRequest struct {
	ProjectID int32 `json:"project_id" binding:"required,min=1"`
	UserID    int32 `json:"user_id" binding:"required,min=1"`
}

func (s *Server) removeProjectLead(c *gin.Context) {
	var req removeProjectLeadRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	err := s.query.DeleteProjectLead(c, sqlc.DeleteProjectLeadParams{
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
		Message:   "Project lead removed successfully",
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

	_, err = s.query.GetUser(c, req.UserID)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}

	if ok := s.checkIfUserIsProjectLead(c, project.ProjectID); !ok {
		c.JSON(http.StatusForbidden, Response{
			IsSuccess: false,
			Message:   "You are not authorized to see this team",
		})
		return
	}

	projectMember, err := s.query.CreateProjectMember(c, sqlc.CreateProjectMemberParams{
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
		Message:   "Project member added successfully",
		Data:      projectMember,
	})
}

////////////////////////

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

	_, err := s.query.GetProjectLead(c, sqlc.GetProjectLeadParams{
		ProjectID: projectID,
		UserID:    user.ID,
	})

	if err != nil {
		return false
	}

	return true

}
