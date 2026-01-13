package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct{}

func NewAnalyticsHandler() *AnalyticsHandler {
	return &AnalyticsHandler{}
}

func (h *AnalyticsHandler) GetDashboard(c *gin.Context) {
	dashboard := models.DashboardAnalytics{}

	// Get project stats
	database.DB.QueryRow(`
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'active' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END), 0)
		FROM projects
	`).Scan(&dashboard.Projects.Total, &dashboard.Projects.Active, &dashboard.Projects.Completed)

	// Get task stats
	database.DB.QueryRow(`
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN due_date < NOW() AND status != 'done' THEN 1 ELSE 0 END), 0)
		FROM tasks
	`).Scan(&dashboard.Tasks.Total, &dashboard.Tasks.Completed, &dashboard.Tasks.InProgress, &dashboard.Tasks.Overdue)

	// Get team stats
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE status = 'active'
	`).Scan(&dashboard.Team.TotalEmployees)

	database.DB.QueryRow(`
		SELECT COUNT(DISTINCT user_id) FROM vacations
		WHERE status = 'approved' AND type IN ('annual', 'unpaid', 'maternity')
		AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE
	`).Scan(&dashboard.Team.OnVacation)

	database.DB.QueryRow(`
		SELECT COUNT(DISTINCT user_id) FROM vacations
		WHERE status = 'approved' AND type = 'sick'
		AND start_date <= CURRENT_DATE AND end_date >= CURRENT_DATE
	`).Scan(&dashboard.Team.OnSickLeave)

	// Get upcoming events (next 7 days)
	eventRows, err := database.DB.Query(`
		SELECT id, title, start_date FROM events
		WHERE start_date >= NOW() AND start_date <= NOW() + INTERVAL '7 days'
		ORDER BY start_date ASC
		LIMIT 5
	`)
	if err == nil {
		defer eventRows.Close()
		for eventRows.Next() {
			var event models.DashboardEvent
			eventRows.Scan(&event.ID, &event.Title, &event.Date)
			dashboard.UpcomingEvents = append(dashboard.UpcomingEvents, event)
		}
	}
	if dashboard.UpcomingEvents == nil {
		dashboard.UpcomingEvents = []models.DashboardEvent{}
	}

	// Get recent activity
	activityRows, err := database.DB.Query(`
		SELECT a.id, a.type, a.description, a.user_id, CONCAT(u.first_name, ' ', u.last_name), a.created_at
		FROM activity_log a
		JOIN users u ON a.user_id = u.id
		ORDER BY a.created_at DESC
		LIMIT 10
	`)
	if err == nil {
		defer activityRows.Close()
		for activityRows.Next() {
			var activity models.ActivityLog
			activityRows.Scan(&activity.ID, &activity.Type, &activity.Description, &activity.User.ID, &activity.User.Name, &activity.Timestamp)
			dashboard.RecentActivity = append(dashboard.RecentActivity, activity)
		}
	}
	if dashboard.RecentActivity == nil {
		dashboard.RecentActivity = []models.ActivityLog{}
	}

	utils.RespondWithSuccess(c, http.StatusOK, dashboard)
}

func (h *AnalyticsHandler) GetProjectReport(c *gin.Context) {
	projectID := c.Param("id")

	// Get project basic info
	var project struct {
		ID        string
		Name      string
		StartDate sql.NullTime
		EndDate   sql.NullTime
	}

	err := database.DB.QueryRow(`
		SELECT id, name, start_date, end_date FROM projects WHERE id = $1
	`, projectID).Scan(&project.ID, &project.Name, &project.StartDate, &project.EndDate)

	if err != nil {
		utils.RespondWithError(c, http.StatusNotFound, "Project not found")
		return
	}

	report := models.ProjectReport{
		Project: models.ProjectReportBasic{
			ID:   project.ID,
			Name: project.Name,
		},
	}

	// Timeline
	if project.StartDate.Valid {
		report.Timeline.StartDate = project.StartDate.Time.Format("2006-01-02")
	}
	if project.EndDate.Valid {
		report.Timeline.EndDate = project.EndDate.Time.Format("2006-01-02")
		daysRemaining := int(time.Until(project.EndDate.Time).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}
		report.Timeline.DaysRemaining = daysRemaining
	}

	// Progress
	database.DB.QueryRow(`
		SELECT 
			COUNT(*),
			COALESCE(SUM(CASE WHEN status = 'done' THEN 1 ELSE 0 END), 0)
		FROM tasks WHERE project_id = $1
	`, projectID).Scan(&report.Progress.TasksTotal, &report.Progress.TasksCompleted)

	if report.Progress.TasksTotal > 0 {
		report.Progress.Overall = (report.Progress.TasksCompleted * 100) / report.Progress.TasksTotal
	}

	// Team
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM project_members WHERE project_id = $1
	`, projectID).Scan(&report.Team.Size)

	// Workload distribution
	workloadRows, err := database.DB.Query(`
		SELECT t.assignee_id, CONCAT(u.first_name, ' ', u.last_name), COUNT(*) as tasks_count
		FROM tasks t
		JOIN users u ON t.assignee_id = u.id
		WHERE t.project_id = $1 AND t.assignee_id IS NOT NULL
		GROUP BY t.assignee_id, u.first_name, u.last_name
		ORDER BY tasks_count DESC
	`, projectID)

	if err == nil {
		defer workloadRows.Close()
		for workloadRows.Next() {
			var item models.WorkloadDistributionItem
			workloadRows.Scan(&item.UserID, &item.Name, &item.TasksAssigned)
			report.Team.WorkloadDistribution = append(report.Team.WorkloadDistribution, item)
		}
	}
	if report.Team.WorkloadDistribution == nil {
		report.Team.WorkloadDistribution = []models.WorkloadDistributionItem{}
	}

	// Burndown data (simplified - just showing remaining tasks per day)
	burndownRows, err := database.DB.Query(`
		WITH RECURSIVE dates AS (
			SELECT COALESCE($2::date, CURRENT_DATE - INTERVAL '30 days') as date
			UNION ALL
			SELECT date + INTERVAL '1 day'
			FROM dates
			WHERE date < CURRENT_DATE
		)
		SELECT d.date, 
			(SELECT COUNT(*) FROM tasks t WHERE t.project_id = $1 AND t.status != 'done' AND t.created_at <= d.date)
		FROM dates d
		ORDER BY d.date
		LIMIT 30
	`, projectID, project.StartDate)

	if err == nil {
		defer burndownRows.Close()
		for burndownRows.Next() {
			var point models.BurndownPoint
			var date time.Time
			burndownRows.Scan(&date, &point.Remaining)
			point.Date = date.Format("2006-01-02")
			report.Burndown = append(report.Burndown, point)
		}
	}
	if report.Burndown == nil {
		report.Burndown = []models.BurndownPoint{}
	}

	utils.RespondWithSuccess(c, http.StatusOK, report)
}
