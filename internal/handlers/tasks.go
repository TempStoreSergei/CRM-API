package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct{}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{}
}

type CreateTaskRequest struct {
	Title          string   `json:"title" binding:"required"`
	Description    string   `json:"description"`
	ProjectID      string   `json:"projectId"`
	AssigneeID     string   `json:"assigneeId"`
	Priority       string   `json:"priority"`
	DueDate        string   `json:"dueDate"`
	EstimatedHours float64  `json:"estimatedHours"`
	Tags           []string `json:"tags"`
}

type UpdateTaskRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	ProjectID      string   `json:"projectId"`
	AssigneeID     string   `json:"assigneeId"`
	Priority       string   `json:"priority"`
	Status         string   `json:"status"`
	DueDate        string   `json:"dueDate"`
	EstimatedHours float64  `json:"estimatedHours"`
	LoggedHours    float64  `json:"loggedHours"`
	Tags           []string `json:"tags"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type AddCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	projectID := c.Query("projectId")
	assigneeID := c.Query("assigneeId")
	status := c.Query("status")
	priority := c.Query("priority")

	// Build query
	query := `
		SELECT t.id, t.title, t.description, t.status, t.priority, 
			t.project_id, p.name as project_name, p.code as project_code,
			t.assignee_id, CONCAT(ua.first_name, ' ', ua.last_name) as assignee_name, ua.avatar_url as assignee_avatar,
			t.reporter_id, CONCAT(ur.first_name, ' ', ur.last_name) as reporter_name,
			t.due_date, t.estimated_hours, t.logged_hours, t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN projects p ON t.project_id = p.id
		LEFT JOIN users ua ON t.assignee_id = ua.id
		LEFT JOIN users ur ON t.reporter_id = ur.id
		WHERE 1=1
	`
	countQuery := "SELECT COUNT(*) FROM tasks t WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if projectID != "" {
		query += fmt.Sprintf(" AND t.project_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND project_id = $%d", argIndex)
		args = append(args, projectID)
		argIndex++
	}

	if assigneeID != "" {
		query += fmt.Sprintf(" AND t.assignee_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND assignee_id = $%d", argIndex)
		args = append(args, assigneeID)
		argIndex++
	}

	if status != "" {
		query += fmt.Sprintf(" AND t.status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if priority != "" {
		query += fmt.Sprintf(" AND t.priority = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND priority = $%d", argIndex)
		args = append(args, priority)
		argIndex++
	}

	// Get total count
	var total int
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get tasks count")
		return
	}

	// Add pagination
	query += fmt.Sprintf(" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get tasks")
		return
	}
	defer rows.Close()

	tasks := []models.TaskListItem{}
	for rows.Next() {
		var task models.TaskListItem
		var description, projectID, projectName, projectCode sql.NullString
		var assigneeID, assigneeName, assigneeAvatar sql.NullString
		var reporterID, reporterName sql.NullString
		var dueDate sql.NullTime
		var estimatedHours, loggedHours sql.NullFloat64

		err := rows.Scan(&task.ID, &task.Title, &description, &task.Status, &task.Priority,
			&projectID, &projectName, &projectCode,
			&assigneeID, &assigneeName, &assigneeAvatar,
			&reporterID, &reporterName,
			&dueDate, &estimatedHours, &loggedHours, &task.CreatedAt, &task.UpdatedAt)

		if err != nil {
			continue
		}

		if description.Valid {
			task.Description = description.String
		}

		if projectID.Valid && projectName.Valid {
			task.Project = &models.TaskProject{
				ID:   projectID.String,
				Name: projectName.String,
				Code: projectCode.String,
			}
		}

		if assigneeID.Valid && assigneeName.Valid {
			task.Assignee = &models.TaskUser{
				ID:   assigneeID.String,
				Name: assigneeName.String,
			}
			if assigneeAvatar.Valid {
				task.Assignee.Avatar = assigneeAvatar.String
			}
		}

		if reporterID.Valid && reporterName.Valid {
			task.Reporter = &models.TaskUser{
				ID:   reporterID.String,
				Name: reporterName.String,
			}
		}

		if dueDate.Valid {
			task.DueDate = dueDate.Time.Format(time.RFC3339)
		}

		if estimatedHours.Valid {
			task.EstimatedHours = estimatedHours.Float64
		}

		if loggedHours.Valid {
			task.LoggedHours = loggedHours.Float64
		}

		// Get tags
		task.Tags = h.getTaskTags(task.ID)

		tasks = append(tasks, task)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": tasks,
		"meta": models.PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: utils.CalculateTotalPages(total, limit),
		},
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, _ := time.Parse(time.RFC3339, req.DueDate)
		dueDate = &parsed
	}

	var taskID string
	err := database.DB.QueryRow(`
		INSERT INTO tasks (title, description, project_id, assignee_id, reporter_id, priority, due_date, estimated_hours)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, req.Title, nullString(req.Description), nullStringPtr(req.ProjectID), nullStringPtr(req.AssigneeID), userID, priority, dueDate, req.EstimatedHours).Scan(&taskID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// Add tags
	for _, tag := range req.Tags {
		if tag != "" {
			database.DB.Exec("INSERT INTO task_tags (task_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING", taskID, tag)
		}
	}

	// Log activity
	logActivity(userID, "task_created", "Created task: "+req.Title, "task", taskID)

	// Send notification if assigned
	if req.AssigneeID != "" && req.AssigneeID != userID {
		createNotification(req.AssigneeID, "task_assigned", "New task assigned", "You have been assigned to: "+req.Title, "task", taskID)
	}

	// Get and return the created task
	task := h.getTaskByID(taskID)
	utils.RespondWithSuccess(c, http.StatusCreated, task)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("id")

	task := h.getTaskByID(taskID)
	if task == nil {
		utils.RespondWithError(c, http.StatusNotFound, "Task not found")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Build update query
	query := "UPDATE tasks SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	if req.Title != "" {
		query += fmt.Sprintf(", title = $%d", argIndex)
		args = append(args, req.Title)
		argIndex++
	}
	if req.Description != "" {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, req.Description)
		argIndex++
	}
	if req.ProjectID != "" {
		query += fmt.Sprintf(", project_id = $%d", argIndex)
		args = append(args, req.ProjectID)
		argIndex++
	}
	if req.AssigneeID != "" {
		query += fmt.Sprintf(", assignee_id = $%d", argIndex)
		args = append(args, req.AssigneeID)
		argIndex++
	}
	if req.Priority != "" {
		query += fmt.Sprintf(", priority = $%d", argIndex)
		args = append(args, req.Priority)
		argIndex++
	}
	if req.Status != "" {
		query += fmt.Sprintf(", status = $%d", argIndex)
		args = append(args, req.Status)
		argIndex++
	}
	if req.DueDate != "" {
		dueDate, _ := time.Parse(time.RFC3339, req.DueDate)
		query += fmt.Sprintf(", due_date = $%d", argIndex)
		args = append(args, dueDate)
		argIndex++
	}
	if req.EstimatedHours > 0 {
		query += fmt.Sprintf(", estimated_hours = $%d", argIndex)
		args = append(args, req.EstimatedHours)
		argIndex++
	}
	if req.LoggedHours > 0 {
		query += fmt.Sprintf(", logged_hours = $%d", argIndex)
		args = append(args, req.LoggedHours)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, taskID)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update task")
		return
	}

	// Update tags if provided
	if req.Tags != nil {
		database.DB.Exec("DELETE FROM task_tags WHERE task_id = $1", taskID)
		for _, tag := range req.Tags {
			if tag != "" {
				database.DB.Exec("INSERT INTO task_tags (task_id, tag) VALUES ($1, $2) ON CONFLICT DO NOTHING", taskID, tag)
			}
		}
	}

	task := h.getTaskByID(taskID)
	utils.RespondWithSuccess(c, http.StatusOK, task)
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	taskID := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate status
	validStatuses := []string{"backlog", "in_progress", "in_review", "done"}
	isValid := false
	for _, s := range validStatuses {
		if req.Status == s {
			isValid = true
			break
		}
	}
	if !isValid {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid status")
		return
	}

	_, err := database.DB.Exec("UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2", req.Status, taskID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update status")
		return
	}

	userID := middleware.GetUserID(c)
	if req.Status == "done" {
		logActivity(userID, "task_completed", "Completed task", "task", taskID)
	}

	task := h.getTaskByID(taskID)
	utils.RespondWithSuccess(c, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete task")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) AddComment(c *gin.Context) {
	taskID := c.Param("id")
	userID := middleware.GetUserID(c)

	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var commentID string
	err := database.DB.QueryRow(`
		INSERT INTO task_comments (task_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`, taskID, userID, req.Content).Scan(&commentID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to add comment")
		return
	}

	// Get the created comment
	var comment models.TaskComment
	var userName sql.NullString
	var userAvatar sql.NullString

	database.DB.QueryRow(`
		SELECT c.id, c.task_id, c.content, c.created_at, c.updated_at, u.id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url
		FROM task_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`, commentID).Scan(&comment.ID, &comment.TaskID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.ID, &userName, &userAvatar)

	if userName.Valid {
		comment.User.Name = userName.String
	}
	if userAvatar.Valid {
		comment.User.Avatar = userAvatar.String
	}

	utils.RespondWithSuccess(c, http.StatusCreated, comment)
}

