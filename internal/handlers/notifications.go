package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

type RegisterDeviceRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required"`
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	readFilter := c.Query("read")

	// Build query
	query := `
		SELECT id, type, title, message, data, read, created_at
		FROM notifications WHERE user_id = $1
	`
	countQuery := "SELECT COUNT(*) FROM notifications WHERE user_id = $1"
	unreadQuery := "SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false"
	args := []interface{}{userID}
	argIndex := 2

	if readFilter != "" {
		isRead := readFilter == "true"
		query += fmt.Sprintf(" AND read = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND read = $%d", argIndex)
		args = append(args, isRead)
		argIndex++
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get notifications count")
		return
	}

	// Get unread count
	var unreadCount int
	database.DB.QueryRow(unreadQuery, userID).Scan(&unreadCount)

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get notifications")
		return
	}
	defer rows.Close()

	notifications := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		var message sql.NullString
		var data []byte

		err := rows.Scan(&n.ID, &n.Type, &n.Title, &message, &data, &n.Read, &n.CreatedAt)
		if err != nil {
			continue
		}

		n.UserID = userID
		if message.Valid {
			n.Message = message.String
		}
		if len(data) > 0 {
			json.Unmarshal(data, &n.Data)
		}

		notifications = append(notifications, n)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": notifications,
		"meta": models.NotificationMeta{
			Total:       total,
			UnreadCount: unreadCount,
			Page:        page,
			Limit:       limit,
		},
	})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")
	userID := middleware.GetUserID(c)

	result, err := database.DB.Exec(`
		UPDATE notifications SET read = true WHERE id = $1 AND user_id = $2
	`, notificationID, userID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to mark notification as read")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.RespondWithError(c, http.StatusNotFound, "Notification not found")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)

	_, err := database.DB.Exec(`
		UPDATE notifications SET read = true WHERE user_id = $1 AND read = false
	`, userID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to mark all notifications as read")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationID := c.Param("id")
	userID := middleware.GetUserID(c)

	result, err := database.DB.Exec("DELETE FROM notifications WHERE id = $1 AND user_id = $2", notificationID, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		utils.RespondWithError(c, http.StatusNotFound, "Notification not found")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Notification deleted"})
}

func (h *NotificationHandler) RegisterDevice(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req RegisterDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate platform
	validPlatforms := []string{"web", "ios", "android"}
	isValid := false
	for _, p := range validPlatforms {
		if req.Platform == p {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.RespondWithError(c, http.StatusBadRequest, "Platform must be 'web', 'ios', or 'android'")
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO notification_devices (user_id, token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, token) DO UPDATE SET platform = $3
	`, userID, req.Token, req.Platform)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to register device")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Device registered successfully"})
}
