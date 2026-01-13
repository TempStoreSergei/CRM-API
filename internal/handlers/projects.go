package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct{}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{}
}

type CreateProjectRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	TeamMembers []string `json:"teamMembers"`
}

type UpdateProjectRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Status          string   `json:"status"`
	StartDate       string   `json:"startDate"`
	EndDate         string   `json:"endDate"`
	BudgetAllocated float64  `json:"budgetAllocated"`
	BudgetSpent     float64  `json:"budgetSpent"`
	TeamMembers     []string `json:"teamMembers"`
}

type AddMemberRequest struct {
	UserID string `json:"userId" binding:"required"`
	Role   string `json:"role"`
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	status := c.Query("status")
	search := c.Query("search")

	// Build query
	query := `SELECT id, code, name, description, status, start_date, end_date, created_at FROM projects WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM projects WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", argIndex, argIndex)
		countQuery += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", argIndex, argIndex)
		args = append(args, searchTerm)
		argIndex++
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get projects count")
		return
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get projects")
		return
	}
	defer rows.Close()

	projects := []models.ProjectListItem{}
	for rows.Next() {
		var proj models.ProjectListItem
		var description sql.NullString
		var startDate, endDate sql.NullTime

		err := rows.Scan(&proj.ID, &proj.Code, &proj.Name, &description, &proj.Status, &startDate, &endDate, &proj.CreatedAt)
		if err != nil {
			continue
		}

		if description.Valid {
			proj.Description = description.String
		}
		if startDate.Valid {
			proj.StartDate = startDate.Time.Format("2006-01-02")
		}
		if endDate.Valid {
			proj.EndDate = endDate.Time.Format("2006-01-02")
		}

		// Get team members
		proj.Team = h.getProjectTeam(proj.ID)

		// Get tasks count
		proj.TasksCount = h.getTasksCount(proj.ID)

		// Calculate progress
		if proj.TasksCount.Total > 0 {
			proj.Progress = (proj.TasksCount.Completed * 100) / proj.TasksCount.Total
		}

		projects = append(projects, proj)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": projects,
		"meta": models.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: utils.CalculateTotalPages(total, limit),
		},
	})
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	// Generate project code
	var lastCode int
	database.DB.QueryRow("SELECT COALESCE(MAX(CAST(SUBSTRING(code FROM 3) AS INTEGER)), 0) FROM projects WHERE code LIKE 'PN%'").Scan(&lastCode)
	code := fmt.Sprintf("PN%07d", lastCode+1)

	// Parse dates
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsed, _ := time.Parse("2006-01-02", req.StartDate)
		startDate = &parsed
	}
	if req.EndDate != "" {
		parsed, _ := time.Parse("2006-01-02", req.EndDate)
		endDate = &parsed
	}

	var projectID string
	err := database.DB.QueryRow(`
		INSERT INTO projects (code, name, description, start_date, end_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, code, req.Name, nullString(req.Description), startDate, endDate, userID).Scan(&projectID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create project")
		return
	}

	// Add team members
	for _, memberID := range req.TeamMembers {
		database.DB.Exec(`
			INSERT INTO project_members (project_id, user_id, role)
			VALUES ($1, $2, 'member')
			ON CONFLICT DO NOTHING
		`, projectID, memberID)
	}

	// Log activity
	logActivity(userID, "project_created", "Created project: "+req.Name, "project", projectID)

	// Get and return the created project
	project := h.getProjectByID(projectID)
	utils.RespondWithSuccess(c, http.StatusCreated, project)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectID := c.Param("id")

	project := h.getProjectByID(projectID)
	if project == nil {
		utils.RespondWithError(c, http.StatusNotFound, "Project not found")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, project)
}

func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	projectID := c.Param("id")

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Build update query
	query := "UPDATE projects SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if req.Name != "" {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, req.Name)
		argIndex++
	}
	if req.Description != "" {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, req.Description)
		argIndex++
	}
	if req.Status != "" {
		query += fmt.Sprintf(", status = $%d", argIndex)
		args = append(args, req.Status)
		argIndex++
	}
	if req.StartDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		query += fmt.Sprintf(", start_date = $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}
	if req.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		query += fmt.Sprintf(", end_date = $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}
	if req.BudgetAllocated > 0 {
		query += fmt.Sprintf(", budget_allocated = $%d", argIndex)
		args = append(args, req.BudgetAllocated)
		argIndex++
	}
	if req.BudgetSpent > 0 {
		query += fmt.Sprintf(", budget_spent = $%d", argIndex)
		args = append(args, req.BudgetSpent)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, projectID)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update project")
		return
	}

	// Update team members if provided
	if req.TeamMembers != nil {
		database.DB.Exec("DELETE FROM project_members WHERE project_id = $1", projectID)
		for _, memberID := range req.TeamMembers {
			database.DB.Exec(`
				INSERT INTO project_members (project_id, user_id, role)
				VALUES ($1, $2, 'member')
				ON CONFLICT DO NOTHING
			`, projectID, memberID)
		}
	}

	project := h.getProjectByID(projectID)
	utils.RespondWithSuccess(c, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	// Archive the project instead of deleting
	_, err := database.DB.Exec("UPDATE projects SET status = 'archived', updated_at = NOW() WHERE id = $1", projectID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Project archived successfully"})
}

func (h *ProjectHandler) AddMember(c *gin.Context) {
	projectID := c.Param("id")

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	_, err := database.DB.Exec(`
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, user_id) DO UPDATE SET role = $3
	`, projectID, req.UserID, role)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to add member")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Member added successfully"})
}

