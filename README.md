# GanttPro - Project Management System

Full-stack project management application dengan Gantt Chart visualization untuk ERP FMLX.

## Tech Stack

**Frontend:** Vue 3, Vue Router, Vite
**Backend:** Golang, Gin, GORM, JWT, bcrypt
**Database:** PostgreSQL

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 16+
- PostgreSQL 12+

### 1. Setup Database

```bash
# Create database
psql -U postgres -c "CREATE DATABASE ganttpro_db;"

# Run setup script
psql -U postgres -d ganttpro_db -f backend/database/migrations/RUN_ME_FIRST.sql
```

### 2. Backend Setup

```bash
cd backend

# Copy and configure environment
cp .env.example .env
# Edit .env dengan database credentials Anda

# Install dependencies dan jalankan
go mod download
go run main.go
```

Backend berjalan di `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend

# Install dan jalankan
npm install
npm run dev
```

Frontend berjalan di `http://localhost:5173`

## Default Login

**Username:** `admin`
**Password:** `admin123`

## Configuration

### Backend (.env)
```env
PORT=8080
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ganttpro_db
JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=168
ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend (.env)
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## Main API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/register` | Register |
| GET | `/api/v1/auth/profile` | Get profile |
| GET | `/api/ppic/gantt-data` | PPIC Gantt data |
| GET | `/health` | Health check |

## Features

- User authentication & authorization (Admin, PPIC, Operator roles)
- PPIC Gantt Chart scheduling
- Machine assignments & tracking
- Job order management
- PEM & Toolpather file management
- Operation plan approvals

## Troubleshooting

**Database connection error:** Pastikan PostgreSQL berjalan dan kredensial di `.env` benar
**CORS error:** Pastikan `ALLOWED_ORIGINS` di backend `.env` sudah benar
**Token expired:** Login ulang (default: token berlaku 168 jam / 7 hari)

---

**Tubes Compro Formulatrix**
