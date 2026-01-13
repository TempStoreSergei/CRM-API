package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type VacationHandler struct{}

func NewVacationHandler() *VacationHandler {
	return &VacationHandler{}
}

type CreateVacationRequest struct {
	Type      string `json:"type" binding:"required"`
	StartDate string `json:"startDate" binding:"required"`
	EndDate   string `json:"endDate" binding:"required"`
	Reason    string `json:"reason"`
}

type RejectVacationRequest struct {
	Reason string `json:"reason"`
}

func (h *VacationHandler) GetVacations(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	status := c.Query("status")
	userID := c.Query("userId")
	yearStr := c.Query("year")

	// Build query
	query := `
		SELECT v.id, v.user_id, u.first_name, u.last_name, u.email, u.avatar_url, u.position,
			v.type, v.start_date, v.end_date, v.total_days, v.status, v.reason,
			v.approved_by, CONCAT(a.first_name, ' ', a.last_name) as approver_name,
			v.created_at
		FROM vacations v
		JOIN users u ON v.user_id = u.id
		LEFT JOIN users a ON v.approved_by = a.id
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM vacations WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		query += fmt.Sprintf(" AND v.status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if userID != "" {
		query += fmt.Sprintf(" AND v.user_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if yearStr != "" {
		year, _ := strconv.Atoi(yearStr)
		query += fmt.Sprintf(" AND EXTRACT(YEAR FROM v.start_date) = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND EXTRACT(YEAR FROM start_date) = $%d", argIndex)
		args = append(args, year)
		argIndex++
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get vacations count")
		return
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY v.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get vacations")
		return
	}
	defer rows.Close()

	vacations := []models.VacationListItem{}
	for rows.Next() {
		var v models.VacationListItem
		var avatar, position, reason sql.NullString
		var approvedBy, approverName sql.NullString
		var startDate, endDate time.Time

		err := rows.Scan(&v.ID, &v.User.ID, &v.User.Name, &v.User.Email, &v.User.Email, &avatar, &position,
			&v.Type, &startDate, &endDate, &v.TotalDays, &v.Status, &reason,
			&approvedBy, &approverName, &v.CreatedAt)

		if err != nil {
			continue
		}

		// Fix user name (concat first and last name)
		var firstName, lastName string
		database.DB.QueryRow("SELECT first_name, last_name FROM users WHERE id = $1", v.User.ID).Scan(&firstName, &lastName)
		v.User.Name = firstName + " " + lastName

		if avatar.Valid {
			v.User.Avatar = avatar.String
		}
		if position.Valid {
			v.User.Position = position.String
		}
		if reason.Valid {
			v.Reason = reason.String
		}

		v.StartDate = startDate.Format("2006-01-02")
		v.EndDate = endDate.Format("2006-01-02")

		if approvedBy.Valid && approverName.Valid {
			v.ApprovedBy = &models.VacationApprover{
				ID:   approvedBy.String,
				Name: approverName.String,
			}
		}

		vacations = append(vacations, v)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": vacations,
		"meta": models.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: utils.CalculateTotalPages(total, limit),
		},
	})
}

func (h *VacationHandler) CreateVacation(c *gin.Context) {
	var req CreateVacationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid start date format (use YYYY-MM-DD)")
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid end date format (use YYYY-MM-DD)")
		return
	}

	if endDate.Before(startDate) {
		utils.RespondWithError(c, http.StatusBadRequest, "End date must be after start date")
		return
	}

	// Calculate total days (excluding weekends)
	totalDays := calculateWorkDays(startDate, endDate)

	// Check vacation balance
	balance := h.getVacationBalance(userID)
	if req.Type == "annual" && totalDays > balance.AnnualRemaining {
		utils.RespondWithError(c, http.StatusBadRequest, "Insufficient vacation days remaining")
		return
	}
	if req.Type == "sick" && totalDays > balance.SickRemaining {
		utils.RespondWithError(c, http.StatusBadRequest, "Insufficient sick days remaining")
		return
	}

	var vacationID string
	err = database.DB.QueryRow(`
		INSERT INTO vacations (user_id, type, start_date, end_date, total_days, reason)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, userID, req.Type, startDate, endDate, totalDays, nullString(req.Reason)).Scan(&vacationID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create vacation request")
		return
	}

	vacation := h.getVacationByID(vacationID)
	utils.RespondWithSuccess(c, http.StatusCreated, vacation)
}

