package api

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"time"
	"yildizskylab/src/db/sqlc"
	"yildizskylab/src/util"

	"github.com/gin-gonic/gin"
)

type createNewsRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (s *Server) createNews(c *gin.Context) {
	var req createNewsRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	file, err := c.FormFile("cover_image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   "Cover image is required",
		})
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}
	defer srcFile.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, srcFile); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "Unable to read cover image",
		})
		return
	}
	imageBytes := buf.Bytes()

	url := genereteUrl()

	userResult, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}
	user := userResult.(sqlc.User)

	savedImage, err := s.query.SaveImage(c, sqlc.SaveImageParams{
		Type:      file.Header.Get("Content-Type"),
		Name:      file.Filename,
		Data:      imageBytes,
		Url:       url,
		CreatedBy: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "Unable to save cover image",
		})
		return
	}

	config, err := util.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "Config load error",
		})
		return
	}
	savedImage.Url = config.Domain + "/images/" + savedImage.Url

	news, err := s.query.CreateNews(c, sqlc.CreateNewsParams{
		Title:        req.Title,
		Description:  req.Description,
		PublishDate:  time.Now(),
		CreatedByID:  sql.NullInt32{Int32: user.ID, Valid: true},
		CoverImageID: sql.NullInt32{Int32: savedImage.ID, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "Unable to save news",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "News created successfully",
		Data:      news,
	})
}

type getAllNewsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type NewsWithDetails struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	PublishDate time.Time `json:"publish_date"`
	Description string    `json:"description"`

	CoverImage struct {
		ID   int    `json:"id"`
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"cover_image"`

	CreatedBy struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		LastName   string `json:"last_name"`
		Email      string `json:"email"`
		University string `json:"university"`
		Department string `json:"department"`
	} `json:"created_by"`
}

func (s *Server) getAllNews(c *gin.Context) {

	var req getAllNewsRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	news, err := s.query.GetNewsWithDetails(c, sqlc.GetNewsWithDetailsParams{
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

	config, err := util.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "Config load error",
		})
		return
	}

	var newsList []NewsWithDetails
	for _, n := range news {
		newsList = append(newsList, NewsWithDetails{
			ID:          int(n.ID),
			Title:       n.Title,
			PublishDate: n.PublishDate,
			Description: n.Description,
			CoverImage: struct {
				ID   int    `json:"id"`
				URL  string `json:"url"`
				Type string `json:"type"`
			}{
				ID:   int(n.ImageID),
				URL:  config.Domain + "/images/" + n.ImageUrl,
				Type: n.ImageType,
			},
			CreatedBy: struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				LastName   string `json:"last_name"`
				Email      string `json:"email"`
				University string `json:"university"`
				Department string `json:"department"`
			}{
				ID:         int(n.UserID),
				Name:       n.UserName,
				LastName:   n.UserLastName,
				Email:      n.UserEmail,
				University: n.UserUniversity,
				Department: n.UserDepartment,
			},
		})
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "News got successfully",
		Data:      newsList,
	})
}

type getNewsRequest struct {
	ID int32 `uri:"id" binding:"required"`
}

func (s *Server) getNews(c *gin.Context) {

	var req getNewsRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	n, err := s.query.GetANewsWithDetails(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{
				IsSuccess: false,
				Message:   "News not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	config, err := util.LoadConfig(".")
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "Config load error",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "News got successfully",
		Data: NewsWithDetails{
			ID:          int(n.ID),
			Title:       n.Title,
			PublishDate: n.PublishDate,
			Description: n.Description,
			CoverImage: struct {
				ID   int    `json:"id"`
				URL  string `json:"url"`
				Type string `json:"type"`
			}{
				ID:   int(n.ImageID),
				URL:  config.Domain + "/images/" + n.ImageUrl,
				Type: n.ImageType,
			},
			CreatedBy: struct {
				ID         int    `json:"id"`
				Name       string `json:"name"`
				LastName   string `json:"last_name"`
				Email      string `json:"email"`
				University string `json:"university"`
				Department string `json:"department"`
			}{
				ID:         int(n.UserID),
				Name:       n.UserName,
				LastName:   n.UserLastName,
				Email:      n.UserEmail,
				University: n.UserUniversity,
				Department: n.UserDepartment,
			},
		},
	})
}
