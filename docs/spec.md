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

### User Login form 
- Login by Usernames and Passwords.
- Login information is managed by JWT tokens.
- We will use Supabase for JWT Authorization.
- Login form has two text boxes for UserName and Password, one button for login.
- The text box for password masks user inputs.

#### Non-Goals
- Password reset.
- Email verification.
- Role based access control.
- Production-level security considerations.

