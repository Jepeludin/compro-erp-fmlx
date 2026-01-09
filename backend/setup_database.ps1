#!/usr/bin/env pwsh
# Script untuk setup database PostgreSQL di Windows
# Jalankan dengan: .\setup_database.ps1

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "  Database Setup Script for GanttPro ERP" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# Check if PostgreSQL is installed
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    Write-Host "ERROR: PostgreSQL tidak ditemukan!" -ForegroundColor Red
    Write-Host "Pastikan PostgreSQL sudah terinstal dan psql ada di PATH" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Download PostgreSQL di: https://www.postgresql.org/download/" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ PostgreSQL ditemukan: $($psqlPath.Source)" -ForegroundColor Green
Write-Host ""

# Get database credentials
Write-Host "Masukkan kredensial PostgreSQL:" -ForegroundColor Yellow
$dbUser = Read-Host "PostgreSQL Username (default: postgres)"
if ([string]::IsNullOrWhiteSpace($dbUser)) {
    $dbUser = "postgres"
}

$dbPassword = Read-Host "PostgreSQL Password" -AsSecureString
$dbPasswordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($dbPassword))

$dbName = Read-Host "Database Name (default: ganttpro_db)"
if ([string]::IsNullOrWhiteSpace($dbName)) {
    $dbName = "ganttpro_db"
}

Write-Host ""
Write-Host "Konfigurasi:" -ForegroundColor Cyan
Write-Host "  User: $dbUser" -ForegroundColor White
Write-Host "  Database: $dbName" -ForegroundColor White
Write-Host ""

# Set PGPASSWORD environment variable
$env:PGPASSWORD = $dbPasswordPlain

# Test connection
Write-Host "Testing koneksi ke PostgreSQL..." -ForegroundColor Yellow
$testResult = psql -U $dbUser -d postgres -c "SELECT version();" 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Tidak bisa connect ke PostgreSQL!" -ForegroundColor Red
    Write-Host $testResult -ForegroundColor Red
    Write-Host ""
    Write-Host "Pastikan:" -ForegroundColor Yellow
    Write-Host "  1. PostgreSQL service sedang running" -ForegroundColor Yellow
    Write-Host "  2. Username dan password benar" -ForegroundColor Yellow
    Write-Host "  3. PostgreSQL berjalan di port 5432" -ForegroundColor Yellow
    exit 1
}
Write-Host "✓ Koneksi berhasil!" -ForegroundColor Green
Write-Host ""

# Check if database exists
Write-Host "Mengecek database '$dbName'..." -ForegroundColor Yellow
$dbExists = psql -U $dbUser -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$dbName'" 2>&1

if ($dbExists -eq "1") {
    Write-Host "⚠ Database '$dbName' sudah ada!" -ForegroundColor Yellow
    $response = Read-Host "Apakah Anda ingin DROP dan buat ulang? (yes/no)"
    
    if ($response -eq "yes" -or $response -eq "y") {
        Write-Host "Dropping database '$dbName'..." -ForegroundColor Yellow
        psql -U $dbUser -d postgres -c "DROP DATABASE $dbName;" 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ Database dropped" -ForegroundColor Green
        } else {
            Write-Host "ERROR: Gagal drop database" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "Setup dibatalkan." -ForegroundColor Yellow
        exit 0
    }
}

# Create database
Write-Host "Membuat database '$dbName'..." -ForegroundColor Yellow
psql -U $dbUser -d postgres -c "CREATE DATABASE $dbName;" 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Database berhasil dibuat" -ForegroundColor Green
} else {
    Write-Host "ERROR: Gagal membuat database" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Run migration script
$migrationFile = Join-Path $PSScriptRoot "database\migrations\00_complete_setup.sql"
if (-not (Test-Path $migrationFile)) {
    Write-Host "ERROR: File migration tidak ditemukan!" -ForegroundColor Red
    Write-Host "Path: $migrationFile" -ForegroundColor Red
    exit 1
}

Write-Host "Menjalankan migration script..." -ForegroundColor Yellow
Write-Host "File: $migrationFile" -ForegroundColor Gray
psql -U $dbUser -d $dbName -f $migrationFile 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Migration berhasil dijalankan" -ForegroundColor Green
} else {
    Write-Host "ERROR: Migration gagal" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Verify tables
Write-Host "Verifikasi table..." -ForegroundColor Yellow
$tableCount = psql -U $dbUser -d $dbName -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"
Write-Host "✓ Total $tableCount tables berhasil dibuat" -ForegroundColor Green

$userCount = psql -U $dbUser -d $dbName -tAc "SELECT COUNT(*) FROM users;"
Write-Host "✓ Total $userCount users" -ForegroundColor Green

$machineCount = psql -U $dbUser -d $dbName -tAc "SELECT COUNT(*) FROM machines;"
Write-Host "✓ Total $machineCount machines" -ForegroundColor Green
Write-Host ""

# Update .env file
Write-Host "Update file .env..." -ForegroundColor Yellow
$envFile = Join-Path $PSScriptRoot ".env"
$envContent = @"
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=$dbUser
DB_PASSWORD=$dbPasswordPlain
DB_NAME=$dbName
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Environment
ENVIRONMENT=development
"@

$envContent | Out-File -FilePath $envFile -Encoding UTF8
Write-Host "✓ File .env berhasil dibuat/diupdate" -ForegroundColor Green
Write-Host ""

# Clear password from environment
$env:PGPASSWORD = $null

# Done
Write-Host "=========================================" -ForegroundColor Green
Write-Host "  Setup Database BERHASIL!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Green
Write-Host ""
Write-Host "Default Login Credentials:" -ForegroundColor Cyan
Write-Host "  Email   : admin@example.com" -ForegroundColor White
Write-Host "  Username: admin" -ForegroundColor White
Write-Host "  Password: admin123" -ForegroundColor White
Write-Host ""
Write-Host "Langkah selanjutnya:" -ForegroundColor Yellow
Write-Host "  1. Jalankan backend: go run main.go" -ForegroundColor White
Write-Host "  2. Test API: http://localhost:8080/api/ppic/gantt-data" -ForegroundColor White
Write-Host ""
Write-Host "PENTING: Segera ganti password default setelah login!" -ForegroundColor Red
Write-Host ""
