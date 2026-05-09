# Task Management Application
- Created At 2025.12.22
- Version 0.1.0
- Created By k.ito

## Overview
The Purpose of this Application.

### Why
Training for making Web Applications made by Golang and React.

### Targets
This document is intended for developers themselves.
This Application is Web Application.

### Technical Assumptions
- Backend: Golang standard net/http.
- Frontend: React (SPA).
- State management: minimal (no Redux in v0.1.0).

## Detailed Specifications
This is minimum specifications.

### User Login form 
- Login by Email and Passwords.
- Login information is managed by JWT tokens.
- Supabase is used for Authentication.
- Login form has two text boxes for Email and Password, one button for login.
- The text box for password masks user inputs.

#### Non-Goals
- User registration.
- Password reset.
- Email verification.
- Role based access control.
- Production-level security considerations.

---

### Task Management form
- Tasks registration.
  - POST /tasks
  - Returns a new list of tasks for the authenticated user.
- Tasks update.
  - PATCH /tasks/{id}
  - Updates title, description, and status.
  - Request body:
    - title: string
    - description: string
    - status: TODO / DOING / DONE
  - Returns only HTTP status code.
- Tasks delete.
  - DELETE /tasks/{id}
  - Returns only HTTP statsu code
- Tasks view in list.
  - GET /tasks
  - Returns a list of tasks for the authenticated user.

#### Task Attributes
- Title
- Description (optional)
- Status (TODO / DOING / DONE)

#### Notes
- Split into multiple components within the same page, not across different pages.
