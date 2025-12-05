package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/barzaevhalid/cloud_storage_backend/services"
	"github.com/gofiber/fiber/v2"
)

type FileHandler struct {
	FileService *services.FileService
}

func NewFileHandler(s *services.FileService) *FileHandler {
	return &FileHandler{FileService: s}
}

// @Router /api/files/upload [post]
// @Security ApiKeyAuth
func (h *FileHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file required"})
	}

	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File too large. Maximum 5 MB",
		})
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	// Save to DB
	err = c.SaveFile(file, "./uploads/"+filename)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	id, err := h.FileService.SaveFileMetadata(1, file, filename)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	fileUrl := fmt.Sprintf("/uploads/%s", filename)

	return c.JSON(fiber.Map{
		"id":           id,
		"filename":     filename,
		"originalName": file.Filename,
		"mimeType":     file.Header.Get("Content-Type"),
		"size":         file.Size,
		"url":          fileUrl,
	})
}

func (h *FileHandler) GetFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	return c.SendFile("./uploads/" + filename)
}

func (h *FileHandler) FindAllFiles(c *fiber.Ctx) error {
	fileType := c.Query("type", "all")
	userId := c.Locals("user_id").(int64)
	files, err := h.FileService.FindAllFiles(userId, fileType)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(files)
}
func (h *FileHandler) DeleteFiles(c *fiber.Ctx) error {
	idsParam := c.Params("ids", "")

	if idsParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ids required"})
	}
	idStrings := strings.Split(idsParam, ",")

	ids := make([]int64, 0, len(idStrings))

	for _, s := range idStrings {
		id, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		ids = append(ids, id)
	}
	userId := c.Locals("user_id").(int64)

	err := h.FileService.MarkDeleted(userId, ids)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "deleted"})
}
