package api

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"yildizskylab/src/db/sqlc"
	"yildizskylab/src/util"

	"github.com/gin-gonic/gin"
)

func (s *Server) createImage(c *gin.Context) {

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
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

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, srcFile); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}
	imageBytes := buf.Bytes()

	userResult, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   "User not found",
		})
		return
	}
	user := userResult.(sqlc.User)

	url := genereteUrl()

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
			Message:   err.Error(),
		})
		return
	}

	config, err := util.LoadConfig(".")
	savedImage.Url = config.Domain + "/images/" + savedImage.Url

	c.JSON(http.StatusOK, Response{
		IsSuccess: true,
		Message:   "Image uploaded successfully",
		Data:      savedImage,
	})

}

type getImageRequest struct {
	ImageUrl string `uri:"url" binding:"required"`
}

func (s *Server) getImage(c *gin.Context) {
	var req getImageRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	image, err := s.query.GetImageByUrl(c, req.ImageUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			IsSuccess: false,
			Message:   err.Error(),
		})
		return
	}

	// Resmin Content-Type'ını belirtelim
	c.Header("Content-Type", image.Type)

	// Resmi byte olarak döndürelim
	c.Writer.Write(image.Data)

}

// UTILS
func genereteUrl() string {
	randomBytes := make([]byte, 64)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes)
}