func (h *ProjectHandler) RemoveMember(c *gin.Context) {
	projectID := c.Param("id")
	userID := c.Param("userId")

	_, err := database.DB.Exec("DELETE FROM project_members WHERE project_id = $1 AND user_id = $2", projectID, userID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to remove member")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Member removed successfully"})
}

func (h *ProjectHandler) getProjectByID(projectID string) interface{} {
	var project struct {
		ID              string
		Code            string
		Name            string
		Description     sql.NullString
		Status          string
		StartDate       sql.NullTime
		EndDate         sql.NullTime
		BudgetAllocated sql.NullFloat64
		BudgetSpent     float64
		BudgetCurrency  string
		CreatedAt       time.Time
		UpdatedAt       time.Time
	}

	err := database.DB.QueryRow(`
		SELECT id, code, name, description, status, start_date, end_date, budget_allocated, COALESCE(budget_spent, 0), COALESCE(budget_currency, 'RUB'), created_at, updated_at
		FROM projects WHERE id = $1
	`, projectID).Scan(&project.ID, &project.Code, &project.Name, &project.Description, &project.Status, &project.StartDate, &project.EndDate, &project.BudgetAllocated, &project.BudgetSpent, &project.BudgetCurrency, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		return nil
	}

	result := gin.H{
		"id":        project.ID,
		"code":      project.Code,
		"name":      project.Name,
		"status":    project.Status,
		"createdAt": project.CreatedAt,
		"updatedAt": project.UpdatedAt,
	}

	if project.Description.Valid {
		result["description"] = project.Description.String
	}
	if project.StartDate.Valid {
		result["startDate"] = project.StartDate.Time.Format("2006-01-02")
	}
	if project.EndDate.Valid {
		result["endDate"] = project.EndDate.Time.Format("2006-01-02")
	}

	// Budget
	if project.BudgetAllocated.Valid {
		result["budget"] = models.ProjectBudget{
			Allocated: project.BudgetAllocated.Float64,
			Spent:     project.BudgetSpent,
			Currency:  project.BudgetCurrency,
		}
	}

	// Team
	result["team"] = h.getProjectTeam(projectID)

	// Tasks count and progress
	tasksCount := h.getTasksCount(projectID)
	result["progress"] = 0
	if tasksCount.Total > 0 {
		result["progress"] = (tasksCount.Completed * 100) / tasksCount.Total
	}

	// Get tasks
	result["tasks"] = h.getProjectTasks(projectID)

	// Get files
	result["files"] = h.getProjectFiles(projectID)

	return result
}

func (h *ProjectHandler) getProjectTeam(projectID string) []models.ProjectTeamMember {
	rows, err := database.DB.Query(`
		SELECT u.id, CONCAT(u.first_name, ' ', u.last_name) as name, u.avatar_url, pm.role, u.position
		FROM project_members pm
		JOIN users u ON pm.user_id = u.id
		WHERE pm.project_id = $1
	`, projectID)

	if err != nil {
		return []models.ProjectTeamMember{}
	}
	defer rows.Close()

	team := []models.ProjectTeamMember{}
	for rows.Next() {
		var member models.ProjectTeamMember
		var avatar, role, position sql.NullString
		rows.Scan(&member.ID, &member.Name, &avatar, &role, &position)
		if avatar.Valid {
			member.Avatar = avatar.String
		}
		if role.Valid {
			member.Role = role.String
		}
		if position.Valid {
			member.Position = position.String
		}
		team = append(team, member)
	}

	return team
}

func (h *ProjectHandler) getTasksCount(projectID string) models.TasksCount {
	var count models.TasksCount
	database.DB.QueryRow(`
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'backlog' THEN 1 ELSE 0 END), 0)
		FROM tasks WHERE project_id = $1
	`, projectID).Scan(&count.Total, &count.Completed, &count.InProgress, &count.Backlog)

	return count
}

func (h *ProjectHandler) getProjectTasks(projectID string) []gin.H {
	rows, err := database.DB.Query(`
		SELECT t.id, t.title, t.status, t.assignee_id, CONCAT(u.first_name, ' ', u.last_name) as assignee_name
		FROM tasks t
		LEFT JOIN users u ON t.assignee_id = u.id
		WHERE t.project_id = $1
		LIMIT 10
	`, projectID)

	if err != nil {
		return []gin.H{}
	}
	defer rows.Close()

	tasks := []gin.H{}
	for rows.Next() {
		var id, title, status string
		var assigneeID, assigneeName sql.NullString
		rows.Scan(&id, &title, &status, &assigneeID, &assigneeName)

		task := gin.H{
			"id":     id,
			"title":  title,
			"status": status,
		}
		if assigneeID.Valid && assigneeName.Valid {
			task["assignee"] = gin.H{"id": assigneeID.String, "name": assigneeName.String}
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func (h *ProjectHandler) getProjectFiles(projectID string) []models.File {
	rows, err := database.DB.Query(`
		SELECT id, name, url, type, size, created_at
		FROM files WHERE entity_type = 'project' AND entity_id = $1
	`, projectID)

	if err != nil {
		return []models.File{}
	}
	defer rows.Close()

	files := []models.File{}
	for rows.Next() {
		var f models.File
		var fileType sql.NullString
		rows.Scan(&f.ID, &f.Name, &f.URL, &fileType, &f.Size, &f.UploadedAt)
		if fileType.Valid {
			f.Type = fileType.String
		}
		files = append(files, f)
	}

	return files
}

func logActivity(userID, activityType, description, entityType, entityID string) {
	database.DB.Exec(`
		INSERT INTO activity_log (user_id, type, description, entity_type, entity_id)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, activityType, description, entityType, entityID)
}
