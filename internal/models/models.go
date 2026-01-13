package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	Password   string         `json:"-"`
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	Avatar     sql.NullString `json:"avatar,omitempty"`
	Phone      sql.NullString `json:"phone,omitempty"`
	Position   sql.NullString `json:"position,omitempty"`
	Department sql.NullString `json:"department,omitempty"`
	Level      sql.NullString `json:"level,omitempty"`
	Location   sql.NullString `json:"location,omitempty"`
	Birthday   sql.NullTime   `json:"birthday,omitempty"`
	Skype      sql.NullString `json:"skype,omitempty"`
	Role       string         `json:"role"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Skills     []string       `json:"skills,omitempty"`
}

type UserResponse struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Avatar     string    `json:"avatar,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Position   string    `json:"position,omitempty"`
	Department string    `json:"department,omitempty"`
	Level      string    `json:"level,omitempty"`
	Location   string    `json:"location,omitempty"`
	Birthday   string    `json:"birthday,omitempty"`
	Skype      string    `json:"skype,omitempty"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	Skills     []string  `json:"skills,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type UserListItem struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Avatar     string    `json:"avatar,omitempty"`
	Position   string    `json:"position,omitempty"`
	Department string    `json:"department,omitempty"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

type UserProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type UserTeamMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

type Project struct {
	ID              string         `json:"id"`
	Code            string         `json:"code"`
	Name            string         `json:"name"`
	Description     sql.NullString `json:"description,omitempty"`
	Status          string         `json:"status"`
	StartDate       sql.NullTime   `json:"startDate,omitempty"`
	EndDate         sql.NullTime   `json:"endDate,omitempty"`
	BudgetAllocated sql.NullFloat64 `json:"budgetAllocated,omitempty"`
	BudgetSpent     float64        `json:"budgetSpent"`
	BudgetCurrency  string         `json:"budgetCurrency"`
	CreatedBy       sql.NullString `json:"createdBy,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

type ProjectListItem struct {
	ID          string              `json:"id"`
	Code        string              `json:"code"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Status      string              `json:"status"`
	Progress    int                 `json:"progress"`
	StartDate   string              `json:"startDate,omitempty"`
	EndDate     string              `json:"endDate,omitempty"`
	Team        []ProjectTeamMember `json:"team"`
	TasksCount  TasksCount          `json:"tasksCount"`
	CreatedAt   time.Time           `json:"createdAt"`
}

type ProjectTeamMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
	Role     string `json:"role,omitempty"`
	Position string `json:"position,omitempty"`
}

type TasksCount struct {
	Total      int `json:"total"`
	Completed  int `json:"completed"`
	InProgress int `json:"inProgress"`
	Backlog    int `json:"backlog"`
}

type ProjectBudget struct {
	Allocated float64 `json:"allocated"`
	Spent     float64 `json:"spent"`
	Currency  string  `json:"currency"`
}

type Task struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Description    sql.NullString `json:"description,omitempty"`
	Status         string         `json:"status"`
	Priority       string         `json:"priority"`
	ProjectID      sql.NullString `json:"projectId,omitempty"`
	AssigneeID     sql.NullString `json:"assigneeId,omitempty"`
	ReporterID     sql.NullString `json:"reporterId,omitempty"`
	DueDate        sql.NullTime   `json:"dueDate,omitempty"`
	EstimatedHours sql.NullFloat64 `json:"estimatedHours,omitempty"`
	LoggedHours    float64        `json:"loggedHours"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	Tags           []string       `json:"tags,omitempty"`
}

