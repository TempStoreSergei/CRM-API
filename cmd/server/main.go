package main

import (
	"log"
	"os"

	"github.com/TempStoreSergei/CRM-API/internal/config"
	"github.com/TempStoreSergei/CRM-API/internal/database"
	"github.com/TempStoreSergei/CRM-API/internal/handlers"
	"github.com/TempStoreSergei/CRM-API/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if exists
	godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Set gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create Gin router
	r := gin.New()

	// Middleware
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware(100, 200)) // 100 requests per second, burst of 200

	// Serve static files (uploads)
	r.Static("/uploads", "./uploads")

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg)
	userHandler := handlers.NewUserHandler()
	employeeHandler := handlers.NewEmployeeHandler()
	projectHandler := handlers.NewProjectHandler()
	taskHandler := handlers.NewTaskHandler()
	eventHandler := handlers.NewEventHandler()
	vacationHandler := handlers.NewVacationHandler()
	conversationHandler := handlers.NewConversationHandler()
	notificationHandler := handlers.NewNotificationHandler()
	analyticsHandler := handlers.NewAnalyticsHandler()
	integrationHandler := handlers.NewIntegrationHandler(cfg)
	fileHandler := handlers.NewFileHandler()
	wsHandler := handlers.NewWebSocketHandler(cfg)

	// API routes
	api := r.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/logout", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Logout)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			// Users
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetUsers)
				users.GET("/me", userHandler.GetCurrentUser)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.POST("/:id/avatar", userHandler.UploadAvatar)
			}

			// Employees
			employees := protected.Group("/employees")
			{
				employees.GET("", employeeHandler.GetEmployees)
				employees.GET("/:id/tasks", employeeHandler.GetEmployeeTasks)
			}

			// Projects
			projects := protected.Group("/projects")
			{
				projects.GET("", projectHandler.GetProjects)
				projects.POST("", projectHandler.CreateProject)
				projects.GET("/:id", projectHandler.GetProject)
				projects.PUT("/:id", projectHandler.UpdateProject)
				projects.DELETE("/:id", projectHandler.DeleteProject)
				projects.POST("/:id/members", projectHandler.AddMember)
				projects.DELETE("/:id/members/:userId", projectHandler.RemoveMember)
			}

			// Tasks
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", taskHandler.GetTasks)
				tasks.POST("", taskHandler.CreateTask)
				tasks.GET("/:id", taskHandler.GetTask)
				tasks.PUT("/:id", taskHandler.UpdateTask)
				tasks.PATCH("/:id/status", taskHandler.UpdateStatus)
				tasks.DELETE("/:id", taskHandler.DeleteTask)
				tasks.POST("/:id/comments", taskHandler.AddComment)
				tasks.GET("/:id/comments", taskHandler.GetComments)
				tasks.POST("/:id/files", taskHandler.AttachFile)
			}

			// Events
			events := protected.Group("/events")
			{
				events.GET("", eventHandler.GetEvents)
				events.POST("", eventHandler.CreateEvent)
				events.GET("/:id", eventHandler.GetEvent)
				events.PUT("/:id", eventHandler.UpdateEvent)
				events.DELETE("/:id", eventHandler.DeleteEvent)
				events.PATCH("/:id/respond", eventHandler.RespondToEvent)
			}

			// Vacations
			vacations := protected.Group("/vacations")
			{
				vacations.GET("", vacationHandler.GetVacations)
				vacations.POST("", vacationHandler.CreateVacation)
				vacations.GET("/balance", vacationHandler.GetBalance)
				vacations.GET("/:id", vacationHandler.GetVacation)
				vacations.PATCH("/:id/approve", vacationHandler.ApproveVacation)
				vacations.PATCH("/:id/reject", vacationHandler.RejectVacation)
				vacations.DELETE("/:id", vacationHandler.DeleteVacation)
			}

			// Conversations (Messenger)
			conversations := protected.Group("/conversations")
			{
				conversations.GET("", conversationHandler.GetConversations)
				conversations.POST("", conversationHandler.CreateConversation)
				conversations.GET("/:id/messages", conversationHandler.GetMessages)
				conversations.POST("/:id/messages", conversationHandler.SendMessage)
			}

			// Notifications
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.GetNotifications)
				notifications.PATCH("/:id/read", notificationHandler.MarkAsRead)
				notifications.PATCH("/read-all", notificationHandler.MarkAllAsRead)
				notifications.DELETE("/:id", notificationHandler.DeleteNotification)
				notifications.POST("/devices", notificationHandler.RegisterDevice)
			}

			// Analytics
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/dashboard", analyticsHandler.GetDashboard)
				analytics.GET("/projects/:id/report", analyticsHandler.GetProjectReport)
			}

			// Integrations
			integrations := protected.Group("/integrations")
			{
				integrations.GET("/google/auth", integrationHandler.GetGoogleAuthURL)
				integrations.POST("/google/callback", integrationHandler.GoogleCallback)
				integrations.POST("/google/sync", integrationHandler.GoogleSync)
				integrations.POST("/telegram/link", integrationHandler.LinkTelegram)
			}

			// Files
			files := protected.Group("/files")
			{
				files.POST("/upload", fileHandler.Upload)
				files.GET("/:id", fileHandler.GetFile)
				files.DELETE("/:id", fileHandler.DeleteFile)
			}
		}
	}

	// WebSocket endpoint
	r.GET("/ws", wsHandler.HandleConnection)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
