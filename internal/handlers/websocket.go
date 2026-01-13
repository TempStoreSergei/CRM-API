package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/TempStoreSergei/CRM-API/internal/config"
	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type WebSocketHandler struct {
	cfg     *config.Config
	clients map[string]*Client
	mu      sync.RWMutex
}

type Client struct {
	ID             string
	UserID         string
	Conn           *websocket.Conn
	Send           chan []byte
	ConversationID string
}

type WSMessage struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type MessageSendData struct {
	ConversationID string `json:"conversationId"`
	Content        string `json:"content"`
	Type           string `json:"type"`
}

type MarkReadData struct {
	ConversationID string   `json:"conversationId"`
	MessageIDs     []string `json:"messageIds"`
}

type TypingData struct {
	ConversationID string `json:"conversationId"`
	IsTyping       bool   `json:"isTyping"`
}

func NewWebSocketHandler(cfg *config.Config) *WebSocketHandler {
	return &WebSocketHandler{
		cfg:     cfg,
		clients: make(map[string]*Client),
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Validate token
	claims, err := utils.ValidateAccessToken(token, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:     claims.UserID + "-" + time.Now().Format("20060102150405"),
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.registerClient(client)

	// Notify others that user is online
	h.broadcastUserStatus(claims.UserID, "online")

	go h.writePump(client)
	h.readPump(client)
}

func (h *WebSocketHandler) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.ID] = client
}

func (h *WebSocketHandler) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients, client.ID)
		close(client.Send)
	}
}

func (h *WebSocketHandler) readPump(client *Client) {
	defer func() {
		h.broadcastUserStatus(client.UserID, "offline")
		h.unregisterClient(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(512 * 1024)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		h.handleMessage(client, message)
	}
}

func (h *WebSocketHandler) writePump(client *Client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			w.Close()

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *WebSocketHandler) handleMessage(client *Client, message []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Invalid message format: %v", err)
		return
	}

	switch wsMsg.Event {
	case "message.send":
		h.handleMessageSend(client, wsMsg.Data)
	case "message.markRead":
		h.handleMarkRead(client, wsMsg.Data)
	case "typing":
		h.handleTyping(client, wsMsg.Data)
	}
}

func (h *WebSocketHandler) handleMessageSend(client *Client, data json.RawMessage) {
	var msgData MessageSendData
	if err := json.Unmarshal(data, &msgData); err != nil {
		return
	}

	msgType := msgData.Type
	if msgType == "" {
		msgType = "text"
	}

	// Save message to database
	var messageID string
	err := database.DB.QueryRow(`
		INSERT INTO messages (conversation_id, sender_id, content, type)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, msgData.ConversationID, client.UserID, msgData.Content, msgType).Scan(&messageID)

	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	// Get sender info
	var senderName, senderAvatar string
	database.DB.QueryRow(`
		SELECT CONCAT(first_name, ' ', last_name), COALESCE(avatar_url, '')
		FROM users WHERE id = $1
	`, client.UserID).Scan(&senderName, &senderAvatar)

	// Broadcast to all participants
	response := gin.H{
		"event": "message.new",
		"data": gin.H{
			"conversationId": msgData.ConversationID,
			"message": gin.H{
				"id":      messageID,
				"content": msgData.Content,
				"type":    msgType,
				"sender": gin.H{
					"id":     client.UserID,
					"name":   senderName,
					"avatar": senderAvatar,
				},
				"createdAt": time.Now(),
			},
		},
	}

	responseBytes, _ := json.Marshal(response)
	h.broadcastToConversation(msgData.ConversationID, responseBytes)
}

func (h *WebSocketHandler) handleMarkRead(client *Client, data json.RawMessage) {
	var readData MarkReadData
	if err := json.Unmarshal(data, &readData); err != nil {
		return
	}

	// Update read status
	for _, messageID := range readData.MessageIDs {
		database.DB.Exec(`
			INSERT INTO message_read_status (message_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, messageID, client.UserID)
	}

	// Broadcast read status
	response := gin.H{
		"event": "message.read",
		"data": gin.H{
			"conversationId": readData.ConversationID,
			"messageIds":     readData.MessageIDs,
			"userId":         client.UserID,
		},
	}

	responseBytes, _ := json.Marshal(response)
	h.broadcastToConversation(readData.ConversationID, responseBytes)
}

func (h *WebSocketHandler) handleTyping(client *Client, data json.RawMessage) {
	var typingData TypingData
	if err := json.Unmarshal(data, &typingData); err != nil {
		return
	}

	// Get user name
	var userName string
	database.DB.QueryRow("SELECT CONCAT(first_name, ' ', last_name) FROM users WHERE id = $1", client.UserID).Scan(&userName)

	eventType := "typing.stop"
	if typingData.IsTyping {
		eventType = "typing.start"
	}

	response := gin.H{
		"event": eventType,
		"data": gin.H{
			"conversationId": typingData.ConversationID,
			"userId":         client.UserID,
			"userName":       userName,
		},
	}

	responseBytes, _ := json.Marshal(response)
	h.broadcastToConversation(typingData.ConversationID, responseBytes)
}

func (h *WebSocketHandler) broadcastToConversation(conversationID string, message []byte) {
	// Get all participants of the conversation
	rows, err := database.DB.Query(`
		SELECT user_id FROM conversation_participants WHERE conversation_id = $1
	`, conversationID)
	if err != nil {
		return
	}
	defer rows.Close()

	participantIDs := []string{}
	for rows.Next() {
		var userID string
		rows.Scan(&userID)
		participantIDs = append(participantIDs, userID)
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		for _, participantID := range participantIDs {
			if client.UserID == participantID {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client.ID)
				}
				break
			}
		}
	}
}

func (h *WebSocketHandler) broadcastUserStatus(userID, status string) {
	response := gin.H{
		"event": "user.status",
		"data": gin.H{
			"userId": userID,
			"status": status,
		},
	}

	responseBytes, _ := json.Marshal(response)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.UserID != userID {
			select {
			case client.Send <- responseBytes:
			default:
			}
		}
	}
}