type TaskListItem struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description,omitempty"`
	Status         string         `json:"status"`
	Priority       string         `json:"priority"`
	Project        *TaskProject   `json:"project,omitempty"`
	Assignee       *TaskUser      `json:"assignee,omitempty"`
	Reporter       *TaskUser      `json:"reporter,omitempty"`
	DueDate        string         `json:"dueDate,omitempty"`
	EstimatedHours float64        `json:"estimatedHours,omitempty"`
	LoggedHours    float64        `json:"loggedHours"`
	Tags           []string       `json:"tags,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type TaskProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type TaskUser struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type TaskComment struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"taskId"`
	User      TaskUser  `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Event struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description,omitempty"`
	Type        string         `json:"type"`
	StartDate   time.Time      `json:"startDate"`
	EndDate     time.Time      `json:"endDate"`
	AllDay      bool           `json:"allDay"`
	Location    sql.NullString `json:"location,omitempty"`
	CreatedBy   sql.NullString `json:"createdBy,omitempty"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
}

type EventListItem struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Description  string             `json:"description,omitempty"`
	Type         string             `json:"type"`
	StartDate    time.Time          `json:"startDate"`
	EndDate      time.Time          `json:"endDate"`
	AllDay       bool               `json:"allDay"`
	Location     string             `json:"location,omitempty"`
	Participants []EventParticipant `json:"participants"`
	CreatedBy    *TaskUser          `json:"createdBy,omitempty"`
	CreatedAt    time.Time          `json:"createdAt"`
}

type EventParticipant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Status string `json:"status"`
}

type EventReminder struct {
	Type          string `json:"type"`
	MinutesBefore int    `json:"minutesBefore"`
}

type Vacation struct {
	ID              string         `json:"id"`
	UserID          string         `json:"userId"`
	Type            string         `json:"type"`
	StartDate       time.Time      `json:"startDate"`
	EndDate         time.Time      `json:"endDate"`
	TotalDays       int            `json:"totalDays"`
	Status          string         `json:"status"`
	Reason          sql.NullString `json:"reason,omitempty"`
	RejectionReason sql.NullString `json:"rejectionReason,omitempty"`
	ApprovedBy      sql.NullString `json:"approvedBy,omitempty"`
	ApprovedAt      sql.NullTime   `json:"approvedAt,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

type VacationListItem struct {
	ID         string           `json:"id"`
	User       VacationUser     `json:"user"`
	Type       string           `json:"type"`
	StartDate  string           `json:"startDate"`
	EndDate    string           `json:"endDate"`
	TotalDays  int              `json:"totalDays"`
	Status     string           `json:"status"`
	Reason     string           `json:"reason,omitempty"`
	ApprovedBy *VacationApprover `json:"approvedBy,omitempty"`
	CreatedAt  time.Time        `json:"createdAt"`
}

type VacationUser struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
	Position string `json:"position,omitempty"`
}

type VacationApprover struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VacationBalance struct {
	AnnualTotal     int `json:"annualTotal"`
	AnnualUsed      int `json:"annualUsed"`
	AnnualRemaining int `json:"annualRemaining"`
	SickTotal       int `json:"sickTotal"`
	SickUsed        int `json:"sickUsed"`
	SickRemaining   int `json:"sickRemaining"`
}

type Conversation struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name,omitempty"`
	AvatarURL string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type ConversationGroup struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Avatar       string                      `json:"avatar,omitempty"`
	Type         string                      `json:"type"`
	Participants []ConversationParticipant   `json:"participants"`
	LastMessage  *ConversationLastMessage    `json:"lastMessage,omitempty"`
	UnreadCount  int                         `json:"unreadCount"`
}

type ConversationDirect struct {
	ID          string                      `json:"id"`
	Participant DirectParticipant           `json:"participant"`
	Type        string                      `json:"type"`
	LastMessage *ConversationLastMessage    `json:"lastMessage,omitempty"`
	UnreadCount int                         `json:"unreadCount"`
}

type ConversationParticipant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type DirectParticipant struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Status string `json:"status"`
}

type ConversationLastMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Sender    *TaskUser `json:"sender,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	ID             string              `json:"id"`
	ConversationID string              `json:"conversationId"`
	Content        string              `json:"content"`
	Type           string              `json:"type"`
	Sender         TaskUser            `json:"sender"`
	Attachments    []MessageAttachment `json:"attachments,omitempty"`
	ReadBy         []string            `json:"readBy,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
}

type MessageAttachment struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"userId"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Read      bool                   `json:"read"`
	CreatedAt time.Time              `json:"createdAt"`
}

type File struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Type       string    `json:"type,omitempty"`
	Size       int       `json:"size"`
	UploadedBy string    `json:"uploadedBy,omitempty"`
	EntityType string    `json:"entityType,omitempty"`
	EntityID   string    `json:"entityId,omitempty"`
	UploadedAt time.Time `json:"uploadedAt"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"token"`
	UserAgent string    `json:"userAgent,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type Employee struct {
	ID         string           `json:"id"`
	UserID     string           `json:"userId"`
	FirstName  string           `json:"firstName"`
	LastName   string           `json:"lastName"`
	Email      string           `json:"email"`
	Avatar     string           `json:"avatar,omitempty"`
	Position   string           `json:"position,omitempty"`
	Level      string           `json:"level,omitempty"`
	Department string           `json:"department,omitempty"`
	Status     string           `json:"status"`
	Workload   EmployeeWorkload `json:"workload"`
}

type EmployeeWorkload struct {
	BacklogTasks    int `json:"backlogTasks"`
	InProgressTasks int `json:"inProgressTasks"`
	InReviewTasks   int `json:"inReviewTasks"`
}

type ActivityLog struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	User        TaskUser  `json:"user"`
	Timestamp   time.Time `json:"timestamp"`
}

type DashboardAnalytics struct {
	Projects       DashboardProjects    `json:"projects"`
	Tasks          DashboardTasks       `json:"tasks"`
	Team           DashboardTeam        `json:"team"`
	UpcomingEvents []DashboardEvent     `json:"upcomingEvents"`
	RecentActivity []ActivityLog        `json:"recentActivity"`
}

type DashboardProjects struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Completed int `json:"completed"`
}

type DashboardTasks struct {
	Total      int `json:"total"`
	Completed  int `json:"completed"`
	InProgress int `json:"inProgress"`
	Overdue    int `json:"overdue"`
}

type DashboardTeam struct {
	TotalEmployees int `json:"totalEmployees"`
	OnVacation     int `json:"onVacation"`
	OnSickLeave    int `json:"onSickLeave"`
}

type DashboardEvent struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
}

type ProjectReport struct {
	Project  ProjectReportBasic    `json:"project"`
	Timeline ProjectReportTimeline `json:"timeline"`
	Progress ProjectReportProgress `json:"progress"`
	Team     ProjectReportTeam     `json:"team"`
	Burndown []BurndownPoint       `json:"burndown"`
}

type ProjectReportBasic struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectReportTimeline struct {
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	DaysRemaining int    `json:"daysRemaining"`
}

type ProjectReportProgress struct {
	Overall        int `json:"overall"`
	TasksCompleted int `json:"tasksCompleted"`
	TasksTotal     int `json:"tasksTotal"`
}

type ProjectReportTeam struct {
	Size                 int                       `json:"size"`
	WorkloadDistribution []WorkloadDistributionItem `json:"workloadDistribution"`
}

type WorkloadDistributionItem struct {
	UserID        string `json:"userId"`
	Name          string `json:"name"`
	TasksAssigned int    `json:"tasksAssigned"`
}

type BurndownPoint struct {
	Date      string `json:"date"`
	Remaining int    `json:"remaining"`
}

// Pagination
type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"totalPages"`
}

type NotificationMeta struct {
	Total       int `json:"total"`
	UnreadCount int `json:"unreadCount"`
	Page        int `json:"page"`
	Limit       int `json:"limit"`
}

type MessageMeta struct {
	HasMore bool `json:"hasMore"`
}
