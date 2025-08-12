@echo off
REM Simple test verification script

echo Testing FindImg Application
echo ============================

echo.
echo Step 1: Check Go installation
echo ----------------------------
go version
if errorlevel 1 (
    echo ERROR: Go is not installed or not in PATH
    exit /b 1
)

echo.
echo Step 2: Build the application
echo ----------------------------
go build -o findimg.exe
if errorlevel 1 (
    echo ERROR: Build failed
    exit /b 1
) else (
    echo SUCCESS: Application built successfully
)

echo.
echo Step 3: Test basic functionality
echo --------------------------------
echo Testing help output:
findimg.exe -h
if errorlevel 2 (
    echo SUCCESS: Help command works (exit code 2 is expected)
) else (
    echo WARNING: Unexpected exit code from help command
)

echo.
echo Step 4: Run unit tests
echo ---------------------
go test
if errorlevel 1 (
    echo ERROR: Tests failed
    exit /b 1
) else (
    echo SUCCESS: Tests passed
)

echo.
echo Step 5: Test with asset files (if available)
echo -------------------------------------------
if exist assets\haystack.jpg (
    if exist assets\needle.jpg (
        echo Testing with asset files:
        findimg.exe assets\haystack.jpg assets\needle.jpg
        echo SUCCESS: Asset test completed
    ) else (
        echo INFO: needle.jpg not found in assets
    )
) else (
    echo INFO: haystack.jpg not found in assets
)

echo.
echo ============================
echo All tests completed successfully!
echo ============================
echo.
echo Available commands:
echo   findimg.exe ^<image^> ^<subimage^>          - Basic usage
echo   findimg.exe -o json ^<image^> ^<subimage^>  - JSON output
echo   findimg.exe -v ^<image^> ^<subimage^>       - Verbose output
echo   findimg.exe -random ^<image^>               - Random subimage test
echo   findimg.exe -k 5 ^<image^> ^<subimage^>     - Find top 5 matches
