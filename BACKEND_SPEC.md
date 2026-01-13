# Техническое задание на разработку Backend для CRM системы

## 1. Общее описание

### 1.1 Назначение системы
CRM система предназначена для управления сотрудниками, проектами, задачами, отпусками и внутренними коммуникациями компании.

### 1.2 Технологический стек (рекомендуемый)
- **Язык программирования**: Node.js (TypeScript) / Go / Python
- **Фреймворк**: NestJS / Express / FastAPI
- **База данных**: PostgreSQL
- **Кэширование**: Redis
- **Очереди сообщений**: RabbitMQ / Bull
- **Real-time коммуникации**: WebSocket (Socket.io)
- **Документация API**: OpenAPI/Swagger
- **Контейнеризация**: Docker

---

## 2. Аутентификация и Авторизация

### 2.1 Endpoints авторизации

#### POST /api/auth/signup
Регистрация нового пользователя.

**Request Body:**
```json
{
  "email": "string",
  "password": "string",
  "firstName": "string",
  "lastName": "string",
  "agent": "string"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "createdAt": "timestamp"
}
```

#### POST /api/auth/login
Авторизация пользователя.

**Request Body:**
```json
{
  "email": "string",
  "password": "string",
  "agent": "string"
}
```

**Response (200 OK):**
```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "expiresIn": 3600,
  "user": {
    "id": "uuid",
    "email": "string",
    "firstName": "string",
    "lastName": "string",
    "role": "string",
    "avatar": "string"
  }
}
```

#### POST /api/auth/refresh
Обновление токенов.

**Request Body:**
```json
{
  "refreshToken": "string",
  "agent": "string"
}
```

**Response (200 OK):**
```json
{
  "accessToken": "string",
  "refreshToken": "string",
  "expiresIn": 3600
}
```

#### POST /api/auth/logout
Выход из системы.

**Headers:** `Authorization: Bearer <accessToken>`

**Response (200 OK):**
```json
{
  "message": "Successfully logged out"
}
```

#### POST /api/auth/forgot-password
Запрос на сброс пароля.

**Request Body:**
```json
{
  "email": "string"
}
```

**Response (200 OK):**
```json
{
  "message": "Password reset instructions sent to email"
}
```

---

## 3. Управление пользователями (Users)

### 3.1 Endpoints пользователей

#### GET /api/users
Получение списка пользователей с пагинацией и фильтрами.

**Query Parameters:**
- `page` (number, default: 1)
- `limit` (number, default: 20)
- `search` (string, optional) - поиск по имени/email
- `role` (string, optional) - фильтр по роли
- `department` (string, optional) - фильтр по отделу

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "email": "string",
      "firstName": "string",
      "lastName": "string",
      "avatar": "string",
      "position": "string",
      "department": "string",
      "role": "string",
      "status": "active|inactive",
      "createdAt": "timestamp"
    }
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "totalPages": 5
  }
}
```

#### GET /api/users/:id
Получение информации о конкретном пользователе.

**Response (200 OK):**
```json
{
  "id": "uuid",
  "email": "string",
  "firstName": "string",
  "lastName": "string",
  "avatar": "string",
  "phone": "string",
  "position": "string",
  "department": "string",
  "location": "string",
  "birthday": "date",
  "skype": "string",
  "role": "string",
  "status": "active|inactive",
  "skills": ["string"],
  "projects": [
    {
      "id": "uuid",
      "name": "string",
      "role": "string"
    }
  ],
  "team": [
    {
      "id": "uuid",
      "name": "string",
      "position": "string"
    }
  ],
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

#### GET /api/users/me
Получение информации о текущем пользователе.

**Headers:** `Authorization: Bearer <accessToken>`

**Response:** Аналогично GET /api/users/:id

#### PUT /api/users/:id
Обновление информации о пользователе.

**Request Body:**
```json
{
  "firstName": "string",
  "lastName": "string",
  "phone": "string",
  "position": "string",
  "department": "string",
  "location": "string",
  "birthday": "date",
  "skype": "string",
  "skills": ["string"]
}
```

**Response (200 OK):** Обновленный объект пользователя

#### POST /api/users/:id/avatar
Загрузка аватара пользователя.

**Content-Type:** `multipart/form-data`

**Response (200 OK):**
```json
{
  "avatar": "string (URL)"
}
```

---

## 4. Управление сотрудниками (Employees)

### 4.1 Endpoints сотрудников

#### GET /api/employees
Получение списка сотрудников с их рабочей нагрузкой.

**Query Parameters:**
- `page`, `limit`
- `department` (string, optional)
- `status` (string, optional): active|on_vacation|sick_leave

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "userId": "uuid",
      "firstName": "string",
      "lastName": "string",
      "email": "string",
      "avatar": "string",
      "position": "string",
      "level": "Junior|Middle|Senior|Lead",
      "department": "string",
      "status": "active|on_vacation|sick_leave",
      "workload": {
        "backlogTasks": 5,
        "inProgressTasks": 3,
        "inReviewTasks": 2
      }
    }
  ],
  "meta": { "total": 50, "page": 1, "limit": 20, "totalPages": 3 }
}
```

#### GET /api/employees/:id/tasks
Получение задач конкретного сотрудника.

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "string",
      "description": "string",
      "status": "backlog|in_progress|in_review|done",
      "priority": "low|medium|high|critical",
      "project": {
        "id": "uuid",
        "name": "string"
      },
      "dueDate": "timestamp",
      "createdAt": "timestamp"
    }
  ]
}
```

