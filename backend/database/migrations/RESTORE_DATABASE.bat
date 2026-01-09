@echo off
REM ============================================================
REM RESTORE DATABASE SCRIPT - GanttPro ERP
REM ============================================================
REM Script untuk restore database PostgreSQL dari backup
REM ============================================================

setlocal enabledelayedexpansion

echo ============================================================
echo    RESTORE DATABASE - GANTTPRO ERP
echo ============================================================
echo.

REM Konfigurasi
set DB_USER=postgres
set DB_NAME=ganttpro_db
set DB_HOST=localhost
set DB_PORT=5432

REM Tampilkan daftar backup yang tersedia
echo Backup files tersedia:
echo.
if exist "backups\*.sql" (
    dir /b /od backups\*.sql
    echo.
) else (
    echo Tidak ada backup ditemukan di folder backups\
    echo.
    pause
    exit /b 1
)

REM Minta input file backup
set /p BACKUP_FILE=Masukkan nama file backup (atau path lengkap):

REM Cek apakah file ada
if not exist "%BACKUP_FILE%" (
    if not exist "backups\%BACKUP_FILE%" (
        echo.
        echo [ERROR] File tidak ditemukan: %BACKUP_FILE%
        echo.
        pause
        exit /b 1
    ) else (
        set BACKUP_FILE=backups\%BACKUP_FILE%
    )
)

echo.
echo File backup: %BACKUP_FILE%
echo Target database: %DB_NAME%
echo.

REM Konfirmasi
set /p CONFIRM=PERINGATAN: Ini akan REPLACE semua data di database %DB_NAME%. Lanjutkan? (y/N):

if /i not "%CONFIRM%"=="y" (
    echo.
    echo Restore dibatalkan.
    pause
    exit /b 0
)

echo.

REM Minta password
set /p DB_PASSWORD=Masukkan password PostgreSQL:

REM Set password untuk psql
set PGPASSWORD=%DB_PASSWORD%

echo.
echo [1/4] Testing koneksi database...
psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d postgres -c "SELECT version();" > nul 2>&1

if %ERRORLEVEL% neq 0 (
    echo [ERROR] Gagal koneksi ke PostgreSQL!
    goto :cleanup
)
echo [OK] Koneksi berhasil.

echo.
echo [2/4] Dropping database lama...
REM Terminate semua koneksi
psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d postgres -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '%DB_NAME%' AND pid <> pg_backend_pid();" > nul 2>&1

REM Drop database
psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d postgres -c "DROP DATABASE IF EXISTS %DB_NAME%;" > nul 2>&1

if %ERRORLEVEL% neq 0 (
    echo [ERROR] Gagal drop database!
    goto :cleanup
)
echo [OK] Database lama di-drop.

echo.
echo [3/4] Creating database baru...
psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d postgres -c "CREATE DATABASE %DB_NAME%;" > nul 2>&1

if %ERRORLEVEL% neq 0 (
    echo [ERROR] Gagal create database!
    goto :cleanup
)
echo [OK] Database baru dibuat.

echo.
echo [4/4] Restoring data dari backup...

REM Cek jenis file backup
echo %BACKUP_FILE% | findstr /i ".dump" > nul
if %ERRORLEVEL% equ 0 (
    REM Custom format backup
    pg_restore -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d %DB_NAME% "%BACKUP_FILE%" 2>&1
) else (
    REM SQL format backup
    psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d %DB_NAME% < "%BACKUP_FILE%" 2>&1
)

if %ERRORLEVEL% neq 0 (
    echo.
    echo [WARNING] Restore selesai dengan beberapa error (mungkin normal).
) else (
    echo [OK] Restore berhasil.
)

echo.
echo ============================================================
echo    RESTORE SELESAI!
echo ============================================================
echo.
echo Database: %DB_NAME%
echo Source: %BACKUP_FILE%
echo.
echo Verifikasi data:
psql -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d %DB_NAME% -c "SELECT 'users' AS table_name, COUNT(*) AS records FROM users UNION ALL SELECT 'machines', COUNT(*) FROM machines UNION ALL SELECT 'ppic_schedules', COUNT(*) FROM ppic_schedules;" 2>&1
echo.
echo ============================================================

:cleanup
REM Clear password dari memory
set PGPASSWORD=
set DB_PASSWORD=

pause
