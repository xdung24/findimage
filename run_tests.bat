@echo off
REM Test runner script for FindImg

echo FindImg Test Runner
echo ===================

echo.
echo 1. Running unit tests...
echo ------------------------
go test -v

echo.
echo 2. Running benchmarks...
echo -----------------------
go test -bench=.

echo.
echo 3. Running tests with coverage...
echo --------------------------------
go test -cover

echo.
echo 4. Building the application...
echo -----------------------------
go build -o findimg.exe
if errorlevel 1 (
    echo Build failed!
    exit /b 1
) else (
    echo Build successful: findimg.exe created
)

echo.
echo 5. Testing command line help...
echo ------------------------------
findimg.exe -h

echo.
echo Tests completed!
echo.
echo To run individual test commands:
echo   go test -v -run TestFindImageBasic
echo   go test -v -run TestRandomSubimage  
echo   go test -bench=BenchmarkFindImageSmall
echo   go test -cover