---

## 5. Управление проектами (Projects)

### 5.1 Endpoints проектов

#### GET /api/projects
Получение списка проектов.

**Query Parameters:**
- `page`, `limit`
- `status` (string, optional): active|completed|archived
- `search` (string, optional)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "code": "string (e.g., PN0001245)",
      "name": "string",
      "description": "string",
      "status": "active|completed|archived",
      "progress": 75,
      "startDate": "date",
      "endDate": "date",
      "team": [
        {
          "id": "uuid",
          "name": "string",
          "avatar": "string",
          "role": "string"
        }
      ],
      "tasksCount": {
        "total": 50,
        "completed": 30,
        "inProgress": 15,
        "backlog": 5
      },
      "createdAt": "timestamp"
    }
  ],
  "meta": { "total": 20, "page": 1, "limit": 10, "totalPages": 2 }
}
```

#### POST /api/projects
Создание нового проекта.

**Request Body:**
```json
{
  "name": "string",
  "description": "string",
  "startDate": "date",
  "endDate": "date",
  "teamMembers": ["uuid"]
}
```

**Response (201 Created):** Созданный объект проекта

#### GET /api/projects/:id
Получение детальной информации о проекте.

**Response (200 OK):**
```json
{
  "id": "uuid",
  "code": "string",
  "name": "string",
  "description": "string",
  "status": "active|completed|archived",
  "progress": 75,
  "startDate": "date",
  "endDate": "date",
  "budget": {
    "allocated": 100000,
    "spent": 75000,
    "currency": "RUB"
  },
  "team": [
    {
      "id": "uuid",
      "name": "string",
      "avatar": "string",
      "role": "string",
      "position": "string"
    }
  ],
  "tasks": [
    {
      "id": "uuid",
      "title": "string",
      "status": "string",
      "assignee": { "id": "uuid", "name": "string" }
    }
  ],
  "files": [
    {
      "id": "uuid",
      "name": "string",
      "url": "string",
      "size": 1024,
      "uploadedAt": "timestamp"
    }
  ],
  "createdAt": "timestamp",
  "updatedAt": "timestamp"
}
```

#### PUT /api/projects/:id
Обновление проекта.

#### DELETE /api/projects/:id
Удаление/архивирование проекта.

#### POST /api/projects/:id/members
Добавление участника в проект.

**Request Body:**
```json
{
  "userId": "uuid",
  "role": "string"
}
```

#### DELETE /api/projects/:id/members/:userId
Удаление участника из проекта.

---

## 6. Управление задачами (Tasks)

### 6.1 Endpoints задач

#### GET /api/tasks
Получение списка задач с фильтрами.

**Query Parameters:**
- `page`, `limit`
- `projectId` (uuid, optional)
- `assigneeId` (uuid, optional)
- `status` (string, optional): backlog|in_progress|in_review|done
- `priority` (string, optional): low|medium|high|critical

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "string",
      "description": "string",
      "status": "backlog|in_progress|in_review|done",
      "priority": "low|medium|high|critical",
      "project": {
        "id": "uuid",
        "name": "string",
        "code": "string"
      },
      "assignee": {
        "id": "uuid",
        "name": "string",
        "avatar": "string"
      },
      "reporter": {
        "id": "uuid",
        "name": "string"
      },
      "dueDate": "timestamp",
      "estimatedHours": 8,
      "loggedHours": 5,
      "tags": ["string"],
      "createdAt": "timestamp",
      "updatedAt": "timestamp"
    }
  ],
  "meta": { "total": 100, "page": 1, "limit": 20, "totalPages": 5 }
}
```