func (h *TaskHandler) GetComments(c *gin.Context) {
	taskID := c.Param("id")

	rows, err := database.DB.Query(`
		SELECT c.id, c.task_id, c.content, c.created_at, c.updated_at, u.id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url
		FROM task_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.task_id = $1
		ORDER BY c.created_at ASC
	`, taskID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get comments")
		return
	}
	defer rows.Close()

	comments := []models.TaskComment{}
	for rows.Next() {
		var comment models.TaskComment
		var userName sql.NullString
		var userAvatar sql.NullString

		rows.Scan(&comment.ID, &comment.TaskID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt, &comment.User.ID, &userName, &userAvatar)

		if userName.Valid {
			comment.User.Name = userName.String
		}
		if userAvatar.Valid {
			comment.User.Avatar = userAvatar.String
		}

		comments = append(comments, comment)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"data": comments})
}

func (h *TaskHandler) AttachFile(c *gin.Context) {
	taskID := c.Param("id")
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "No file uploaded")
		return
	}
	defer file.Close()

	// In a real application, you would upload to S3 or similar storage
	// For now, we'll just record the file info
	fileURL := "/uploads/tasks/" + header.Filename

	var fileID string
	err = database.DB.QueryRow(`
		INSERT INTO files (name, url, type, size, uploaded_by, entity_type, entity_id)
		VALUES ($1, $2, $3, $4, $5, 'task', $6)
		RETURNING id
	`, header.Filename, fileURL, header.Header.Get("Content-Type"), header.Size, userID, taskID).Scan(&fileID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to attach file")
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, models.File{
		ID:         fileID,
		Name:       header.Filename,
		URL:        fileURL,
		Type:       header.Header.Get("Content-Type"),
		Size:       int(header.Size),
		UploadedAt: time.Now(),
	})
}

