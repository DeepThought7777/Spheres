@echo off
REM Check if two arguments are provided
if "%~3"=="" (
    echo Error: Please provide three arguments.
    exit /b 1
)

REM Start the spheres.exe program with the provided arguments
start "%1" spheres.exe %2 %3

REM Exit the batch file
exit /b 0