package handlers

import (
	"database/sql"
	"net/http"

	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/TempStoreSergei/CRM-API/internal/models"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type ConversationHandler struct{}

func NewConversationHandler() *ConversationHandler {
	return &ConversationHandler{}
}

type CreateConversationRequest struct {
	Type           string   `json:"type" binding:"required"`
	Name           string   `json:"name"`
	ParticipantIDs []string `json:"participantIds" binding:"required"`
}

type SendMessageRequest struct {
	Content       string   `json:"content" binding:"required"`
	Type          string   `json:"type"`
	AttachmentIDs []string `json:"attachmentIds"`
}

func (h *ConversationHandler) GetConversations(c *gin.Context) {
	userID := middleware.GetUserID(c)

	groups := []models.ConversationGroup{}
	directMessages := []models.ConversationDirect{}

	// Get group conversations
	groupRows, err := database.DB.Query(`
		SELECT c.id, c.name, c.avatar_url
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1 AND c.type = 'group'
		ORDER BY c.created_at DESC
	`, userID)

	if err == nil {
		defer groupRows.Close()
		for groupRows.Next() {
			var group models.ConversationGroup
			var name, avatar sql.NullString

			groupRows.Scan(&group.ID, &name, &avatar)

			group.Type = "group"
			if name.Valid {
				group.Name = name.String
			}
			if avatar.Valid {
				group.Avatar = avatar.String
			}

			// Get participants
			group.Participants = h.getConversationParticipants(group.ID)

			// Get last message
			group.LastMessage = h.getLastMessage(group.ID)

			// Get unread count
			group.UnreadCount = h.getUnreadCount(group.ID, userID)

			groups = append(groups, group)
		}
	}

	// Get direct message conversations
	dmRows, err := database.DB.Query(`
		SELECT c.id
		FROM conversations c
		JOIN conversation_participants cp ON c.id = cp.conversation_id
		WHERE cp.user_id = $1 AND c.type = 'direct'
		ORDER BY c.created_at DESC
	`, userID)

	if err == nil {
		defer dmRows.Close()
		for dmRows.Next() {
			var dm models.ConversationDirect
			dmRows.Scan(&dm.ID)

			dm.Type = "direct"

			// Get the other participant
			var participantID string
			var participantName string
			var participantAvatar sql.NullString

			database.DB.QueryRow(`
				SELECT u.id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url
				FROM conversation_participants cp
				JOIN users u ON cp.user_id = u.id
				WHERE cp.conversation_id = $1 AND cp.user_id != $2
				LIMIT 1
			`, dm.ID, userID).Scan(&participantID, &participantName, &participantAvatar)

			dm.Participant = models.DirectParticipant{
				ID:     participantID,
				Name:   participantName,
				Status: "offline", // In a real app, this would be tracked in real-time
			}
			if participantAvatar.Valid {
				dm.Participant.Avatar = participantAvatar.String
			}

			// Get last message
			dm.LastMessage = h.getLastMessage(dm.ID)

			// Get unread count
			dm.UnreadCount = h.getUnreadCount(dm.ID, userID)

			directMessages = append(directMessages, dm)
		}
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"groups":         groups,
		"directMessages": directMessages,
	})
}