func (h *VacationHandler) GetVacation(c *gin.Context) {
	vacationID := c.Param("id")

	vacation := h.getVacationByID(vacationID)
	if vacation == nil {
		utils.RespondWithError(c, http.StatusNotFound, "Vacation request not found")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, vacation)
}

func (h *VacationHandler) ApproveVacation(c *gin.Context) {
	vacationID := c.Param("id")
	approverID := middleware.GetUserID(c)
	approverRole := middleware.GetUserRole(c)

	// Only managers and admins can approve
	if approverRole != "manager" && approverRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Only managers can approve vacation requests")
		return
	}

	// Get vacation details
	var vacation struct {
		UserID    string
		Type      string
		TotalDays int
		Status    string
	}

	err := database.DB.QueryRow("SELECT user_id, type, total_days, status FROM vacations WHERE id = $1", vacationID).Scan(&vacation.UserID, &vacation.Type, &vacation.TotalDays, &vacation.Status)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Vacation request not found")
		return
	}

	if vacation.Status != "pending" {
		utils.RespondWithError(c, http.StatusBadRequest, "Vacation request is not pending")
		return
	}

	// Update vacation status
	_, err = database.DB.Exec(`
		UPDATE vacations SET status = 'approved', approved_by = $1, approved_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`, approverID, vacationID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to approve vacation")
		return
	}

	// Update vacation balance
	currentYear := time.Now().Year()
	if vacation.Type == "annual" {
		database.DB.Exec(`
			UPDATE vacation_balances SET annual_used = annual_used + $1
			WHERE user_id = $2 AND year = $3
		`, vacation.TotalDays, vacation.UserID, currentYear)
	} else if vacation.Type == "sick" {
		database.DB.Exec(`
			UPDATE vacation_balances SET sick_used = sick_used + $1
			WHERE user_id = $2 AND year = $3
		`, vacation.TotalDays, vacation.UserID, currentYear)
	}

	// Notify user
	createNotification(vacation.UserID, "vacation_approved", "Vacation Approved", "Your vacation request has been approved", "vacation", vacationID)

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Vacation request approved"})
}

func (h *VacationHandler) RejectVacation(c *gin.Context) {
	vacationID := c.Param("id")
	approverRole := middleware.GetUserRole(c)

	// Only managers and admins can reject
	if approverRole != "manager" && approverRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "Only managers can reject vacation requests")
		return
	}

	var req RejectVacationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Reason is optional
	}

	// Get vacation details
	var vacation struct {
		UserID string
		Status string
	}

	err := database.DB.QueryRow("SELECT user_id, status FROM vacations WHERE id = $1", vacationID).Scan(&vacation.UserID, &vacation.Status)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Vacation request not found")
		return
	}

	if vacation.Status != "pending" {
		utils.RespondWithError(c, http.StatusBadRequest, "Vacation request is not pending")
		return
	}

	// Update vacation status
	_, err = database.DB.Exec(`
		UPDATE vacations SET status = 'rejected', rejection_reason = $1, updated_at = NOW()
		WHERE id = $2
	`, nullString(req.Reason), vacationID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to reject vacation")
		return
	}

	// Notify user
	createNotification(vacation.UserID, "vacation_rejected", "Vacation Rejected", "Your vacation request has been rejected", "vacation", vacationID)

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Vacation request rejected"})
}

func (h *VacationHandler) DeleteVacation(c *gin.Context) {
	vacationID := c.Param("id")
	userID := middleware.GetUserID(c)
	userRole := middleware.GetUserRole(c)

	// Get vacation details
	var vacation struct {
		UserID string
		Status string
	}

	err := database.DB.QueryRow("SELECT user_id, status FROM vacations WHERE id = $1", vacationID).Scan(&vacation.UserID, &vacation.Status)
	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Vacation request not found")
		return
	}

	// Only owner or admin can delete, and only pending requests
	if vacation.UserID != userID && userRole != "admin" {
		utils.RespondWithError(c, http.StatusForbidden, "You can only cancel your own vacation requests")
		return
	}

	if vacation.Status != "pending" {
		utils.RespondWithError(c, http.StatusBadRequest, "Only pending vacation requests can be cancelled")
		return
	}

	_, err = database.DB.Exec("DELETE FROM vacations WHERE id = $1", vacationID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to cancel vacation request")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Vacation request cancelled"})
}