#### POST /api/tasks
Создание новой задачи.

**Request Body:**
```json
{
  "title": "string",
  "description": "string",
  "projectId": "uuid",
  "assigneeId": "uuid",
  "priority": "low|medium|high|critical",
  "dueDate": "timestamp",
  "estimatedHours": 8,
  "tags": ["string"]
}
```

**Response (201 Created):** Созданный объект задачи

#### GET /api/tasks/:id
Получение детальной информации о задаче.

#### PUT /api/tasks/:id
Обновление задачи.

#### PATCH /api/tasks/:id/status
Изменение статуса задачи.

**Request Body:**
```json
{
  "status": "backlog|in_progress|in_review|done"
}
```

#### DELETE /api/tasks/:id
Удаление задачи.

#### POST /api/tasks/:id/comments
Добавление комментария к задаче.

**Request Body:**
```json
{
  "content": "string"
}
```

#### GET /api/tasks/:id/comments
Получение комментариев к задаче.

#### POST /api/tasks/:id/files
Прикрепление файла к задаче.

---

## 7. Управление событиями (Events)

### 7.1 Endpoints событий

#### GET /api/events
Получение списка событий.

**Query Parameters:**
- `startDate` (date)
- `endDate` (date)
- `type` (string, optional): meeting|deadline|reminder|holiday

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "title": "string",
      "description": "string",
      "type": "meeting|deadline|reminder|holiday",
      "startDate": "timestamp",
      "endDate": "timestamp",
      "allDay": false,
      "location": "string",
      "participants": [
        {
          "id": "uuid",
          "name": "string",
          "avatar": "string",
          "status": "pending|accepted|declined"
        }
      ],
      "createdBy": {
        "id": "uuid",
        "name": "string"
      },
      "createdAt": "timestamp"
    }
  ]
}
```

#### POST /api/events
Создание нового события.

**Request Body:**
```json
{
  "title": "string",
  "description": "string",
  "type": "meeting|deadline|reminder|holiday",
  "startDate": "timestamp",
  "endDate": "timestamp",
  "allDay": false,
  "location": "string",
  "participantIds": ["uuid"],
  "reminders": [
    {
      "type": "email|push",
      "minutesBefore": 30
    }
  ]
}
```

**Response (201 Created):** Созданный объект события

#### PUT /api/events/:id
Обновление события.

#### DELETE /api/events/:id
Удаление события.

#### PATCH /api/events/:id/respond
Ответ на приглашение на событие.

**Request Body:**
```json
{
  "status": "accepted|declined"
}
```

---

## 8. Управление отпусками (Vacations)

### 8.1 Endpoints отпусков

#### GET /api/vacations
Получение списка заявок на отпуск.

**Query Parameters:**
- `page`, `limit`
- `status` (string, optional): pending|approved|rejected
- `userId` (uuid, optional)
- `year` (number, optional)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "user": {
        "id": "uuid",
        "name": "string",
        "email": "string",
        "avatar": "string",
        "position": "string"
      },
      "type": "annual|sick|unpaid|maternity",
      "startDate": "date",
      "endDate": "date",
      "totalDays": 14,
      "status": "pending|approved|rejected",
      "reason": "string",
      "approvedBy": {
        "id": "uuid",
        "name": "string"
      },
      "createdAt": "timestamp"
    }
  ],
  "meta": { "total": 30, "page": 1, "limit": 20, "totalPages": 2 }
}
```

#### POST /api/vacations
Создание заявки на отпуск.

**Request Body:**
```json
{
  "type": "annual|sick|unpaid|maternity",
  "startDate": "date",
  "endDate": "date",
  "reason": "string"
}
```