func (h *ConversationHandler) CreateConversation(c *gin.Context) {
	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserID(c)

	// Validate type
	if req.Type != "group" && req.Type != "direct" {
		utils.RespondWithError(c, http.StatusBadRequest, "Type must be 'group' or 'direct'")
		return
	}

	// For direct messages, check if conversation already exists
	if req.Type == "direct" && len(req.ParticipantIDs) == 1 {
		otherUserID := req.ParticipantIDs[0]
		var existingID string
		err := database.DB.QueryRow(`
			SELECT c.id FROM conversations c
			JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
			JOIN conversation_participants cp2 ON c.id = cp2.conversation_id
			WHERE c.type = 'direct' AND cp1.user_id = $1 AND cp2.user_id = $2
		`, userID, otherUserID).Scan(&existingID)

		if err == nil {
			// Return existing conversation
			utils.RespondWithSuccess(c, http.StatusOK, gin.H{"id": existingID})
			return
		}
	}

	// Create conversation
	var conversationID string
	err := database.DB.QueryRow(`
		INSERT INTO conversations (type, name)
		VALUES ($1, $2)
		RETURNING id
	`, req.Type, nullString(req.Name)).Scan(&conversationID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create conversation")
		return
	}

	// Add creator as participant
	database.DB.Exec(`
		INSERT INTO conversation_participants (conversation_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, conversationID, userID)

	// Add other participants
	for _, participantID := range req.ParticipantIDs {
		database.DB.Exec(`
			INSERT INTO conversation_participants (conversation_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, conversationID, participantID)
	}

	utils.RespondWithSuccess(c, http.StatusCreated, gin.H{"id": conversationID})
}

func (h *ConversationHandler) GetMessages(c *gin.Context) {
	conversationID := c.Param("id")
	userID := middleware.GetUserID(c)

	page, limit := utils.GetPaginationParams(c)
	offset := (page - 1) * limit

	// Check if user is a participant
	var participantExists int
	database.DB.QueryRow(`
		SELECT 1 FROM conversation_participants
		WHERE conversation_id = $1 AND user_id = $2
	`, conversationID, userID).Scan(&participantExists)

	if participantExists == 0 {
		utils.RespondWithError(c, http.StatusForbidden, "You are not a participant in this conversation")
		return
	}

	rows, err := database.DB.Query(`
		SELECT m.id, m.content, m.type, m.sender_id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC
		LIMIT $2 OFFSET $3
	`, conversationID, limit+1, offset)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get messages")
		return
	}
	defer rows.Close()

	messages := []models.Message{}
	count := 0
	for rows.Next() {
		count++
		if count > limit {
			break
		}

		var msg models.Message
		var senderAvatar sql.NullString

		rows.Scan(&msg.ID, &msg.Content, &msg.Type, &msg.Sender.ID, &msg.Sender.Name, &senderAvatar, &msg.CreatedAt)

		msg.ConversationID = conversationID
		if senderAvatar.Valid {
			msg.Sender.Avatar = senderAvatar.String
		}

		// Get read by
		msg.ReadBy = h.getMessageReadBy(msg.ID)

		// Get attachments
		msg.Attachments = h.getMessageAttachments(msg.ID)

		messages = append(messages, msg)
	}

	// Update last read timestamp
	database.DB.Exec(`
		UPDATE conversation_participants SET last_read_at = NOW()
		WHERE conversation_id = $1 AND user_id = $2
	`, conversationID, userID)

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"data": messages,
		"meta": models.MessageMeta{HasMore: count > limit},
	})
}

