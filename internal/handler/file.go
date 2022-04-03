package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"restaurant-assistant/internal/domain"
)

const (
	maxUploadSize = 5 << 20 // 5 megabytes
)

var (
	imageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}
)

type uploadResponse struct {
	URL string `json:"url"`
}



// UploadImage @Summary upload Image
// @Tags image
// @Description Endpoint to upload and update images
// @ModuleID UploadImage
// @Accept mpfd
// @Produce json
// @Param file formData file true "file"
// @Param id path string true "dish, restaurant id"
// @Success 200 {object} uploadResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurant/image/{id} [post]
// @Router /dish/image/{id} [post]
func (h *Handler) UploadImage(c *gin.Context) {
	id := c.Param("id")
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Error().Err(err)
		newResponse(c, http.StatusBadRequest, "invalid image name in multipart/form")

		return
	}

	defer file.Close()

	buffer := make([]byte, fileHeader.Size)

	if _, err := file.Read(buffer); err != nil {
		log.Error().Err(err)
		newResponse(c, http.StatusBadRequest, "cannot read the file")

		return
	}

	contentType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := imageTypes[contentType]; !ex {
		newResponse(c, http.StatusBadRequest, "file type is not supported")

		return
	}

	tempFilename := fmt.Sprintf("%s--%s", id, fileHeader.Filename)

	f, err := os.OpenFile(tempFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o666)
	if err != nil {
		log.Error().Err(err)
		newResponse(c, http.StatusInternalServerError, "failed to create temp file")

		return
	}

	if _, err := io.Copy(f, bytes.NewReader(buffer)); err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to write chunk to temp file")

		return
	}

	if err := f.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close temp file")
	}

	path := c.FullPath()

	url, err := h.services.File.UploadAndSaveFile(c.Request.Context(), domain.File{
		ContentType: contentType,
		Name:        tempFilename,
		Size:        fileHeader.Size,
	},
		id,
		path)

	if err != nil {
		newResponse(c, http.StatusInternalServerError, "failed to upload image")

		return
	}

	c.JSON(http.StatusOK, &uploadResponse{url})
}