**Response (201 Created):** Созданный объект заявки

#### GET /api/vacations/:id
Получение информации о заявке на отпуск.

#### PATCH /api/vacations/:id/approve
Одобрение заявки на отпуск (для менеджера).

#### PATCH /api/vacations/:id/reject
Отклонение заявки на отпуск (для менеджера).

**Request Body:**
```json
{
  "reason": "string"
}
```

#### DELETE /api/vacations/:id
Отмена заявки на отпуск.

#### GET /api/vacations/balance
Получение баланса отпускных дней для текущего пользователя.

**Response (200 OK):**
```json
{
  "annualTotal": 28,
  "annualUsed": 10,
  "annualRemaining": 18,
  "sickTotal": 10,
  "sickUsed": 2,
  "sickRemaining": 8
}
```

---

## 9. Мессенджер (Messenger)

### 9.1 REST Endpoints

#### GET /api/conversations
Получение списка диалогов.

**Response (200 OK):**
```json
{
  "groups": [
    {
      "id": "uuid",
      "name": "string",
      "avatar": "string",
      "type": "group",
      "participants": [
        { "id": "uuid", "name": "string", "avatar": "string" }
      ],
      "lastMessage": {
        "id": "uuid",
        "content": "string",
        "sender": { "id": "uuid", "name": "string" },
        "timestamp": "timestamp"
      },
      "unreadCount": 5
    }
  ],
  "directMessages": [
    {
      "id": "uuid",
      "participant": {
        "id": "uuid",
        "name": "string",
        "avatar": "string",
        "status": "online|offline"
      },
      "type": "direct",
      "lastMessage": {
        "id": "uuid",
        "content": "string",
        "timestamp": "timestamp"
      },
      "unreadCount": 2
    }
  ]
}
```

#### POST /api/conversations
Создание нового диалога/группы.

**Request Body:**
```json
{
  "type": "group|direct",
  "name": "string (for groups)",
  "participantIds": ["uuid"]
}
```

#### GET /api/conversations/:id/messages
Получение сообщений диалога.

**Query Parameters:**
- `page`, `limit`
- `before` (timestamp, optional) - для пагинации

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "content": "string",
      "type": "text|file|image",
      "sender": {
        "id": "uuid",
        "name": "string",
        "avatar": "string"
      },
      "attachments": [
        {
          "id": "uuid",
          "name": "string",
          "url": "string",
          "type": "string",
          "size": 1024
        }
      ],
      "readBy": ["uuid"],
      "createdAt": "timestamp"
    }
  ],
  "meta": { "hasMore": true }
}
```

#### POST /api/conversations/:id/messages
Отправка сообщения.

**Request Body:**
```json
{
  "content": "string",
  "type": "text|file|image",
  "attachmentIds": ["uuid"]
}
```

### 9.2 WebSocket Events

#### Подключение
```
ws://api.domain.com/ws?token=<accessToken>
```

#### События от сервера к клиенту

**message.new**
```json
{
  "event": "message.new",
  "data": {
    "conversationId": "uuid",
    "message": { /* Message object */ }
  }
}
```

**message.read**
```json
{
  "event": "message.read",
  "data": {
    "conversationId": "uuid",
    "messageIds": ["uuid"],
    "userId": "uuid"
  }
}
```

**user.status**
```json
{
  "event": "user.status",
  "data": {
    "userId": "uuid",
    "status": "online|offline"
  }
}
```

**typing.start / typing.stop**
```json
{
  "event": "typing.start",
  "data": {
    "conversationId": "uuid",
    "userId": "uuid",
    "userName": "string"
  }
}
```

#### События от клиента к серверу

**message.send**
```json
{
  "event": "message.send",
  "data": {
    "conversationId": "uuid",
    "content": "string",
    "type": "text"
  }
}
```

**message.markRead**
```json
{
  "event": "message.markRead",
  "data": {
    "conversationId": "uuid",
    "messageIds": ["uuid"]
  }
}
```

**typing**
```json
{
  "event": "typing",
  "data": {
    "conversationId": "uuid",
    "isTyping": true
  }
}
```

---

## 10. Уведомления (Notifications)

### 10.1 Endpoints уведомлений

#### GET /api/notifications
Получение списка уведомлений.

**Query Parameters:**
- `page`, `limit`
- `read` (boolean, optional)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "type": "task_assigned|event_reminder|vacation_approved|message_received",
      "title": "string",
      "message": "string",
      "data": {
        "entityType": "task|event|vacation|message",
        "entityId": "uuid"
      },
      "read": false,
      "createdAt": "timestamp"
    }
  ],
  "meta": {
    "total": 50,
    "unreadCount": 10,
    "page": 1,
    "limit": 20
  }
}
```