func (h *ConversationHandler) SendMessage(c *gin.Context) {
	conversationID := c.Param("id")
	userID := middleware.GetUserID(c)

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if user is a participant
	var participantExists int
	database.DB.QueryRow(`
		SELECT 1 FROM conversation_participants
		WHERE conversation_id = $1 AND user_id = $2
	`, conversationID, userID).Scan(&participantExists)

	if participantExists == 0 {
		utils.RespondWithError(c, http.StatusForbidden, "You are not a participant in this conversation")
		return
	}

	msgType := req.Type
	if msgType == "" {
		msgType = "text"
	}

	var messageID string
	err := database.DB.QueryRow(`
		INSERT INTO messages (conversation_id, sender_id, content, type)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, conversationID, userID, req.Content, msgType).Scan(&messageID)

	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to send message")
		return
	}

	// Link attachments if provided
	for _, attachmentID := range req.AttachmentIDs {
		database.DB.Exec(`
			UPDATE files SET entity_type = 'message', entity_id = $1
			WHERE id = $2
		`, messageID, attachmentID)
	}

	// Mark as read by sender
	database.DB.Exec(`
		INSERT INTO message_read_status (message_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, messageID, userID)

	// Get the created message
	var msg models.Message
	var senderName, senderAvatar sql.NullString

	database.DB.QueryRow(`
		SELECT m.id, m.content, m.type, m.sender_id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = $1
	`, messageID).Scan(&msg.ID, &msg.Content, &msg.Type, &msg.Sender.ID, &senderName, &senderAvatar, &msg.CreatedAt)

	msg.ConversationID = conversationID
	if senderName.Valid {
		msg.Sender.Name = senderName.String
	}
	if senderAvatar.Valid {
		msg.Sender.Avatar = senderAvatar.String
	}
	msg.Attachments = h.getMessageAttachments(messageID)
	msg.ReadBy = []string{userID}

	utils.RespondWithSuccess(c, http.StatusCreated, msg)
}

func (h *ConversationHandler) getConversationParticipants(conversationID string) []models.ConversationParticipant {
	rows, err := database.DB.Query(`
		SELECT u.id, CONCAT(u.first_name, ' ', u.last_name), u.avatar_url
		FROM conversation_participants cp
		JOIN users u ON cp.user_id = u.id
		WHERE cp.conversation_id = $1
	`, conversationID)

	if err != nil {
		return []models.ConversationParticipant{}
	}
	defer rows.Close()

	participants := []models.ConversationParticipant{}
	for rows.Next() {
		var p models.ConversationParticipant
		var avatar sql.NullString
		rows.Scan(&p.ID, &p.Name, &avatar)
		if avatar.Valid {
			p.Avatar = avatar.String
		}
		participants = append(participants, p)
	}

	return participants
}

func (h *ConversationHandler) getLastMessage(conversationID string) *models.ConversationLastMessage {
	var msg models.ConversationLastMessage
	var senderID, senderName sql.NullString

	err := database.DB.QueryRow(`
		SELECT m.id, m.content, m.sender_id, CONCAT(u.first_name, ' ', u.last_name), m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.conversation_id = $1
		ORDER BY m.created_at DESC
		LIMIT 1
	`, conversationID).Scan(&msg.ID, &msg.Content, &senderID, &senderName, &msg.Timestamp)

	if err != nil {
		return nil
	}

	if senderID.Valid {
		msg.Sender = &models.TaskUser{ID: senderID.String, Name: senderName.String}
	}

	return &msg
}

func (h *ConversationHandler) getUnreadCount(conversationID, userID string) int {
	var count int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM messages m
		WHERE m.conversation_id = $1
		AND m.sender_id != $2
		AND m.created_at > COALESCE(
			(SELECT last_read_at FROM conversation_participants
			WHERE conversation_id = $1 AND user_id = $2),
			'1970-01-01'
		)
	`, conversationID, userID).Scan(&count)
	return count
}

func (h *ConversationHandler) getMessageReadBy(messageID string) []string {
	rows, err := database.DB.Query("SELECT user_id FROM message_read_status WHERE message_id = $1", messageID)
	if err != nil {
		return []string{}
	}
	defer rows.Close()

	userIDs := []string{}
	for rows.Next() {
		var userID string
		rows.Scan(&userID)
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func (h *ConversationHandler) getMessageAttachments(messageID string) []models.MessageAttachment {
	rows, err := database.DB.Query(`
		SELECT id, name, url, type, size
		FROM files WHERE entity_type = 'message' AND entity_id = $1
	`, messageID)

	if err != nil {
		return []models.MessageAttachment{}
	}
	defer rows.Close()

	attachments := []models.MessageAttachment{}
	for rows.Next() {
		var a models.MessageAttachment
		var fileType sql.NullString
		rows.Scan(&a.ID, &a.Name, &a.URL, &fileType, &a.Size)
		if fileType.Valid {
			a.Type = fileType.String
		}
		attachments = append(attachments, a)
	}

	return attachments
}
