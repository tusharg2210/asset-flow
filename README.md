# 🏢 AssetFlow

**A full-stack, enterprise-grade Asset Management System built with React, Go, and PostgreSQL.**

AssetFlow is a scalable application designed to track, manage, and optimize organizational assets. It supports end-to-end workflows including employee management, asset directory tracking, allocations, reservations with conflict detection, ticketing for maintenance, physical audit cycles, and real-time dashboard analytics.

---

## ✨ Core Features

*   🔐 **Role-Based Access Control (RBAC):** Distinct workflows for `Admin`, `Asset Manager`, and `Employee` user tiers using JWT authentication.
*   📦 **Asset Directory & Lifecycle:** Comprehensive registration for physical assets (electronics, vehicles, furniture), tracking condition, location, and acquisition details.
*   🔄 **Allocations & Transfers:** Assign assets to specific employees or departments and maintain an immutable historical chain of custody.
*   📅 **Resource Booking:** Smart booking calendar for shared resources (e.g., conference rooms, projectors) with strict overlap conflict detection.
*   🛠 **Maintenance Workflows:** Raise, track, and resolve maintenance tickets for broken or malfunctioning assets, automatically preventing their booking while disabled.
*   🔍 **Audit Cycles:** Verify the physical presence and condition of scoped assets via discrepancy flagging.
*   📊 **Real-Time Analytics:** Interactive KPI dashboards detailing resource utilization, idling periods, and actionable alerts.

---

## 🛠 Tech Stack & Architecture

AssetFlow employs a modern, decoupled architecture connecting a dynamic single-page React frontend with a high-performance RESTful Go API, powered by a live Supabase PostgreSQL database.

### Frontend
*   **Framework:** React 18, Vite, TypeScript
*   **Styling:** Tailwind CSS (utility-first UI)
*   **Networking:** Axios (with centralized interceptors)
*   **Icons:** Lucide React

### Backend
*   **Language:** Go (Golang)
*   **Framework:** Echo (v4) for HTTP routing and middleware
*   **Database Interface:** `pgxpool` for high-performance PostgreSQL connection pooling
*   **Security:** JWT (JSON Web Tokens) & `bcrypt` for password hashing
*   **Database:** Supabase PostgreSQL

### Data Flow Overview
1.  The client (React) requests an action via Axios, attaching a Bearer JWT.
2.  The Echo server routes the HTTP request through standard logging and authentication middleware.
3.  The request hits the **Handler layer** (parsing inputs/validation).
4.  The request cascades down to the **Repository layer**, which executes the raw SQL queries using `pgxpool`.
5.  Standardized JSON responses (`util.Success` or `util.Fail`) are returned to the client and rendered dynamically.

---

## 📁 Repository Structure

### 🖥 Frontend (`/frontend`)
*   `src/api/` - Contains `axiosClient.ts` (interceptors) and `endpoints.ts` (centralized API route constants).
*   `src/components/` - Reusable UI elements (`AppLayout`, `Modal`, `Badge`, `FormInput`).
*   `src/pages/` - The core application views:
    *   `DashboardPage.tsx` - KPI charts and system-wide alerts.
    *   `AssetsPage.tsx` - Filterable grid/list directory of all company assets.
    *   `AllocationsPage.tsx` - Assignments to individuals.
    *   `MaintenancePage.tsx` - Kanban-style ticketing for broken assets.
    *   `BookingsPage.tsx` - Timeline view for reserving shared resources.
    *   `AuditPage.tsx` - Discrepancy flagging.
    *   `OrganizationPage.tsx` - Management of categories, departments, and user creation.
    *   `NotificationsPage.tsx` - System logs and approval alerts.

### ⚙️ Backend (`/backend`)
*   `cmd/` - Contains `main.go`, the application's entry point that spins up the Echo server.
*   `internal/`
    *   `config/` - Environment variables setup (`.env` parsing).
    *   `db/` - Connection pooling logic and an automated `seed.go` script to bulk-populate the DB.
    *   `handler/` - HTTP request parsing (e.g., `asset.go`, `auth.go`, `booking.go`).
    *   `middleware/` - Application-level intercepts (`auth_middleware.go` for JWT/Roles, `logger.go`).
    *   `model/` - Struct definitions mapping exactly to the PostgreSQL schema (`asset.go`, `enums.go`).
    *   `repository/` - The SQL query logic, isolating the database from the business handlers.
    *   `router/` - Defines the public and protected API routes mapped to their respective handlers.

---

## 🔌 API Endpoints Summary

All routes are prefixed with `/api`. Protected routes require a valid JWT Bearer token.

| Module | Endpoints |
| :--- | :--- |
| **Auth** | `POST /auth/login`, `POST /auth/signup` |
| **Assets** | `GET /assets`, `POST /assets`, `GET /assets/:id` |
| **Allocations** | `GET /allocations`, `POST /allocations`, `PUT /allocations/:id/return` |
| **Bookings** | `GET /assets/:id/bookings`, `POST /bookings/create`, `PUT /bookings/:id/cancel` |
| **Maintenance** | `GET /maintenance`, `POST /maintenance`, `PUT /maintenance/:id` |
| **Audits** | `GET /audits/cycle/:id/assets`, `POST /audits/cycle/:id/report` |
| **Dashboard** | `GET /dashboard/metrics`, `GET /dashboard/alerts` |
| **Reports** | `GET /reports/analytics` |

---

## 🚀 Setup & Installation

### Prerequisites
*   [Node.js](https://nodejs.org/) (v18+)
*   [Go](https://go.dev/) (1.21+)
*   A [Supabase](https://supabase.com/) PostgreSQL database (or local PostgreSQL instance).

### 1. Database Initialization
Create a new Supabase project and execute the provided `schema.sql` (if available) to generate the tables, or use the built-in Go struct models.

### 2. Backend Setup
1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install Go modules:
    ```bash
    go mod tidy
    ```
3.  Create a `.env` file in the `/backend` root:
    ```env
    # Server config
    PORT=8080

    # JWT Config
    JWT_SECRET=super_secret_jwt_key_change_me
    JWT_EXPIRATION_HOURS=24

    # DB Config (Supabase Postgres)
    DB_HOST=db.your-supabase-url.supabase.co
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=your_password
    DB_NAME=postgres
    DB_SSLMODE=disable
    ```
4.  Run the application:
    ```bash
    go run cmd/main.go
    ```
    *The server will start on port 8080 and automatically run the `seed.go` script if the database is empty, generating 100+ demo assets, users, and allocations.*

### 3. Frontend Setup
1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  The frontend is configured to talk to `http://127.0.0.1:8080` by default via `src/api/axiosClient.ts`.
4.  Start the Vite development server:
    ```bash
    npm run dev
    ```
    *The application will be accessible at http://localhost:5173.*

### 4. Logging In
If the database was freshly seeded by the backend, use the following credentials to access the full admin suite:
*   **Email:** `admin@company.com`
*   **Password:** `SecurePassword123!`

---

## 🔒 Security Posture
*   **JWT Integrity:** Stateless authentication requiring token validation on every protected route.
*   **Role Validation:** Admin endpoints immediately reject (`403 Forbidden`) requests originating from `Employee` tokens.
*   **SQL Injection Prevention:** Go's `pgxpool` strictly uses parameterized queries (`$1, $2`).
*   **CORS Policies:** Configured in `cmd/main.go` to strictly allow trusted origins.

---

*AssetFlow - Designed for absolute visibility over your organization's physical footprint.*
