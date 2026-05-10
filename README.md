# task-app

A simple task management app for learning Go and React.

## Overview

This project is a small full-stack web application.  
It uses Go for the backend, React + TypeScript for the frontend, PostgreSQL for the database, and Supabase Auth for login.

### Current Status

- Email/password login with Supabase Auth
- Task list API with JWT authentication (`GET /tasks`)
- User-specific task retrieval
- Basic backend layers: Handler / Usecase / Repository
- Docker + PostgreSQL setup

### Specifications

See `docs/spec.md` for details.

Current scope:

- User login with Supabase
- Task CRUD operations
- Task fields:
  - Title
  - Description (optional)
  - Status (`TODO` / `DOING` / `DONE`)

---

## Tech Stack

### Backend

- Go 1.23+
- net/http
- PostgreSQL
- JWT verification with ES256 / JWKS
- Docker

### Frontend

- React 18 + TypeScript
- Vite
- Supabase JS Client

### Infrastructure

- Docker Compose
- PostgreSQL

---

## Getting Started

### 1. Clone this repository

```bash
git clone https://github.com/koh5963/task-app.git
cd task-app
```

### 2. Set up environment variables

Create backend environment variables.

```bash
cp backend/.env.example backend/.env
```

Create frontend environment variables if needed.

```bash
cp frontend/.env.example frontend/.env
```

Then edit each `.env` file for your local environment.

### 3. Start Docker services

```bash
docker-compose up -d
```

### 4. Start the backend

```bash
cd backend
go run ./cmd/api
```

### 5. Start the frontend

```bash
cd frontend
npm install
npm run dev
```

---

## Local URLs

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

---

## Project Structure

```text
task-app/
├── backend/
│   ├── cmd/
│   │   └── api/
│   ├── internal/
│   │   ├── auth/
│   │   ├── domain/
│   │   ├── handler/
│   │   ├── infra/
│   │   └── usecase/
│   └── .env.example
├── frontend/
│   ├── src/
│   └── vite.config.ts
├── db/
│   └── init/
├── docs/
│   └── spec.md
├── docker-compose.yaml
└── README.md
```

---

## API

Authentication is handled by Supabase Auth.  
The frontend sends the access token to the backend with the `Authorization` header.

```http
Authorization: Bearer <access_token>
```

The backend verifies the JWT and uses the `sub` claim as the user ID.

### Task APIs

| Method | Endpoint | Description |
| --- | --- | --- |
| GET | `/tasks` | Get tasks for the logged-in user |
| POST | `/tasks` | Create a new task |
| PATCH | `/tasks/{id}` | Update a task |
| DELETE | `/tasks/{id}` | Delete a task |

---

## TODO

- [x] Add task creation API (`POST /tasks`)
- [x] Add task update API (`PATCH /tasks/{id}`)
- [x] Add task delete API (`DELETE /tasks/{id}`)
- [x] Connect create/update/delete from the frontend
- [ ] Improve loading and error messages
- [ ] Improve UI styling
- [ ] Add filtering by task status
- [ ] Add due date support
- [ ] Extract auth logic into a small reusable package

---

## License

MIT License