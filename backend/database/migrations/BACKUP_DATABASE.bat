@echo off
REM ============================================================
REM BACKUP DATABASE SCRIPT - GanttPro ERP
REM ============================================================
REM Script untuk backup database PostgreSQL
REM ============================================================

setlocal enabledelayedexpansion

echo ============================================================
echo    BACKUP DATABASE - GANTTPRO ERP
echo ============================================================
echo.

REM Konfigurasi
set DB_USER=postgres
set DB_NAME=ganttpro_db
set DB_HOST=localhost
set DB_PORT=5432

REM Buat folder backup jika belum ada
if not exist "backups" mkdir backups

REM Generate timestamp untuk nama file
for /f "tokens=2-4 delims=/ " %%a in ('date /t') do (set mydate=%%c%%a%%b)
for /f "tokens=1-2 delims=/: " %%a in ('time /t') do (set mytime=%%a%%b)
set timestamp=%mydate%_%mytime%

set BACKUP_FILE=backups\ganttpro_backup_%timestamp%.sql
set CUSTOM_BACKUP=backups\ganttpro_backup_%timestamp%.dump

echo [1/3] Memulai backup database...
echo Database: %DB_NAME%
echo Backup file: %BACKUP_FILE%
echo.

REM Minta password
set /p DB_PASSWORD=Masukkan password PostgreSQL:

REM Set password untuk pg_dump
set PGPASSWORD=%DB_PASSWORD%

REM Backup SQL format (readable)
echo [2/3] Creating SQL backup...
pg_dump -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d %DB_NAME% > "%BACKUP_FILE%" 2>&1

if %ERRORLEVEL% neq 0 (
    echo.
    echo [ERROR] Backup gagal!
    echo Pastikan PostgreSQL terinstall dan kredensial benar.
    goto :cleanup
)

echo [OK] SQL backup berhasil: %BACKUP_FILE%
echo.

REM Backup custom format (compressed, untuk pg_restore)
echo [3/3] Creating compressed backup...
pg_dump -U %DB_USER% -h %DB_HOST% -p %DB_PORT% -d %DB_NAME% -F c -f "%CUSTOM_BACKUP%" 2>&1

if %ERRORLEVEL% neq 0 (
    echo [WARNING] Compressed backup gagal, tapi SQL backup sudah ada.
) else (
    echo [OK] Compressed backup berhasil: %CUSTOM_BACKUP%
)

echo.
echo ============================================================
echo    BACKUP SELESAI!
echo ============================================================
echo.
echo File backup:
dir /b backups\ganttpro_backup_%timestamp%.*
echo.
echo Untuk restore:
echo   SQL format: psql -U postgres -d ganttpro_db ^< %BACKUP_FILE%
echo   Custom format: pg_restore -U postgres -d ganttpro_db %CUSTOM_BACKUP%
echo.
echo ============================================================

:cleanup
REM Clear password dari memory
set PGPASSWORD=
set DB_PASSWORD=

pause
