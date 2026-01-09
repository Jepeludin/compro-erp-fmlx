@echo off
REM ============================================================
REM Script Setup Database Migration - GanttPro ERP (Batch)
REM ============================================================
REM Untuk Windows Command Prompt
REM ============================================================

echo ============================================================
echo    GANTTPRO ERP - DATABASE MIGRATION SETUP
echo ============================================================
echo.

REM Jalankan PowerShell script
powershell.exe -ExecutionPolicy Bypass -File "%~dp0setup_database_migration.ps1"

pause
