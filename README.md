# CRM-API

A complete CRM (Customer Relationship Management) system backend built with Go.

## Features

- **Authentication & Authorization**: JWT-based authentication with refresh tokens, role-based access control
- **User Management**: User registration, profile management, avatar uploads
- **Employee Management**: Track employees, workload, and status (active, vacation, sick leave)
- **Project Management**: Create, update, and manage projects with team members
- **Task Management**: CRUD operations, status tracking, comments, file attachments
- **Event Management**: Calendar events with participants and reminders
- **Vacation Management**: Request, approve/reject vacations, track vacation balance
- **Messenger**: Real-time chat with WebSocket support, group and direct messages
- **Notifications**: Push notification support, in-app notifications
- **Analytics**: Dashboard with project statistics, burndown charts
- **Integrations**: Google Calendar, Telegram (optional)
- **File Management**: File uploads with metadata tracking

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Real-time**: WebSocket (Gorilla)
- **Authentication**: JWT
- **Containerization**: Docker

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 15+
- Redis (optional, for caching)
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/TempStoreSergei/CRM-API.git
cd CRM-API
```

2. Copy environment variables:
```bash
cp .env.example .env
```

3. Update `.env` with your configuration

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run cmd/server/main.go
```

### Using Docker

1. Build and run with Docker Compose:
```bash
docker compose up -d
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Register new user
- `POST /api/auth/login` - Login
- `POST /api/auth/refresh` - Refresh tokens
- `POST /api/auth/logout` - Logout
- `POST /api/auth/forgot-password` - Password reset request

### Users
- `GET /api/users` - List users
- `GET /api/users/me` - Get current user
- `GET /api/users/:id` - Get user by ID
- `PUT /api/users/:id` - Update user
- `POST /api/users/:id/avatar` - Upload avatar

### Employees
- `GET /api/employees` - List employees with workload
- `GET /api/employees/:id/tasks` - Get employee tasks

### Projects
- `GET /api/projects` - List projects
- `POST /api/projects` - Create project
- `GET /api/projects/:id` - Get project
- `PUT /api/projects/:id` - Update project
- `DELETE /api/projects/:id` - Delete project
- `POST /api/projects/:id/members` - Add member
- `DELETE /api/projects/:id/members/:userId` - Remove member

### Tasks
- `GET /api/tasks` - List tasks
- `POST /api/tasks` - Create task
- `GET /api/tasks/:id` - Get task
- `PUT /api/tasks/:id` - Update task
- `PATCH /api/tasks/:id/status` - Update status
- `DELETE /api/tasks/:id` - Delete task
- `POST /api/tasks/:id/comments` - Add comment
- `GET /api/tasks/:id/comments` - Get comments
- `POST /api/tasks/:id/files` - Attach file

### Events
- `GET /api/events` - List events
- `POST /api/events` - Create event
- `GET /api/events/:id` - Get event
- `PUT /api/events/:id` - Update event
- `DELETE /api/events/:id` - Delete event
- `PATCH /api/events/:id/respond` - Respond to event

### Vacations
- `GET /api/vacations` - List vacation requests
- `POST /api/vacations` - Create vacation request
- `GET /api/vacations/balance` - Get vacation balance
- `GET /api/vacations/:id` - Get vacation request
- `PATCH /api/vacations/:id/approve` - Approve request
- `PATCH /api/vacations/:id/reject` - Reject request
- `DELETE /api/vacations/:id` - Cancel request

### Conversations (Messenger)
- `GET /api/conversations` - List conversations
- `POST /api/conversations` - Create conversation
- `GET /api/conversations/:id/messages` - Get messages
- `POST /api/conversations/:id/messages` - Send message

### Notifications
- `GET /api/notifications` - List notifications
- `PATCH /api/notifications/:id/read` - Mark as read
- `PATCH /api/notifications/read-all` - Mark all as read
- `DELETE /api/notifications/:id` - Delete notification
- `POST /api/notifications/devices` - Register device

### Analytics
- `GET /api/analytics/dashboard` - Get dashboard data
- `GET /api/analytics/projects/:id/report` - Get project report

### Integrations
- `GET /api/integrations/google/auth` - Get Google auth URL
- `POST /api/integrations/google/callback` - Google OAuth callback
- `POST /api/integrations/google/sync` - Sync Google Calendar
- `POST /api/integrations/telegram/link` - Link Telegram

### Files
- `POST /api/files/upload` - Upload file
- `GET /api/files/:id` - Get file info
- `DELETE /api/files/:id` - Delete file

### WebSocket
- `GET /ws?token=<accessToken>` - WebSocket connection for real-time messaging

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection URL | Required |
| `REDIS_URL` | Redis connection URL | Optional |
| `JWT_SECRET` | JWT signing secret | Required |
| `JWT_REFRESH_SECRET` | Refresh token secret | Required |
| `JWT_EXPIRES_IN` | Access token expiry | `15m` |
| `JWT_REFRESH_EXPIRES` | Refresh token expiry | `168h` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | Optional |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | Optional |
| `ENVIRONMENT` | Environment (development/production) | `development` |

## Database Schema

The application automatically runs migrations on startup. See [BACKEND_SPEC.md](BACKEND_SPEC.md) for the complete database schema.

## Security

- JWT authentication with short-lived access tokens
- Refresh token rotation
- Password hashing with bcrypt
- Rate limiting
- CORS middleware
- Role-based access control (admin, manager, employee)

## License

MIT