# GanttPro - Project Management System

Full-stack project management application dengan Gantt Chart visualization.

## 🏗️ Tech Stack

### Frontend
- **Vue 3** - Progressive JavaScript Framework
- **Vue Router** - Official router for Vue.js
- **Vite** - Next Generation Frontend Tooling

### Backend
- **Golang** - High-performance backend language
- **Gin** - Web framework for Go
- **GORM** - ORM library for Go
- **JWT** - JSON Web Tokens for authentication
- **bcrypt** - Password hashing

### Database
- **MySQL** or **PostgreSQL** - Relational database

## 📁 Project Structure

```
Project/
├── backend/                 # Golang Backend API
│   ├── config/             # Configuration management
│   ├── database/           # Database connection & migration
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # Middleware (CORS, Auth)
│   ├── models/             # Database models
│   ├── repository/         # Data access layer
│   ├── routes/             # API routes
│   ├── services/           # Business logic
│   ├── utils/              # Utility functions
│   ├── scripts/            # Database scripts
│   ├── .env.example        # Environment variables example
│   ├── go.mod              # Go dependencies
│   ├── main.go             # Entry point
│   └── README.md           # Backend documentation
│
├── frontend/               # Vue.js Frontend
│   ├── public/             # Static assets
│   ├── src/
│   │   ├── assets/         # Images, styles
│   │   ├── components/     # Vue components
│   │   ├── router/         # Vue Router config
│   │   ├── services/       # API services
│   │   ├── App.vue         # Root component
│   │   └── main.js         # Entry point
│   ├── .env                # Environment variables
│   ├── package.json        # NPM dependencies
│   └── vite.config.js      # Vite configuration
│
├── API_TESTING.md          # API testing guide
├── QUICKSTART.md           # Quick start guide
└── README.md               # This file
```

## 🚀 Quick Start

### Prerequisites
- **Go** 1.21 or higher
- **Node.js** 16 or higher
- **MySQL** 5.7+ or **PostgreSQL** 12+
- **Git**

### 1. Clone Repository
```bash
git clone <repository-url>
cd Project
```

### 2. Backend Setup

```bash
cd backend

# Copy environment file
cp .env.example .env

# Edit .env with your database credentials
# nano .env  # or use any text editor

# Download Go dependencies
go mod download

# Create database
# For MySQL:
mysql -u root -p -e "CREATE DATABASE ganttpro_db;"

# For PostgreSQL:
psql -U postgres -c "CREATE DATABASE ganttpro_db;"

# Run database setup script (optional, for sample users)
# MySQL:
mysql -u root -p ganttpro_db < scripts/setup_database.sql

# PostgreSQL:
psql -U postgres -d ganttpro_db -f scripts/setup_database_postgres.sql

# Run backend server
go run main.go
```

Backend akan berjalan di: `http://localhost:8080`

### 3. Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend akan berjalan di: `http://localhost:5173`

## 🔐 Default Login Credentials

Setelah menjalankan database setup script:

**Admin Account:**
- Username: `admin`
- Password: `admin123`

**Test Account:**
- Username: `testuser`
- Password: `admin123`

## 📡 API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | User login | No |
| GET | `/api/v1/auth/profile` | Get user profile | Yes |

### Health Check

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | API health check | No |

## 🔧 Configuration

### Backend (.env)
```env
PORT=8080
ENV=development

# Database
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ganttpro_db

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY_HOURS=24

# CORS
ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend (.env)
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

## 🧪 Testing

### Test API dengan PowerShell

```powershell
# Health check
Invoke-RestMethod -Uri "http://localhost:8080/health"

# Login
$body = @{
    userid = "admin"
    password = "admin123"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $body
```

Lihat [API_TESTING.md](./API_TESTING.md) untuk panduan testing lengkap.

## 🏗️ Development

### Backend Development

```bash
# Run with hot reload (install air first)
go install github.com/cosmtrek/air@latest
air

# Build for production
go build -o ganttpro-backend

# Run tests
go test ./...
```

### Frontend Development

```bash
# Run dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## 📚 Documentation

- [Backend README](./backend/README.md) - Detailed backend documentation
- [API Testing Guide](./API_TESTING.md) - Complete API testing guide
- [Quick Start Guide](./QUICKSTART.md) - Quick setup instructions

## 🔒 Security Features

- ✅ Password hashing with bcrypt
- ✅ JWT authentication
- ✅ CORS protection
- ✅ SQL injection prevention (GORM)
- ✅ XSS protection
- ✅ Secure password validation

## 🌟 Features (Current)

- ✅ User authentication (Login/Register)
- ✅ JWT-based session management
- ✅ Protected routes
- ✅ Role-based access control
- ✅ Modern UI design
- ✅ Responsive layout

## 🚧 Roadmap (Future Features)

- [ ] Project management
- [ ] Task creation and management
- [ ] Gantt chart visualization
- [ ] Team collaboration
- [ ] Real-time updates
- [ ] File attachments
- [ ] Activity timeline
- [ ] Email notifications
- [ ] Dashboard analytics

## 🐛 Troubleshooting

### Backend tidak bisa connect ke database
```bash
# Cek apakah database sudah dibuat
# MySQL:
mysql -u root -p -e "SHOW DATABASES;"

# PostgreSQL:
psql -U postgres -l

# Cek kredensial di .env
# Pastikan DB_USER, DB_PASSWORD, DB_NAME sudah benar
```

### CORS error di frontend
```bash
# Pastikan backend running
# Cek ALLOWED_ORIGINS di backend/.env
# Harus include http://localhost:5173
```

### Token expired
```bash
# JWT token berlaku 24 jam (default)
# Login ulang untuk mendapatkan token baru
```

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

MIT License - feel free to use this project for learning or production.

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- Vue.js team for the amazing framework
- Gin team for the excellent Go web framework
- GORM team for the powerful ORM
- All open source contributors

---

**Made with ❤️ using Vue.js and Golang**
=======
# compro-erp-fmlx
tubes Compro Formulatrix
