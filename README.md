AssetFlow Asset Management System

AssetFlow is a scalable and modular full-stack application for managing organizational assets, allocations, bookings, maintenance, and audits. It features a modern React frontend and a robust Go backend connected to a live Supabase PostgreSQL database.

The application follows a clean, domain-driven architecture where each business module is organized into its own package with dedicated handler, service, and repository layers on the backend, seamlessly communicating with the frontend React UI.

---

Features

- Real-time API Integration (No Mock Data)
- Supabase PostgreSQL Database Integration
- JWT Authentication & Authorization
- Role-Based Access Control (RBAC)
- Organization & Employee Management
- Asset Registration & Asset Directory
- Asset Allocation & Transfers
- Resource Booking with Conflict Detection
- Maintenance Workflow
- Asset Audits & Verification
- Dashboard KPIs & Analytics
- Standardized JSON Responses
- Request Logging Middleware

---

Project Structure

assetflow-backend/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── auth/
│   ├── organization/
│   ├── assets/
│   ├── allocations/
│   ├── bookings/
│   ├── maintenance/
│   ├── audits/
│   ├── dashboard/
│   └── utils/
│
├── migrations/
├── pkg/
├── go.mod
├── go.sum
└── .env

---

Architecture

The backend follows a layered architecture:

HTTP Request
      │
      ▼
 Handler
      │
      ▼
 Service
      │
      ▼
 Repository
      │
      ▼
 PostgreSQL

Handler

Responsible for:

- Parsing HTTP requests
- Input validation
- Calling business services
- Returning standardized JSON responses

Service

Responsible for:

- Business rules
- Validation
- Authorization logic
- Coordinating repositories

Repository

Responsible for:

- SQL queries
- Database transactions
- Data persistence

---

Domain Modules

Authentication

Handles:

- User Registration
- Login
- JWT Generation
- Password Hashing
- User Authentication

---

Organization

Handles:

- Departments
- Employees
- Employee Directory

---

Assets

Handles:

- Asset Registration
- Asset Details
- Asset Status
- Asset Categories

---

Allocations

Handles:

- Asset Assignment
- Asset Transfers
- Asset Returns
- Allocation History

---

Bookings

Handles:

- Resource Booking
- Availability Checks
- Booking Approval
- Overlap Validation

---

Maintenance

Handles:

- Maintenance Requests
- Scheduled Maintenance
- Completion Tracking
- Maintenance History

---

Audits

Handles:

- Asset Verification
- Audit Logs
- Compliance Checks

---

Dashboard

Provides:

- Asset Statistics
- Allocation Summary
- Booking Metrics
- Maintenance Metrics
- Audit KPIs

---

Models

The "models" package contains the application's core domain entities that map directly to the PostgreSQL schema.

Examples include:

- User
- Asset
- Allocation
- Booking
- Maintenance
- Audit

Shared enums such as:

- User Roles
- Asset Status
- Booking Status

are defined in:

internal/models/enums.go

---

Middleware

Authentication Middleware

- JWT Verification
- Role Validation
- Protected Routes

Logger Middleware

- Request Logging
- Response Status
- Execution Time

---

Utilities

Common reusable utilities include:

- JWT helpers
- Password hashing
- Standard API responses

---

Database

Database configuration is managed inside:

internal/database/postgres.go

Database migrations are located in:

migrations/

Example:

000001_init_schema.up.sql
000001_init_schema.down.sql

---

Environment Variables

Create a ".env" file in the project root.

PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=assetflow
DB_SSLMODE=disable

JWT_SECRET=your-secret-key

---

Getting Started

1. Clone the repository

git clone https://github.com/your-username/assetflow-backend.git

cd assetflow-backend

2. Install dependencies

go mod tidy

3. Configure environment

Create the ".env" file and update the database credentials.

4. Run database migrations

Apply the SQL migration files located in:

migrations/

using your preferred migration tool.

5. Run the application

go run cmd/server/main.go

The API will start on:

http://localhost:8080

---

API Design

The API follows RESTful conventions.

Example endpoints:

POST   /api/auth/login
POST   /api/auth/register

GET    /api/assets
POST   /api/assets

GET    /api/allocations
POST   /api/allocations

GET    /api/bookings
POST   /api/bookings

GET    /api/maintenance
POST   /api/maintenance

GET    /api/audits

GET    /api/dashboard

---

Security

- JWT Authentication
- Password Hashing with bcrypt
- Role-Based Access Control (RBAC)
- Protected API Routes
- Centralized Authentication Middleware

---

Tech Stack

**Frontend:**
- React (Vite)
- TypeScript
- Tailwind CSS
- Axios

**Backend:**
- Go (Golang)
- PostgreSQL (Supabase)
- JWT Authentication
- bcrypt
- SQL Migrations
- REST API

---

Future Improvements

- Refresh Tokens
- API Documentation (Swagger/OpenAPI)
- Redis Caching
- Background Workers
- Email Notifications
- File Upload Support
- Docker & Docker Compose
- CI/CD Pipeline
- Unit & Integration Testing
- Prometheus & Grafana Monitoring

---

License

This project is licensed under the MIT License.
