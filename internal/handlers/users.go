package handlers

import (
	"database/sql"
	"fmt"
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

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type UpdateUserRequest struct {
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Phone      string   `json:"phone"`
	Position   string   `json:"position"`
	Department string   `json:"department"`
	Location   string   `json:"location"`
	Birthday   string   `json:"birthday"`
	Skype      string   `json:"skype"`
	Skills     []string `json:"skills"`
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	search := c.Query("search")
	role := c.Query("role")
	department := c.Query("department")

	// Build query
	query := `
		SELECT id, email, first_name, last_name, avatar_url, position, department, role, status, created_at
		FROM users WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query += fmt.Sprintf(" AND (LOWER(first_name) LIKE $%d OR LOWER(last_name) LIKE $%d OR LOWER(email) LIKE $%d)", argIndex, argIndex, argIndex)
		countQuery += fmt.Sprintf(" AND (LOWER(first_name) LIKE $%d OR LOWER(last_name) LIKE $%d OR LOWER(email) LIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, searchTerm)
		argIndex++
	}

	if role != "" {
		query += fmt.Sprintf(" AND role = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND role = $%d", argIndex)
		args = append(args, role)
		argIndex++
	}

	if department != "" {
		query += fmt.Sprintf(" AND department = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND department = $%d", argIndex)
		args = append(args, department)
		argIndex++
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get users count")
		return
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get users")
		return
	}
	defer rows.Close()

	users := []models.UserListItem{}
	for rows.Next() {
		var user models.UserListItem
		var avatar, position, department sql.NullString

		err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &avatar, &position, &department, &user.Role, &user.Status, &user.CreatedAt)
		if err != nil {
			continue
		}

		if avatar.Valid {
			user.Avatar = avatar.String
		}
		if position.Valid {
			user.Position = position.String
		}
		if department.Valid {
			user.Department = department.String
		}

		users = append(users, user)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": users,
		"meta": models.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: utils.CalculateTotalPages(total, limit),
		},
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.getUserByID(userID)
	if err == sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.getUserByID(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)
	currentUserRole := middleware.GetUserRole(c)

	// Check if user can update this profile
	if userID != currentUserID && currentUserRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "You can only update your own profile")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse birthday if provided
	var birthday *time.Time
	if req.Birthday != "" {
		parsed, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			utils.RespondWithError(c, http.StatusBadRequest, "Invalid birthday format (use YYYY-MM-DD)")
			return
		}
		birthday = &parsed
	}

	// Update user
	_, err := database.DB.Exec(`
		UPDATE users SET
			first_name = COALESCE(NULLIF($1, ''), first_name),
			last_name = COALESCE(NULLIF($2, ''), last_name),
			phone = $3,
			position = $4,
			department = $5,
			location = $6,
			birthday = $7,
			skype = $8,
			updated_at = NOW()
		WHERE id = $9
	`, req.FirstName, req.LastName, nullString(req.Phone), nullString(req.Position),
		nullString(req.Department), nullString(req.Location), birthday, nullString(req.Skype), userID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Update skills
	if req.Skills != nil {
		// Delete existing skills
		database.DB.Exec("DELETE FROM user_skills WHERE user_id = $1", userID)

		// Insert new skills
		for _, skill := range req.Skills {
			if skill != "" {
				database.DB.Exec("INSERT INTO user_skills (user_id, skill) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, skill)
			}
		}
	}

	// Get updated user
	user, err := h.getUserByID(userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get updated user")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, user)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := middleware.GetUserID(c)
	currentUserRole := middleware.GetUserRole(c)

	// Check if user can update this avatar
	if userID != currentUserID && currentUserRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "You can only update your own avatar")
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		utils.RespondWithError(c, http.StatusBadRequest, "Only image files are allowed")
		return
	}

	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create upload directory")
		return
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	out, err := os.Create(filepath)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Update user avatar URL
	avatarURL := "/uploads/avatars/" + filename
	_, err = database.DB.Exec("UPDATE users SET avatar_url = $1, updated_at = NOW() WHERE id = $2", avatarURL, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update avatar")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"avatar": avatarURL})
}

func (h *UserHandler) getUserByID(userID string) (*models.UserResponse, error) {
	var user models.User
	var avatar, phone, position, department, level, location, skype sql.NullString
	var birthday sql.NullTime

	err := database.DB.QueryRow(`
		SELECT id, email, first_name, last_name, avatar_url, phone, position, department, level, location, birthday, skype, role, status, created_at, updated_at
		FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &avatar, &phone, &position, &department, &level, &location, &birthday, &skype, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	response := &models.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if avatar.Valid {
		response.Avatar = avatar.String
	}
	if phone.Valid {
		response.Phone = phone.String
	}
	if position.Valid {
		response.Position = position.String
	}
	if department.Valid {
		response.Department = department.String
	}
	if level.Valid {
		response.Level = level.String
	}
	if location.Valid {
		response.Location = location.String
	}
	if birthday.Valid {
		response.Birthday = birthday.Time.Format("2006-01-02")
	}
	if skype.Valid {
		response.Skype = skype.String
	}

	// Get skills
	rows, err := database.DB.Query("SELECT skill FROM user_skills WHERE user_id = $1", userID)
	if err == nil {
		defer rows.Close()
		skills := []string{}
		for rows.Next() {
			var skill string
			if rows.Scan(&skill) == nil {
				skills = append(skills, skill)
			}
		}
		response.Skills = skills
	}

	return response, nil
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