func (h *TaskHandler) getTaskByID(taskID string) interface{} {
	var task struct {
		ID             string
		Title          string
		Description    sql.NullString
		Status         string
		Priority       string
		ProjectID      sql.NullString
		ProjectName    sql.NullString
		ProjectCode    sql.NullString
		AssigneeID     sql.NullString
		AssigneeName   sql.NullString
		AssigneeAvatar sql.NullString
		ReporterID     sql.NullString
		ReporterName   sql.NullString
		DueDate        sql.NullTime
		EstimatedHours sql.NullFloat64
		LoggedHours    float64
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}

	err := database.DB.QueryRow(`
		SELECT t.id, t.title, t.description, t.status, t.priority,
			t.project_id, p.name, p.code,
			t.assignee_id, CONCAT(ua.first_name, ' ', ua.last_name), ua.avatar_url,
			t.reporter_id, CONCAT(ur.first_name, ' ', ur.last_name),
			t.due_date, t.estimated_hours, COALESCE(t.logged_hours, 0), t.created_at, t.updated_at
		FROM tasks t
		LEFT JOIN projects p ON t.project_id = p.id
		LEFT JOIN users ua ON t.assignee_id = ua.id
		LEFT JOIN users ur ON t.reporter_id = ur.id
		WHERE t.id = $1
	`, taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority,
		&task.ProjectID, &task.ProjectName, &task.ProjectCode,
		&task.AssigneeID, &task.AssigneeName, &task.AssigneeAvatar,
		&task.ReporterID, &task.ReporterName,
		&task.DueDate, &task.EstimatedHours, &task.LoggedHours, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil
	}

	result := gin.H{
		"id":          task.ID,
		"title":       task.Title,
		"status":      task.Status,
		"priority":    task.Priority,
		"loggedHours": task.LoggedHours,
		"createdAt":   task.CreatedAt,
		"updatedAt":   task.UpdatedAt,
		"tags":        h.getTaskTags(taskID),
	}

	if task.Description.Valid {
		result["description"] = task.Description.String
	}
	if task.ProjectID.Valid {
		result["project"] = gin.H{
			"id":   task.ProjectID.String,
			"name": task.ProjectName.String,
			"code": task.ProjectCode.String,
		}
	}
	if task.AssigneeID.Valid {
		result["assignee"] = gin.H{
			"id":     task.AssigneeID.String,
			"name":   task.AssigneeName.String,
			"avatar": task.AssigneeAvatar.String,
		}
	}
	if task.ReporterID.Valid {
		result["reporter"] = gin.H{
			"id":   task.ReporterID.String,
			"name": task.ReporterName.String,
		}
	}
	if task.DueDate.Valid {
		result["dueDate"] = task.DueDate.Time
	}
	if task.EstimatedHours.Valid {
		result["estimatedHours"] = task.EstimatedHours.Float64
	}

	return result
}

func (h *TaskHandler) getTaskTags(taskID string) []string {
	rows, err := database.DB.Query("SELECT tag FROM task_tags WHERE task_id = $1", taskID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	tags := []string{}
	for rows.Next() {
		var tag string
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	return tags
}

func nullStringPtr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func createNotification(userID, notifType, title, message, entityType, entityID string) {
	database.DB.Exec(`
		INSERT INTO notifications (user_id, type, title, message, data)
		VALUES ($1, $2, $3, $4, $5)
	`, userID, notifType, title, message, fmt.Sprintf(`{"entityType":"%s","entityId":"%s"}`, entityType, entityID))
}
