package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(databaseURL string) error {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func RunMigrations() error {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`,

		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			first_name VARCHAR(100) NOT NULL,
			last_name VARCHAR(100) NOT NULL,
			avatar_url VARCHAR(500),
			phone VARCHAR(20),
			position VARCHAR(100),
			department VARCHAR(100),
			level VARCHAR(20),
			location VARCHAR(200),
			birthday DATE,
			skype VARCHAR(100),
			role VARCHAR(20) DEFAULT 'employee',
			status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS user_skills (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			skill VARCHAR(100) NOT NULL,
			UNIQUE(user_id, skill)
		)`,

		`CREATE TABLE IF NOT EXISTS projects (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			code VARCHAR(20) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(20) DEFAULT 'active',
			start_date DATE,
			end_date DATE,
			budget_allocated DECIMAL(15,2),
			budget_spent DECIMAL(15,2) DEFAULT 0,
			budget_currency VARCHAR(10) DEFAULT 'RUB',
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS project_members (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			role VARCHAR(50),
			joined_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(project_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS tasks (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			status VARCHAR(20) DEFAULT 'backlog',
			priority VARCHAR(20) DEFAULT 'medium',
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			assignee_id UUID REFERENCES users(id),
			reporter_id UUID REFERENCES users(id),
			due_date TIMESTAMP,
			estimated_hours DECIMAL(5,2),
			logged_hours DECIMAL(5,2) DEFAULT 0,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS task_tags (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
			tag VARCHAR(50) NOT NULL,
			UNIQUE(task_id, tag)
		)`,

		`CREATE TABLE IF NOT EXISTS task_comments (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id),
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS events (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			type VARCHAR(20) NOT NULL,
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			all_day BOOLEAN DEFAULT FALSE,
			location VARCHAR(255),
			created_by UUID REFERENCES users(id),
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS event_participants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID REFERENCES events(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			status VARCHAR(20) DEFAULT 'pending',
			UNIQUE(event_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS event_reminders (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			event_id UUID REFERENCES events(id) ON DELETE CASCADE,
			type VARCHAR(20) NOT NULL,
			minutes_before INTEGER NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS vacations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			type VARCHAR(20) NOT NULL,
			start_date DATE NOT NULL,
			end_date DATE NOT NULL,
			total_days INTEGER NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			reason TEXT,
			rejection_reason TEXT,
			approved_by UUID REFERENCES users(id),
			approved_at TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS vacation_balances (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			year INTEGER NOT NULL,
			annual_total INTEGER DEFAULT 28,
			annual_used INTEGER DEFAULT 0,
			sick_total INTEGER DEFAULT 10,
			sick_used INTEGER DEFAULT 0,
			UNIQUE(user_id, year)
		)`,

		`CREATE TABLE IF NOT EXISTS conversations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type VARCHAR(20) NOT NULL,
			name VARCHAR(255),
			avatar_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS conversation_participants (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			last_read_at TIMESTAMP,
			UNIQUE(conversation_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
			sender_id UUID REFERENCES users(id),
			content TEXT NOT NULL,
			type VARCHAR(20) DEFAULT 'text',
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS message_read_status (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			read_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(message_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			type VARCHAR(50) NOT NULL,
			title VARCHAR(255) NOT NULL,
			message TEXT,
			data JSONB,
			read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS notification_devices (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(500) NOT NULL,
			platform VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id, token)
		)`,

		`CREATE TABLE IF NOT EXISTS files (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			url VARCHAR(500) NOT NULL,
			type VARCHAR(100),
			size INTEGER,
			uploaded_by UUID REFERENCES users(id),
			entity_type VARCHAR(50),
			entity_id UUID,
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS refresh_tokens (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			token VARCHAR(500) NOT NULL,
			user_agent TEXT,
			expires_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		`CREATE TABLE IF NOT EXISTS google_integrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			access_token TEXT,
			refresh_token TEXT,
			token_expiry TIMESTAMP,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS telegram_integrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			telegram_user_id VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id),
			UNIQUE(telegram_user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS activity_log (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id),
			type VARCHAR(50) NOT NULL,
			description TEXT,
			entity_type VARCHAR(50),
			entity_id UUID,
			created_at TIMESTAMP DEFAULT NOW()
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_users_department ON users(department)`,
		`CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_code ON projects(code)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_assignee_id ON tasks(assignee_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status)`,
		`CREATE INDEX IF NOT EXISTS idx_events_start_date ON events(start_date)`,
		`CREATE INDEX IF NOT EXISTS idx_events_type ON events(type)`,
		`CREATE INDEX IF NOT EXISTS idx_vacations_user_id ON vacations(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_vacations_status ON vacations(status)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read)`,
		`CREATE INDEX IF NOT EXISTS idx_files_entity ON files(entity_type, entity_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token)`,
		`CREATE INDEX IF NOT EXISTS idx_activity_log_user_id ON activity_log(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activity_log_created_at ON activity_log(created_at)`,
	}

	for _, migration := range migrations {
		if _, err := DB.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration: %w\nSQL: %s", err, migration)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
