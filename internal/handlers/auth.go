package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/config"
	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type SignupRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Agent     string `json:"agent"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Agent    string `json:"agent"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	Agent        string `json:"agent"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SignupResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"`
}

type LoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresIn    int64     `json:"expiresIn"`
	User         LoginUser `json:"user"`
}

type LoginUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Avatar    string `json:"avatar,omitempty"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user already exists
	var existingID string
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		utils.RespondWithError(c, http.StatusConflict, "User with this email already exists")
		return
	}
	if err != sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	var userID string
	var createdAt time.Time
	err = database.DB.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, req.Email, hashedPassword, req.FirstName, req.LastName).Scan(&userID, &createdAt)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Initialize vacation balance for current year
	currentYear := time.Now().Year()
	_, err = database.DB.Exec(`
		INSERT INTO vacation_balances (user_id, year, annual_total, annual_used, sick_total, sick_used)
		VALUES ($1, $2, 28, 0, 10, 0)
		ON CONFLICT (user_id, year) DO NOTHING
	`, userID, currentYear)
	if err != nil {
		// Log but don't fail
	}

	utils.RespondWithSuccess(c, http.StatusCreated, SignupResponse{
		ID:        userID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: createdAt,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Get user from database
	var user struct {
		ID           string
		Email        string
		PasswordHash string
		FirstName    string
		LastName     string
		Role         string
		Avatar       sql.NullString
	}

	err := database.DB.QueryRow(`
		SELECT id, email, password_hash, first_name, last_name, role, avatar_url
		FROM users WHERE email = $1
	`, req.Email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Role, &user.Avatar)

	if err == sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate tokens
	tokenPair, err := utils.GenerateTokenPair(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTExpiresIn)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Store refresh token
	expiresAt := time.Now().Add(h.cfg.JWTRefreshExpires)
	_, err = database.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, user_agent, expires_at)
		VALUES ($1, $2, $3, $4)
	`, user.ID, tokenPair.RefreshToken, req.Agent, expiresAt)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to store refresh token")
		return
	}

	avatar := ""
	if user.Avatar.Valid {
		avatar = user.Avatar.String
	}

	utils.RespondWithSuccess(c, http.StatusOK, LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		User: LoginUser{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			Avatar:    avatar,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Find and validate refresh token
	var tokenData struct {
		UserID    string
		ExpiresAt time.Time
	}

	err := database.DB.QueryRow(`
		SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1
	`, req.RefreshToken).Scan(&tokenData.UserID, &tokenData.ExpiresAt)

	if err == sql.ErrNoRows {
		utils.RespondWithError(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	if time.Now().After(tokenData.ExpiresAt) {
		// Delete expired token
		database.DB.Exec("DELETE FROM refresh_tokens WHERE token = $1", req.RefreshToken)
		utils.RespondWithError(c, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	// Get user details
	var user struct {
		Email string
		Role  string
	}

	err = database.DB.QueryRow(`
		SELECT email, role FROM users WHERE id = $1
	`, tokenData.UserID).Scan(&user.Email, &user.Role)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "User not found")
		return
	}

	// Generate new token pair
	newTokenPair, err := utils.GenerateTokenPair(tokenData.UserID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTExpiresIn)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Delete old refresh token and create new one
	database.DB.Exec("DELETE FROM refresh_tokens WHERE token = $1", req.RefreshToken)

	expiresAt := time.Now().Add(h.cfg.JWTRefreshExpires)
	_, err = database.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, user_agent, expires_at)
		VALUES ($1, $2, $3, $4)
	`, tokenData.UserID, newTokenPair.RefreshToken, req.Agent, expiresAt)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to store refresh token")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, RefreshResponse{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
		ExpiresIn:    newTokenPair.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Delete all refresh tokens for this user
	_, err := database.DB.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user exists
	var userID string
	err := database.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		// Don't reveal if user exists
		utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Password reset instructions sent to email"})
		return
	}
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Database error")
		return
	}

	// In a real application, you would:
	// 1. Generate a password reset token
	// 2. Store it in the database with an expiry
	// 3. Send an email with a reset link

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Password reset instructions sent to email"})
}