#### PATCH /api/notifications/:id/read
Отметить уведомление как прочитанное.

#### PATCH /api/notifications/read-all
Отметить все уведомления как прочитанные.

#### DELETE /api/notifications/:id
Удаление уведомления.

### 10.2 Push-уведомления

Для отправки push-уведомлений использовать Firebase Cloud Messaging (FCM) или Web Push API.

**Регистрация устройства:**

#### POST /api/notifications/devices
**Request Body:**
```json
{
  "token": "string",
  "platform": "web|ios|android"
}
```

---

## 11. Аналитика и статистика (Analytics)

### 11.1 Endpoints аналитики

#### GET /api/analytics/dashboard
Получение сводной информации для дашборда.

**Response (200 OK):**
```json
{
  "projects": {
    "total": 15,
    "active": 10,
    "completed": 5
  },
  "tasks": {
    "total": 150,
    "completed": 80,
    "inProgress": 50,
    "overdue": 10
  },
  "team": {
    "totalEmployees": 50,
    "onVacation": 5,
    "onSickLeave": 2
  },
  "upcomingEvents": [
    {
      "id": "uuid",
      "title": "string",
      "date": "timestamp"
    }
  ],
  "recentActivity": [
    {
      "id": "uuid",
      "type": "task_completed|project_created|user_joined",
      "description": "string",
      "user": { "id": "uuid", "name": "string" },
      "timestamp": "timestamp"
    }
  ]
}
```

#### GET /api/analytics/projects/:id/report
Получение отчета по проекту.

**Response (200 OK):**
```json
{
  "project": {
    "id": "uuid",
    "name": "string"
  },
  "timeline": {
    "startDate": "date",
    "endDate": "date",
    "daysRemaining": 30
  },
  "progress": {
    "overall": 75,
    "tasksCompleted": 45,
    "tasksTotal": 60
  },
  "team": {
    "size": 8,
    "workloadDistribution": [
      { "userId": "uuid", "name": "string", "tasksAssigned": 10 }
    ]
  },
  "burndown": [
    { "date": "date", "remaining": 40 },
    { "date": "date", "remaining": 35 }
  ]
}
```

---

## 12. Интеграции

### 12.1 Google Calendar

#### GET /api/integrations/google/auth
Получение URL для авторизации в Google.

**Response (200 OK):**
```json
{
  "authUrl": "string"
}
```

#### POST /api/integrations/google/callback
Обработка callback от Google OAuth.

**Request Body:**
```json
{
  "code": "string"
}
```

#### POST /api/integrations/google/sync
Синхронизация событий с Google Calendar.

**Request Body:**
```json
{
  "direction": "import|export|both"
}
```

### 12.2 Telegram Bot (опционально)

#### POST /api/integrations/telegram/link
Связывание Telegram аккаунта.

**Request Body:**
```json
{
  "telegramUserId": "string"
}
```

---

## 13. Загрузка файлов

### 13.1 Endpoints файлов

#### POST /api/files/upload
Загрузка файла.

**Content-Type:** `multipart/form-data`

**Response (201 Created):**
```json
{
  "id": "uuid",
  "name": "string",
  "url": "string",
  "type": "string",
  "size": 1024,
  "uploadedAt": "timestamp"
}
```

#### GET /api/files/:id
Получение информации о файле.

#### DELETE /api/files/:id
Удаление файла.

---

## 14. Требования к безопасности

### 14.1 Аутентификация
- JWT токены с коротким временем жизни (15-60 минут)
- Refresh токены с более длительным временем жизни (7-30 дней)
- Хранение refresh токенов в HttpOnly cookies
- Blacklist для отозванных токенов

### 14.2 Авторизация
- Role-Based Access Control (RBAC)
- Роли: admin, manager, employee
- Проверка прав доступа на уровне API

