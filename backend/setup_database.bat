@echo off
REM Simple batch script untuk setup database
REM Alternative jika PowerShell script tidak bisa dijalankan

echo =========================================
echo   Database Setup Script for GanttPro ERP
echo =========================================
echo.

REM Check if psql exists
where psql >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: PostgreSQL tidak ditemukan!
    echo Pastikan PostgreSQL sudah terinstal dan psql ada di PATH
    echo.
    echo Download PostgreSQL di: https://www.postgresql.org/download/
    pause
    exit /b 1
)

echo [OK] PostgreSQL ditemukan
echo.

REM Get database credentials
set /p DB_USER="PostgreSQL Username (default: postgres): "
if "%DB_USER%"=="" set DB_USER=postgres

set /p DB_PASSWORD="PostgreSQL Password: "

set /p DB_NAME="Database Name (default: ganttpro_db): "
if "%DB_NAME%"=="" set DB_NAME=ganttpro_db

echo.
echo Konfigurasi:
echo   User: %DB_USER%
echo   Database: %DB_NAME%
echo.

REM Set password environment variable
set PGPASSWORD=%DB_PASSWORD%

REM Test connection
echo Testing koneksi ke PostgreSQL...
psql -U %DB_USER% -d postgres -c "SELECT version();" >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Tidak bisa connect ke PostgreSQL!
    echo.
    echo Pastikan:
    echo   1. PostgreSQL service sedang running
    echo   2. Username dan password benar
    echo   3. PostgreSQL berjalan di port 5432
    pause
    exit /b 1
)
echo [OK] Koneksi berhasil!
echo.

REM Check if database exists
echo Mengecek database '%DB_NAME%'...
psql -U %DB_USER% -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='%DB_NAME%'" > temp_check.txt
set /p DB_EXISTS=<temp_check.txt
del temp_check.txt

if "%DB_EXISTS%"=="1" (
    echo PERHATIAN: Database '%DB_NAME%' sudah ada!
    set /p CONFIRM="Apakah Anda ingin DROP dan buat ulang? (yes/no): "
    
    if /i "!CONFIRM!"=="yes" (
        echo Dropping database '%DB_NAME%'...
        psql -U %DB_USER% -d postgres -c "DROP DATABASE %DB_NAME%;" >nul 2>nul
        if !errorlevel! equ 0 (
            echo [OK] Database dropped
        ) else (
            echo ERROR: Gagal drop database
            pause
            exit /b 1
        )
    ) else (
        echo Setup dibatalkan.
        pause
        exit /b 0
    )
)

REM Create database
echo Membuat database '%DB_NAME%'...
psql -U %DB_USER% -d postgres -c "CREATE DATABASE %DB_NAME%;" >nul 2>nul
if %errorlevel% equ 0 (
    echo [OK] Database berhasil dibuat
) else (
    echo ERROR: Gagal membuat database
    pause
    exit /b 1
)
echo.

REM Run migration
echo Menjalankan migration script...
psql -U %DB_USER% -d %DB_NAME% -f database\migrations\00_complete_setup.sql >nul 2>nul
if %errorlevel% equ 0 (
    echo [OK] Migration berhasil dijalankan
) else (
    echo ERROR: Migration gagal
    pause
    exit /b 1
)
echo.

REM Verify tables
echo Verifikasi table...
for /f %%i in ('psql -U %DB_USER% -d %DB_NAME% -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';"') do set TABLE_COUNT=%%i
echo [OK] Total %TABLE_COUNT% tables berhasil dibuat

for /f %%i in ('psql -U %DB_USER% -d %DB_NAME% -tAc "SELECT COUNT(*) FROM users;"') do set USER_COUNT=%%i
echo [OK] Total %USER_COUNT% users

for /f %%i in ('psql -U %DB_USER% -d %DB_NAME% -tAc "SELECT COUNT(*) FROM machines;"') do set MACHINE_COUNT=%%i
echo [OK] Total %MACHINE_COUNT% machines
echo.

REM Create .env file
echo Update file .env...
(
echo # Database Configuration
echo DB_HOST=localhost
echo DB_PORT=5432
echo DB_USER=%DB_USER%
echo DB_PASSWORD=%DB_PASSWORD%
echo DB_NAME=%DB_NAME%
echo DB_SSLMODE=disable
echo.
echo # Server Configuration
echo PORT=8080
echo.
echo # JWT Configuration
echo JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
echo JWT_EXPIRY=24h
echo.
echo # Environment
echo ENVIRONMENT=development
) > .env
echo [OK] File .env berhasil dibuat/diupdate
echo.

REM Clear password
set PGPASSWORD=

REM Done
echo =========================================
echo   Setup Database BERHASIL!
echo =========================================
echo.
echo Default Login Credentials:
echo   Email   : admin@example.com
echo   Username: admin
echo   Password: admin123
echo.
echo Langkah selanjutnya:
echo   1. Jalankan backend: go run main.go
echo   2. Test API: http://localhost:8080/api/ppic/gantt-data
echo.
echo PENTING: Segera ganti password default setelah login!
echo.
pause
