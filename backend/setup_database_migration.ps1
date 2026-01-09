# ============================================================
# Script Setup Database Migration - GanttPro ERP
# ============================================================
# Script ini akan:
# 1. Test koneksi PostgreSQL
# 2. Buat database ganttpro_db
# 3. Jalankan migration lengkap
# 4. Verifikasi instalasi
# ============================================================

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "   GANTTPRO ERP - DATABASE MIGRATION SETUP" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

# ============================================================
# KONFIGURASI
# ============================================================

$DB_USER = "postgres"
$DB_NAME = "ganttpro_db"
$DB_PORT = "5432"
$DB_HOST = "localhost"
$MIGRATION_FILE = "database\migrations\MIGRATION_COMPLETE.sql"

# ============================================================
# FUNGSI HELPER
# ============================================================

function Test-Command {
    param($Command)
    try {
        if (Get-Command $Command -ErrorAction SilentlyContinue) {
            return $true
        }
        return $false
    }
    catch {
        return $false
    }
}

function Write-Step {
    param($Message)
    Write-Host "`n>>> $Message" -ForegroundColor Yellow
}

function Write-Success {
    param($Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Error {
    param($Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

# ============================================================
# STEP 1: CEK PREREQUISITES
# ============================================================

Write-Step "Mengecek prerequisites..."

# Cek PostgreSQL
if (-not (Test-Command "psql")) {
    Write-Error "PostgreSQL tidak ditemukan!"
    Write-Host "`nSilakan install PostgreSQL terlebih dahulu:" -ForegroundColor Yellow
    Write-Host "https://www.postgresql.org/download/windows/" -ForegroundColor Cyan
    Write-Host "`nSetelah install, tambahkan PostgreSQL ke PATH:" -ForegroundColor Yellow
    Write-Host "C:\Program Files\PostgreSQL\15\bin" -ForegroundColor Cyan
    exit 1
}
Write-Success "PostgreSQL terinstall"

# Cek Go
if (-not (Test-Command "go")) {
    Write-Error "Go tidak ditemukan!"
    Write-Host "`nSilakan install Go terlebih dahulu:" -ForegroundColor Yellow
    Write-Host "https://go.dev/dl/" -ForegroundColor Cyan
    exit 1
}
Write-Success "Go terinstall"

# Cek migration file
if (-not (Test-Path $MIGRATION_FILE)) {
    Write-Error "File migration tidak ditemukan: $MIGRATION_FILE"
    exit 1
}
Write-Success "File migration ditemukan"

# ============================================================
# STEP 2: INPUT PASSWORD
# ============================================================

Write-Step "Masukkan kredensial PostgreSQL"
Write-Host "Database User: $DB_USER" -ForegroundColor Cyan

# Minta password
$DB_PASSWORD_SECURE = Read-Host "Password PostgreSQL" -AsSecureString
$DB_PASSWORD = [Runtime.InteropServices.Marshal]::PtrToStringAuto(
    [Runtime.InteropServices.Marshal]::SecureStringToBSTR($DB_PASSWORD_SECURE)
)

# Set environment variable untuk psql
$env:PGPASSWORD = $DB_PASSWORD

# ============================================================
# STEP 3: TEST KONEKSI POSTGRESQL
# ============================================================

Write-Step "Testing koneksi ke PostgreSQL..."

$testConnection = psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c "SELECT version();" 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Error "Gagal koneksi ke PostgreSQL!"
    Write-Host "`nPesan error: $testConnection" -ForegroundColor Red
    Write-Host "`nPastikan:" -ForegroundColor Yellow
    Write-Host "1. PostgreSQL service sedang running" -ForegroundColor Cyan
    Write-Host "2. Username dan password benar" -ForegroundColor Cyan
    Write-Host "3. Port 5432 tidak diblokir firewall" -ForegroundColor Cyan
    $env:PGPASSWORD = $null
    exit 1
}

Write-Success "Koneksi ke PostgreSQL berhasil!"
Write-Host "   PostgreSQL Version: $($testConnection -join ' ' | Select-String -Pattern 'PostgreSQL \d+\.\d+')" -ForegroundColor Gray

# ============================================================
# STEP 4: CEK & BUAT DATABASE
# ============================================================

Write-Step "Mengecek database '$DB_NAME'..."

# Cek apakah database sudah ada
$dbExists = psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -t -c "SELECT 1 FROM pg_database WHERE datname='$DB_NAME';" 2>&1

if ($dbExists -match "1") {
    Write-Host "   Database '$DB_NAME' sudah ada." -ForegroundColor Yellow

    $overwrite = Read-Host "   Apakah Anda ingin DROP dan buat ulang? (y/N)"

    if ($overwrite -eq "y" -or $overwrite -eq "Y") {
        Write-Host "   Dropping database..." -ForegroundColor Yellow

        # Terminate existing connections
        psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c @"
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = '$DB_NAME'
  AND pid <> pg_backend_pid();
"@ | Out-Null

        # Drop database
        psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;" | Out-Null

        if ($LASTEXITCODE -ne 0) {
            Write-Error "Gagal drop database!"
            $env:PGPASSWORD = $null
            exit 1
        }

        # Create new database
        psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c "CREATE DATABASE $DB_NAME;" | Out-Null

        if ($LASTEXITCODE -ne 0) {
            Write-Error "Gagal membuat database!"
            $env:PGPASSWORD = $null
            exit 1
        }

        Write-Success "Database '$DB_NAME' berhasil dibuat ulang"
    } else {
        Write-Host "   Menggunakan database yang sudah ada." -ForegroundColor Cyan
    }
} else {
    Write-Host "   Membuat database baru '$DB_NAME'..." -ForegroundColor Cyan

    psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d postgres -c "CREATE DATABASE $DB_NAME;" | Out-Null

    if ($LASTEXITCODE -ne 0) {
        Write-Error "Gagal membuat database!"
        $env:PGPASSWORD = $null
        exit 1
    }

    Write-Success "Database '$DB_NAME' berhasil dibuat"
}

# ============================================================
# STEP 5: JALANKAN MIGRATION
# ============================================================

Write-Step "Menjalankan migration script..."
Write-Host "   File: $MIGRATION_FILE" -ForegroundColor Gray

$migrationOutput = psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -f $MIGRATION_FILE 2>&1

if ($LASTEXITCODE -ne 0) {
    Write-Error "Migration gagal!"
    Write-Host "`nOutput:" -ForegroundColor Red
    Write-Host $migrationOutput
    $env:PGPASSWORD = $null
    exit 1
}

Write-Success "Migration berhasil dijalankan!"

# Show migration summary
$migrationOutput | Select-String -Pattern "NOTICE:" | ForEach-Object {
    $line = $_.ToString() -replace "NOTICE:\s+", ""
    if ($line -match "====") {
        Write-Host "   $line" -ForegroundColor DarkGray
    } else {
        Write-Host "   $line" -ForegroundColor Cyan
    }
}

# ============================================================
# STEP 6: VERIFIKASI DATABASE
# ============================================================

Write-Step "Memverifikasi instalasi..."

# Hitung jumlah tabel
$tableCount = psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -t -c @"
SELECT COUNT(*)
FROM information_schema.tables
WHERE table_schema = 'public' AND table_type = 'BASE TABLE';
"@ 2>&1

# Hitung jumlah users
$userCount = psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -t -c "SELECT COUNT(*) FROM users;" 2>&1

# Hitung jumlah machines
$machineCount = psql -U $DB_USER -h $DB_HOST -p $DB_PORT -d $DB_NAME -t -c "SELECT COUNT(*) FROM machines;" 2>&1

Write-Success "Verifikasi selesai"
Write-Host "   Total Tables: $($tableCount.Trim())" -ForegroundColor Cyan
Write-Host "   Total Users: $($userCount.Trim())" -ForegroundColor Cyan
Write-Host "   Total Machines: $($machineCount.Trim())" -ForegroundColor Cyan

# ============================================================
# STEP 7: UPDATE .ENV FILE
# ============================================================

Write-Step "Membuat/update file .env..."

$envContent = @"
# Server Configuration
PORT=8080

# Database Configuration
DB_DRIVER=postgres
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME

# JWT Secret (GANTI dengan random string untuk production!)
JWT_SECRET=ganttpro-secret-key-change-in-production-$(Get-Random)

# Environment
ENV=development

# Upload Configuration
UPLOAD_DIR=./uploads
MAX_UPLOAD_SIZE=10485760

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
"@

$envPath = ".env"
$envContent | Out-File -FilePath $envPath -Encoding utf8

Write-Success "File .env berhasil dibuat/update"

# Clear password dari memory
$env:PGPASSWORD = $null

# ============================================================
# STEP 8: TAMPILKAN INFORMASI
# ============================================================

Write-Host ""
Write-Host "============================================================" -ForegroundColor Green
Write-Host "   SETUP SELESAI!" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Green
Write-Host ""

Write-Host "📋 Default User Credentials:" -ForegroundColor Cyan
Write-Host "   ┌─────────────────────────────────────────────────┐" -ForegroundColor Gray
Write-Host "   │ Admin    | admin     | admin123               │" -ForegroundColor White
Write-Host "   │ PPIC     | ppic      | ppic123                │" -ForegroundColor White
Write-Host "   │ PEM      | pem       | pem123                 │" -ForegroundColor White
Write-Host "   │ Operator | operator  | operator123            │" -ForegroundColor White
Write-Host "   └─────────────────────────────────────────────────┘" -ForegroundColor Gray
Write-Host ""

Write-Host "📁 Database Info:" -ForegroundColor Cyan
Write-Host "   Host: $DB_HOST" -ForegroundColor White
Write-Host "   Port: $DB_PORT" -ForegroundColor White
Write-Host "   Database: $DB_NAME" -ForegroundColor White
Write-Host "   Tables: $($tableCount.Trim())" -ForegroundColor White
Write-Host ""

Write-Host "🚀 Next Steps:" -ForegroundColor Cyan
Write-Host "   1. Jalankan backend:" -ForegroundColor White
Write-Host "      go run main.go" -ForegroundColor Yellow
Write-Host ""
Write-Host "   2. Test API:" -ForegroundColor White
Write-Host "      http://localhost:8080/health" -ForegroundColor Yellow
Write-Host ""
Write-Host "   3. Jalankan frontend (di folder frontend):" -ForegroundColor White
Write-Host "      npm install" -ForegroundColor Yellow
Write-Host "      npm run dev" -ForegroundColor Yellow
Write-Host ""

Write-Host "⚠️  PENTING: Ganti password default setelah login pertama kali!" -ForegroundColor Red
Write-Host ""
Write-Host "============================================================" -ForegroundColor Green
Write-Host ""

# ============================================================
# STEP 9: TANYA JALANKAN BACKEND
# ============================================================

$runBackend = Read-Host "Apakah Anda ingin menjalankan backend sekarang? (y/N)"

if ($runBackend -eq "y" -or $runBackend -eq "Y") {
    Write-Step "Menjalankan backend server..."
    Write-Host "   Tekan Ctrl+C untuk stop server" -ForegroundColor Yellow
    Write-Host ""

    # Download dependencies jika belum
    if (-not (Test-Path "go.sum")) {
        Write-Host "   Downloading Go dependencies..." -ForegroundColor Cyan
        go mod download
    }

    # Jalankan backend
    go run main.go
} else {
    Write-Host ""
    Write-Host "Terima kasih! Jalankan 'go run main.go' saat Anda siap." -ForegroundColor Cyan
    Write-Host ""
}
