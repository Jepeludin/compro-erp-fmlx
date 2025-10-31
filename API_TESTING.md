# API Testing Guide

## Using PowerShell (Windows)

### 1. Health Check
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
```

### 2. Register New User
```powershell
$body = @{
    userid = "testuser"
    password = "password123"
    name = "Test User"
    email = "test@example.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -ContentType "application/json" -Body $body
```

### 3. Login
```powershell
$loginBody = @{
    userid = "testuser"
    password = "password123"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody

# Save token
$token = $response.data.token
Write-Host "Token: $token"
```

### 4. Get Profile (Protected Route)
```powershell
# Use the token from login
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/profile" -Method GET -Headers $headers
```

## Using cURL

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "userid": "testuser",
    "password": "password123",
    "name": "Test User",
    "email": "test@example.com"
  }'
```

### 3. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "userid": "testuser",
    "password": "password123"
  }'
```

### 4. Get Profile
```bash
# Replace YOUR_TOKEN with actual token from login
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Using Frontend (Vue.js)

The API service is already configured in `frontend/src/services/api.js`

### Example Usage in Vue Component:

```javascript
import api from '@/services/api';

// Login
async function login() {
  try {
    const response = await api.login('testuser', 'password123');
    console.log('Login success:', response);
    // Token automatically saved to localStorage
  } catch (error) {
    console.error('Login failed:', error.message);
  }
}

// Register
async function register() {
  try {
    const userData = {
      userid: 'newuser',
      password: 'pass123',
      name: 'New User',
      email: 'new@example.com'
    };
    const response = await api.register(userData);
    console.log('Register success:', response);
  } catch (error) {
    console.error('Register failed:', error.message);
  }
}

// Get Profile
async function getProfile() {
  try {
    const response = await api.getProfile();
    console.log('Profile:', response);
  } catch (error) {
    console.error('Get profile failed:', error.message);
  }
}

// Check if authenticated
if (api.isAuthenticated()) {
  console.log('User is logged in');
  const user = api.getCurrentUser();
  console.log('Current user:', user);
}
```

## Expected Responses

### Success Login Response:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "userid": "testuser",
      "name": "Test User",
      "email": "test@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2025-01-01T00:00:00Z"
    }
  }
}
```

### Error Response:
```json
{
  "error": "invalid credentials"
}
```

## Testing with Postman

1. Import collection or create new requests
2. Set base URL: `http://localhost:8080`

### Collection Structure:
```
GanttPro API
├── Health Check (GET /health)
├── Auth
│   ├── Register (POST /api/v1/auth/register)
│   ├── Login (POST /api/v1/auth/login)
│   └── Get Profile (GET /api/v1/auth/profile)
└── ...
```

### Setting up Bearer Token in Postman:
1. After login, copy the token from response
2. Go to Authorization tab
3. Select "Bearer Token"
4. Paste the token

## Common Issues

### CORS Error
- Make sure backend is running on port 8080
- Check ALLOWED_ORIGINS in .env includes your frontend URL

### Connection Refused
- Ensure backend server is running: `go run main.go`
- Check if port 8080 is available

### Invalid Token
- Token expires after 24 hours (configurable in .env)
- Login again to get new token

### Database Connection Error
- Verify database is running
- Check DB credentials in .env
- Ensure database exists
