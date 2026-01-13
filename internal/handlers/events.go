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

type EventHandler struct{}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

type CreateEventRequest struct {
	Title          string                  `json:"title" binding:"required"`
	Description    string                  `json:"description"`
	Type           string                  `json:"type" binding:"required"`
	StartDate      string                  `json:"startDate" binding:"required"`
	EndDate        string                  `json:"endDate" binding:"required"`
	AllDay         bool                    `json:"allDay"`
	Location       string                  `json:"location"`
	ParticipantIDs []string                `json:"participantIds"`
	Reminders      []models.EventReminder  `json:"reminders"`
}

type UpdateEventRequest struct {
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	Type           string                  `json:"type"`
	StartDate      string                  `json:"startDate"`
	EndDate        string                  `json:"endDate"`
	AllDay         bool                    `json:"allDay"`
	Location       string                  `json:"location"`
	ParticipantIDs []string                `json:"participantIds"`
	Reminders      []models.EventReminder  `json:"reminders"`
}

type RespondEventRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")
	eventType := c.Query("type")

	// Build query
	query := `
		SELECT e.id, e.title, e.description, e.type, e.start_date, e.end_date, e.all_day, e.location, e.created_by, e.created_at
		FROM events e WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if startDateStr != "" {
		startDate, _ := time.Parse("2006-01-02", startDateStr)
		query += fmt.Sprintf(" AND e.start_date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDateStr != "" {
		endDate, _ := time.Parse("2006-01-02", endDateStr)
		query += fmt.Sprintf(" AND e.end_date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	if eventType != "" {
		query += fmt.Sprintf(" AND e.type = $%d", argIndex)
		args = append(args, eventType)
		argIndex++
	}

	query += " ORDER BY e.start_date ASC"

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get events")
		return
	}
	defer rows.Close()

	events := []models.EventListItem{}
	for rows.Next() {
		var event models.EventListItem
		var description, location, createdBy sql.NullString

		err := rows.Scan(&event.ID, &event.Title, &description, &event.Type, &event.StartDate, &event.EndDate, &event.AllDay, &location, &createdBy, &event.CreatedAt)
		if err != nil {
			continue
		}

		if description.Valid {
			event.Description = description.String
		}
		if location.Valid {
			event.Location = location.String
		}

		// Get participants
		event.Participants = h.getEventParticipants(event.ID)

		// Get creator
		if createdBy.Valid {
			var creatorName sql.NullString
			database.DB.QueryRow("SELECT CONCAT(first_name, ' ', last_name) FROM users WHERE id = $1", createdBy.String).Scan(&creatorName)
			if creatorName.Valid {
				event.CreatedBy = &models.TaskUser{ID: createdBy.String, Name: creatorName.String}
			}
		}

		events = append(events, event)
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"data": events})
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	startDate, err := time.Parse(time.RFC3339, req.StartDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid start date format")
		return
	}

	endDate, err := time.Parse(time.RFC3339, req.EndDate)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid end date format")
		return
	}

	var eventID string
	err = database.DB.QueryRow(`
		INSERT INTO events (title, description, type, start_date, end_date, all_day, location, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`, req.Title, nullString(req.Description), req.Type, startDate, endDate, req.AllDay, nullString(req.Location), userID).Scan(&eventID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create event")
		return
	}

	// Add participants
	for _, participantID := range req.ParticipantIDs {
		database.DB.Exec(`
			INSERT INTO event_participants (event_id, user_id, status)
			VALUES ($1, $2, 'pending')
			ON CONFLICT DO NOTHING
		`, eventID, participantID)

		// Send notification
		if participantID != userID {
			createNotification(participantID, "event_reminder", "Event Invitation", "You have been invited to: "+req.Title, "event", eventID)
		}
	}

	// Add reminders
	for _, reminder := range req.Reminders {
		database.DB.Exec(`
			INSERT INTO event_reminders (event_id, type, minutes_before)
			VALUES ($1, $2, $3)
		`, eventID, reminder.Type, reminder.MinutesBefore)
	}

	event := h.getEventByID(eventID)
	utils.RespondWithSuccess(c, http.StatusCreated, event)
}

func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID := c.Param("id")

	event := h.getEventByID(eventID)
	if event == nil {
		utils.RespondWithError(c, http.StatusNotFound, "Event not found")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, event)
}

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	eventID := c.Param("id")

	var req UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Build update query
	query := "UPDATE events SET updated_at = NOW()"
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
	if req.Type != "" {
		query += fmt.Sprintf(", type = $%d", argIndex)
		args = append(args, req.Type)
		argIndex++
	}
	if req.StartDate != "" {
		startDate, _ := time.Parse(time.RFC3339, req.StartDate)
		query += fmt.Sprintf(", start_date = $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}
	if req.EndDate != "" {
		endDate, _ := time.Parse(time.RFC3339, req.EndDate)
		query += fmt.Sprintf(", end_date = $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}
	if req.Location != "" {
		query += fmt.Sprintf(", location = $%d", argIndex)
		args = append(args, req.Location)
		argIndex++
	}

	query += fmt.Sprintf(", all_day = $%d", argIndex)
	args = append(args, req.AllDay)
	argIndex++

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, eventID)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update event")
		return
	}

	// Update participants if provided
	if req.ParticipantIDs != nil {
		database.DB.Exec("DELETE FROM event_participants WHERE event_id = $1", eventID)
		for _, participantID := range req.ParticipantIDs {
			database.DB.Exec(`
				INSERT INTO event_participants (event_id, user_id, status)
				VALUES ($1, $2, 'pending')
				ON CONFLICT DO NOTHING
			`, eventID, participantID)
		}
	}

	// Update reminders if provided
	if req.Reminders != nil {
		database.DB.Exec("DELETE FROM event_reminders WHERE event_id = $1", eventID)
		for _, reminder := range req.Reminders {
			database.DB.Exec(`
				INSERT INTO event_reminders (event_id, type, minutes_before)
				VALUES ($1, $2, $3)
			`, eventID, reminder.Type, reminder.MinutesBefore)
		}
	}

	event := h.getEventByID(eventID)
	utils.RespondWithSuccess(c, http.StatusOK, event)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")

	_, err := database.DB.Exec("DELETE FROM events WHERE id = $1", eventID)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func (h *EventHandler) RespondToEvent(c *gin.Context) {
	eventID := c.Param("id")
	userID := middleware.GetUserID(c)

	var req RespondEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Validate status
	if req.Status != "accepted" && req.Status != "declined" {
		utils.RespondWithError(c, http.StatusBadRequest, "Status must be 'accepted' or 'declined'")
		return
	}

	_, err := database.DB.Exec(`
		UPDATE event_participants SET status = $1 WHERE event_id = $2 AND user_id = $3
	`, req.Status, eventID, userID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to respond to event")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Response recorded"})
}

func (h *EventHandler) getEventByID(eventID string) interface{} {
	var event struct {
		ID          string
		Title       string
		Description sql.NullString
		Type        string
		StartDate   time.Time
		EndDate     time.Time
		AllDay      bool
		Location    sql.NullString
		CreatedBy   sql.NullString
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	err := database.DB.QueryRow(`
		SELECT id, title, description, type, start_date, end_date, all_day, location, created_by, created_at, updated_at
		FROM events WHERE id = $1
	`, eventID).Scan(&event.ID, &event.Title, &event.Description, &event.Type, &event.StartDate, &event.EndDate, &event.AllDay, &event.Location, &event.CreatedBy, &event.CreatedAt, &event.UpdatedAt)

	if err != nil {
		return nil
	}

	result := gin.H{
		"id":           event.ID,
		"title":        event.Title,
		"type":         event.Type,
		"startDate":    event.StartDate,
		"endDate":      event.EndDate,
		"allDay":       event.AllDay,
		"participants": h.getEventParticipants(eventID),
		"createdAt":    event.CreatedAt,
		"updatedAt":    event.UpdatedAt,
	}

	if event.Description.Valid {
		result["description"] = event.Description.String
	}
	if event.Location.Valid {
		result["location"] = event.Location.String
	}
	if event.CreatedBy.Valid {
		var creatorName sql.NullString
		database.DB.QueryRow("SELECT CONCAT(first_name, ' ', last_name) FROM users WHERE id = $1", event.CreatedBy.String).Scan(&creatorName)
		if creatorName.Valid {
			result["createdBy"] = gin.H{"id": event.CreatedBy.String, "name": creatorName.String}
		}
	}

	return result
}

func (h *EventHandler) getEventParticipants(eventID string) []models.EventParticipant {
	rows, err := database.DB.Query(`
		SELECT u.id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url, ep.status
		FROM event_participants ep
		JOIN users u ON ep.user_id = u.id
		WHERE ep.event_id = $1
	`, eventID)

	if err != nil {
		return []models.EventParticipant{}
	}
	defer rows.Close()

	participants := []models.EventParticipant{}
	for rows.Next() {
		var p models.EventParticipant
		var avatar sql.NullString
		rows.Scan(&p.ID, &p.Name, &avatar, &p.Status)
		if avatar.Valid {
			p.Avatar = avatar.String
		}
		participants = append(participants, p)
	}

	return participants
}
