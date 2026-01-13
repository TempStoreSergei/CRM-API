package handlers

import (
	"net/http"

	"github.com/TempStoreSergei/CRM-API/internal/config"
	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type IntegrationHandler struct {
	cfg *config.Config
}

func NewIntegrationHandler(cfg *config.Config) *IntegrationHandler {
	return &IntegrationHandler{cfg: cfg}
}

type GoogleCallbackRequest struct {
	Code string `json:"code" binding:"required"`
}

type GoogleSyncRequest struct {
	Direction string `json:"direction" binding:"required"`
}

type TelegramLinkRequest struct {
	TelegramUserID string `json:"telegramUserId" binding:"required"`
}

func (h *IntegrationHandler) GetGoogleAuthURL(c *gin.Context) {
	if h.cfg.GoogleClientID == "" {
		utils.RespondWithError(c, http.StatusServiceUnavailable, "Google integration not configured")
		return
	}

	// Build OAuth URL
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=" + h.cfg.GoogleClientID +
		"&redirect_uri=" + "http://localhost:8080/api/integrations/google/callback" +
		"&response_type=code" +
		"&scope=https://www.googleapis.com/auth/calendar" +
		"&access_type=offline"

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"authUrl": authURL})
}

func (h *IntegrationHandler) GoogleCallback(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req GoogleCallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if h.cfg.GoogleClientID == "" || h.cfg.GoogleClientSecret == "" {
		utils.RespondWithError(c, http.StatusServiceUnavailable, "Google integration not configured")
		return
	}

	// In a real application, you would:
	// 1. Exchange the code for tokens using Google's OAuth API
	// 2. Store the tokens in the database
	// 3. Set up refresh token handling

	// For now, we'll just store a placeholder
	_, err := database.DB.Exec(`
		INSERT INTO google_integrations (user_id, access_token, refresh_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET access_token = $2, refresh_token = $3, updated_at = NOW()
	`, userID, "placeholder_access_token", "placeholder_refresh_token")

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to save Google integration")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Google account linked successfully"})
}

func (h *IntegrationHandler) GoogleSync(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req GoogleSyncRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate direction
	if req.Direction != "import" && req.Direction != "export" && req.Direction != "both" {
		utils.RespondWithError(c, http.StatusBadRequest, "Direction must be 'import', 'export', or 'both'")
		return
	}

	// Check if user has Google integration
	var integrationID string
	err := database.DB.QueryRow("SELECT id FROM google_integrations WHERE user_id = $1", userID).Scan(&integrationID)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Google account not linked")
		return
	}

	// In a real application, you would:
	// 1. Use the stored tokens to access Google Calendar API
	// 2. Sync events based on the direction
	// 3. Handle pagination and rate limiting

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Calendar sync initiated", "direction": req.Direction})
}

func (h *IntegrationHandler) LinkTelegram(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req TelegramLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO telegram_integrations (user_id, telegram_user_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET telegram_user_id = $2
	`, userID, req.TelegramUserID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to link Telegram account")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Telegram account linked successfully"})
}
