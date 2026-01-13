package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct{}

func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{}
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	department := c.Query("department")
	status := c.Query("status")

	// Build query
	query := `
		SELECT u.id, u.id as user_id, u.first_name, u.last_name, u.email, u.avatar_url, u.position, u.level, u.department, u.status
		FROM users u WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if department != "" {
		query += fmt.Sprintf(" AND u.department = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND department = $%d", argIndex)
		args = append(args, department)
		argIndex++
	}

	// Map frontend status to internal representation
	if status != "" {
		switch status {
		case "on_vacation":
			query += fmt.Sprintf(" AND u.id IN (SELECT user_id FROM vacations WHERE status = 'approved' AND type = 'annual' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)")
			countQuery += fmt.Sprintf(" AND id IN (SELECT user_id FROM vacations WHERE status = 'approved' AND type = 'annual' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)")
		case "sick_leave":
			query += fmt.Sprintf(" AND u.id IN (SELECT user_id FROM vacations WHERE status = 'approved' AND type = 'sick' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)")
			countQuery += fmt.Sprintf(" AND id IN (SELECT user_id FROM vacations WHERE status = 'approved' AND type = 'sick' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)")
		case "active":
			query += fmt.Sprintf(" AND u.status = 'active' AND u.id NOT IN (SELECT user_id FROM vacations WHERE status = 'approved' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)")
			countQuery += fmt.Sprintf(" AND status = 'active' AND id NOT IN (SELECT user_id FROM vacations WHERE status = 'approved' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE)")
		}
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get employees count")
		return
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY u.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get employees")
		return
	}
	defer rows.Close()

	employees := []models.Employee{}
	for rows.Next() {
		var emp models.Employee
		var avatar, position, level, department sql.NullString

		err := rows.Scan(&emp.ID, &emp.UserID, &emp.FirstName, &emp.LastName, &emp.Email, &avatar, &position, &level, &department, &emp.Status)
		if err != nil {
			continue
		}

		if avatar.Valid {
			emp.Avatar = avatar.String
		}
		if position.Valid {
			emp.Position = position.String
		}
		if level.Valid {
			emp.Level = level.String
		}
		if department.Valid {
			emp.Department = department.String
		}

		// Determine actual status (vacation, sick leave, etc.)
		emp.Status = h.getEmployeeStatus(emp.UserID, emp.Status)

		// Get workload
		emp.Workload = h.getEmployeeWorkload(emp.UserID)

		employees = append(employees, emp)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": employees,
		"meta": models.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: utils.CalculateTotalPages(total, limit),
		},
	})
}

func (h *EmployeeHandler) GetEmployeeTasks(c *gin.Context) {
	employeeID := c.Param("id")

	// Get tasks for this employee
	rows, err := database.DB.Query(`
		SELECT t.id, t.title, t.description, t.status, t.priority, t.project_id, p.name as project_name, t.due_date, t.created_at
		FROM tasks t
		LEFT JOIN projects p ON t.project_id = p.id
		WHERE t.assignee_id = $1
		ORDER BY t.created_at DESC
	`, employeeID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get tasks")
		return
	}
	defer rows.Close()

	tasks := []gin.H{}
	for rows.Next() {
		var id, title, status, priority string
		var description, projectID, projectName sql.NullString
		var dueDate sql.NullTime
		var createdAt interface{}

		err := rows.Scan(&id, &title, &description, &status, &priority, &projectID, &projectName, &dueDate, &createdAt)
		if err != nil {
			continue
		}

		task := gin.H{
			"id":          id,
			"title":       title,
			"description": "",
			"status":      status,
			"priority":    priority,
			"createdAt":   createdAt,
		}

		if description.Valid {
			task["description"] = description.String
		}

		if projectID.Valid && projectName.Valid {
			task["project"] = gin.H{
				"id":   projectID.String,
				"name": projectName.String,
			}
		}

		if dueDate.Valid {
			task["dueDate"] = dueDate.Time
		}

		tasks = append(tasks, task)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"data": tasks})
}

func (h *EmployeeHandler) getEmployeeStatus(userID, baseStatus string) string {
	// Check if user is on vacation
	var vacationType string
	err := database.DB.QueryRow(`
		SELECT type FROM vacations 
		WHERE user_id = $1 AND status = 'approved' AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE
		LIMIT 1
	`, userID).Scan(&vacationType)

	if err == nil {
		switch strings.ToLower(vacationType) {
		case "annual", "unpaid", "maternity":
			return "on_vacation"
		case "sick":
			return "sick_leave"
		}
	}

	return baseStatus
}

func (h *EmployeeHandler) getEmployeeWorkload(userID string) models.EmployeeWorkload {
	workload := models.EmployeeWorkload{}

	database.DB.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN status = 'backlog' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'in_review' THEN 1 ELSE 0 END), 0)
		FROM tasks WHERE assignee_id = $1
	`, userID).Scan(&workload.BacklogTasks, &workload.InProgressTasks, &workload.InReviewTasks)

	return workload
}
