package handlers

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file size (max 50MB)
	if header.Size > 50*1024*1024 {
		utils.RespondWithError(c, http.StatusBadRequest, "File too large (max 50MB)")
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads/files"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	out, err := os.Create(filePath)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Save file record to database
	fileURL := "/uploads/files/" + filename
	var fileID string
	err = database.DB.QueryRow(`
		INSERT INTO files (name, url, type, size, uploaded_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, header.Filename, fileURL, contentType, header.Size, userID).Scan(&fileID)

	if err != nil {
		// Clean up file on database error
		os.Remove(filePath)
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save file record")
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, models.File{
		ID:         fileID,
		Name:       header.Filename,
		URL:        fileURL,
		Type:       contentType,
		Size:       int(header.Size),
		UploadedAt: time.Now(),
	})
}

func (h *FileHandler) GetFile(c *gin.Context) {
	fileID := c.Param("id")

	var file models.File
	var fileType sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, name, url, type, size, created_at
		FROM files WHERE id = $1
	`, fileID).Scan(&file.ID, &file.Name, &file.URL, &fileType, &file.Size, &file.UploadedAt)

	if err == sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusNotFound, "File not found")
		return
	}
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get file")
		return
	}

	if fileType.Valid {
		file.Type = fileType.String
	}

	utils.RespondWithSuccess(c, http.StatusOK, file)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	// Get file details
	var file struct {
		URL        string
		UploadedBy sql.NullString
	}

	err := database.DB.QueryRow("SELECT url, uploaded_by FROM files WHERE id = $1", fileID).Scan(&file.URL, &file.UploadedBy)
	if err == sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusNotFound, "File not found")
		return
	}
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get file")
		return
	}

	// Check permissions
	if file.UploadedBy.Valid && file.UploadedBy.String != userID && userRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "You can only delete your own files")
		return
	}

	// Delete from database
	_, err = database.DB.Exec("DELETE FROM files WHERE id = $1", fileID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	// Delete physical file
	if strings.HasPrefix(file.URL, "/uploads/") {
		os.Remove("." + file.URL)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "File deleted successfully"})
}