### 14.3 Защита данных
- HTTPS обязателен
- Rate limiting для API endpoints
- Валидация всех входных данных
- Защита от SQL injection, XSS, CSRF
- Шифрование чувствительных данных в БД

---

## 15. Требования к производительности

- Время ответа API < 200ms для 95% запросов
- Поддержка до 1000 одновременных WebSocket соединений
- Кэширование часто запрашиваемых данных
- Пагинация для всех list endpoints
- Оптимизация запросов к БД (индексы, N+1 prevention)

---

## 16. Модель данных (Database Schema)

### 16.1 Основные таблицы

```sql
-- Users
CREATE TABLE users (
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
);

-- Projects
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(20) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  status VARCHAR(20) DEFAULT 'active',
  start_date DATE,
  end_date DATE,
  budget_allocated DECIMAL(15,2),
  budget_spent DECIMAL(15,2) DEFAULT 0,
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Project Members
CREATE TABLE project_members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  role VARCHAR(50),
  joined_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(project_id, user_id)
);

-- Tasks
CREATE TABLE tasks (
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
);

-- Events
CREATE TABLE events (
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
);

-- Event Participants
CREATE TABLE event_participants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID REFERENCES events(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  status VARCHAR(20) DEFAULT 'pending',
  UNIQUE(event_id, user_id)
);

-- Vacations
CREATE TABLE vacations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type VARCHAR(20) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  total_days INTEGER NOT NULL,
  status VARCHAR(20) DEFAULT 'pending',
  reason TEXT,
  approved_by UUID REFERENCES users(id),
  approved_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Conversations
CREATE TABLE conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  type VARCHAR(20) NOT NULL,
  name VARCHAR(255),
  avatar_url VARCHAR(500),
  created_at TIMESTAMP DEFAULT NOW()
);

-- Conversation Participants
CREATE TABLE conversation_participants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  last_read_at TIMESTAMP,
  UNIQUE(conversation_id, user_id)
);

-- Messages
CREATE TABLE messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id),
  content TEXT NOT NULL,
  type VARCHAR(20) DEFAULT 'text',
  created_at TIMESTAMP DEFAULT NOW()
);

-- Notifications
CREATE TABLE notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  type VARCHAR(50) NOT NULL,
  title VARCHAR(255) NOT NULL,
  message TEXT,
  data JSONB,
  read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Files
CREATE TABLE files (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  url VARCHAR(500) NOT NULL,
  type VARCHAR(100),
  size INTEGER,
  uploaded_by UUID REFERENCES users(id),
  entity_type VARCHAR(50),
  entity_id UUID,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Refresh Tokens
CREATE TABLE refresh_tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  token VARCHAR(500) NOT NULL,
  user_agent TEXT,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);
```

---

## 17. Логирование и мониторинг

### 17.1 Требования к логированию
- Все HTTP запросы с временем выполнения
- Ошибки с полным stack trace
- Действия пользователей (audit log)
- WebSocket события

### 17.2 Метрики
- Количество активных пользователей
- Время ответа API (p50, p95, p99)
- Количество ошибок по типам
- Использование ресурсов сервера

---

## 18. Развертывание

### 18.1 Окружения
- Development (local)
- Staging
- Production

### 18.2 CI/CD
- Автоматическое тестирование
- Линтинг кода
- Сборка Docker образов
- Автоматическое развертывание на staging
- Ручное подтверждение для production

### 18.3 Конфигурация
Все конфигурационные параметры должны быть вынесены в переменные окружения:
- DATABASE_URL
- REDIS_URL
- JWT_SECRET
- JWT_REFRESH_SECRET
- GOOGLE_CLIENT_ID
- GOOGLE_CLIENT_SECRET
- AWS_S3_BUCKET (для файлов)
- SMTP настройки для email

---

## 19. Дополнительные требования

### 19.1 Документация
- OpenAPI/Swagger спецификация для всех endpoints
- README с инструкциями по запуску
- Описание архитектуры

### 19.2 Тестирование
- Unit тесты с покрытием > 80%
- Integration тесты для критических путей
- E2E тесты для основных сценариев

### 19.3 Миграции
- Версионирование схемы БД
- Возможность отката миграций
- Seed данные для разработки