func (h *VacationHandler) GetBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)

	balance := h.getVacationBalance(userID)
	utils.RespondWithSuccess(c, http.StatusOK, balance)
}

func (h *VacationHandler) getVacationByID(vacationID string) interface{} {
	var v struct {
		ID              string
		UserID          string
		FirstName       string
		LastName        string
		Email           string
		Avatar          sql.NullString
		Position        sql.NullString
		Type            string
		StartDate       time.Time
		EndDate         time.Time
		TotalDays       int
		Status          string
		Reason          sql.NullString
		RejectionReason sql.NullString
		ApprovedBy      sql.NullString
		ApproverName    sql.NullString
		CreatedAt       time.Time
	}

	err := database.DB.QueryRow(`
		SELECT v.id, v.user_id, u.first_name, u.last_name, u.email, u.avatar_url, u.position,
			v.type, v.start_date, v.end_date, v.total_days, v.status, v.reason, v.rejection_reason,
			v.approved_by, CONCAT(a.first_name, ' ', a.last_name),
			v.created_at
		FROM vacations v
		JOIN users u ON v.user_id = u.id
		LEFT JOIN users a ON v.approved_by = a.id
		WHERE v.id = $1
	`, vacationID).Scan(&v.ID, &v.UserID, &v.FirstName, &v.LastName, &v.Email, &v.Avatar, &v.Position,
		&v.Type, &v.StartDate, &v.EndDate, &v.TotalDays, &v.Status, &v.Reason, &v.RejectionReason,
		&v.ApprovedBy, &v.ApproverName, &v.CreatedAt)

	if err != nil {
		return nil
	}

	result := gin.H{
		"id": v.ID,
		"user": gin.H{
			"id":    v.UserID,
			"name":  v.FirstName + " " + v.LastName,
			"email": v.Email,
		},
		"type":      v.Type,
		"startDate": v.StartDate.Format("2006-01-02"),
		"endDate":   v.EndDate.Format("2006-01-02"),
		"totalDays": v.TotalDays,
		"status":    v.Status,
		"createdAt": v.CreatedAt,
	}

	if v.Avatar.Valid {
		result["user"].(gin.H)["avatar"] = v.Avatar.String
	}
	if v.Position.Valid {
		result["user"].(gin.H)["position"] = v.Position.String
	}
	if v.Reason.Valid {
		result["reason"] = v.Reason.String
	}
	if v.RejectionReason.Valid {
		result["rejectionReason"] = v.RejectionReason.String
	}
	if v.ApprovedBy.Valid && v.ApproverName.Valid {
		result["approvedBy"] = gin.H{
			"id":   v.ApprovedBy.String,
			"name": v.ApproverName.String,
		}
	}

	return result
}

func (h *VacationHandler) getVacationBalance(userID string) models.VacationBalance {
	currentYear := time.Now().Year()

	var balance models.VacationBalance
	err := database.DB.QueryRow(`
		SELECT annual_total, annual_used, sick_total, sick_used
		FROM vacation_balances
		WHERE user_id = $1 AND year = $2
	`, userID, currentYear).Scan(&balance.AnnualTotal, &balance.AnnualUsed, &balance.SickTotal, &balance.SickUsed)

	if err != nil {
		// Create default balance if not exists
		database.DB.Exec(`
			INSERT INTO vacation_balances (user_id, year, annual_total, annual_used, sick_total, sick_used)
			VALUES ($1, $2, 28, 0, 10, 0)
			ON CONFLICT (user_id, year) DO NOTHING
		`, userID, currentYear)

		balance = models.VacationBalance{
			AnnualTotal: 28,
			AnnualUsed:  0,
			SickTotal:   10,
			SickUsed:    0,
		}
	}

	balance.AnnualRemaining = balance.AnnualTotal - balance.AnnualUsed
	balance.SickRemaining = balance.SickTotal - balance.SickUsed

	return balance
}

func calculateWorkDays(startDate, endDate time.Time) int {
	workDays := 0
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workDays++
		}
	}
	return workDays
}
