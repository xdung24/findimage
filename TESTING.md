# FindImg Testing Guide

This guide provides comprehensive examples and tests for the FindImg image template matching tool.

## Quick Start

### 1. Build and Test
```cmd
# Build the application
go build -o findimg.exe

# Run unit tests
go test -v

# Run verification script
test_verification.bat
```

### 2. Basic Usage Examples
```cmd
# Find needle in haystack (if asset files exist)
findimg.exe assets\haystack.jpg assets\needle.jpg

# Get JSON output
findimg.exe -o json assets\haystack.jpg assets\needle.jpg

# Verbose output
findimg.exe -v assets\haystack.jpg assets\needle.jpg

# Random subimage test
findimg.exe -random assets\haystack.jpg

# Find top 3 matches
findimg.exe -k 3 assets\haystack.jpg assets\needle.jpg
```

## Test Files Created

### Unit Tests (`main_test.go`)
- **TestFindImageBasic**: Tests basic image template matching
- **TestRandomSubimage**: Tests random subimage generation
- **TestMatchMethods**: Tests Match struct methods (CenterX, CenterY, Scale)
- **TestSumOfAbsDiff**: Tests image difference calculation
- **TestIntegrationWithFiles**: Integration test using actual asset files
- **BenchmarkFindImageSmall/Large**: Performance benchmarks

### Example Scripts
- **`examples/run_examples.bat`**: Command-line usage examples
- **`examples/example_usage.go`**: Programmatic usage examples
- **`test_verification.bat`**: Complete test verification script
- **`run_tests.bat`**: Comprehensive test runner

## Running Tests

### Unit Tests
```cmd
# Run all tests
go test -v

# Run specific test
go test -v -run TestFindImageBasic

# Run with coverage
go test -cover

# Run benchmarks
go test -bench=.
```

### Integration Tests
```cmd
# Test with actual files (requires assets)
go test -v -run TestIntegrationWithFiles

# Command line examples
examples\run_examples.bat
```

### Verification
```cmd
# Complete verification
test_verification.bat
```

## Creating Your Own Tests

### 1. Unit Test Template
```go
func TestYourFunction(t *testing.T) {
    // Create test images
    haystack := createTestImageWithSquare(100, 100, 20,
        color.RGBA{255, 255, 255, 255}, // white background
        color.RGBA{255, 0, 0, 255})     // red square
    
    needle := createTestImage(20, 20, color.RGBA{255, 0, 0, 255})
    
    // Test the function
    opts := Opts{k: 1, verbose: false}
    matches := findImage(haystack, needle, opts)
    
    // Verify results
    if len(matches) == 0 {
        t.Fatal("Expected to find matches")
    }
}
```

### 2. Command Line Test
```cmd
# Create your own test images and run
findimg.exe your_haystack.jpg your_needle.jpg
```

### 3. Benchmark Template
```go
func BenchmarkYourFunction(b *testing.B) {
    // Setup test data
    haystack := createTestImage(500, 500, color.RGBA{255, 255, 255, 255})
    needle := createTestImage(50, 50, color.RGBA{255, 0, 0, 255})
    opts := Opts{k: 1}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        findImage(haystack, needle, opts)
    }
}
```

## Expected Results

### Successful Test Output
```
=== RUN   TestFindImageBasic
--- PASS: TestFindImageBasic (0.01s)
=== RUN   TestRandomSubimage
--- PASS: TestRandomSubimage (0.00s)
=== RUN   TestMatchMethods
--- PASS: TestMatchMethods (0.00s)
=== RUN   TestSumOfAbsDiff
--- PASS: TestSumOfAbsDiff (0.00s)
=== RUN   TestIntegrationWithFiles
--- SKIP: TestIntegrationWithFiles (0.00s)
    main_test.go:XXX: Skipping integration test: assets/haystack.jpg not found
PASS
ok      findimg 0.XXXs
```

### Command Line Output Format
```
# Text format (default)
0.950000   40   40   20   20   50   50

# JSON format (-o json)
{
  "matches": [
    {
      "bounds": {"x": 40, "y": 40, "width": 20, "height": 20},
      "confident": 0.95,
      "centerX": 50,
      "centerY": 50
    }
  ]
}
```

## Performance Testing

### Small Images (100x100)
```cmd
go test -bench=BenchmarkFindImageSmall
```

### Large Images (500x500)
```cmd
go test -bench=BenchmarkFindImageLarge
```

### Memory Usage
```cmd
go test -bench=. -benchmem
```

## Troubleshooting Tests

### Common Issues

1. **Build Failures**
   - Ensure Go is properly installed
   - Check that all dependencies are available
   - Run `go mod tidy` if needed

2. **Test Failures**
   - Check that test images are created correctly
   - Verify expected vs actual results
   - Use `-v` flag for detailed output

3. **Asset File Tests Skipped**
   - Normal if assets/haystack.jpg and assets/needle.jpg don't exist
   - Create your own test images or obtain the asset files

4. **Performance Issues**
   - Use smaller test images for faster tests
   - Adjust timeout values if needed
   - Consider running benchmarks separately

### Debug Tips
```cmd
# Verbose test output
go test -v

# Run only failing test
go test -v -run TestFailingFunction

# Show all test output including prints
go test -v -args -test.v

# Trace test execution
go test -trace=trace.out
```

## Continuous Integration

For automated testing, use:
```cmd
# Basic CI pipeline
go build
go test -short
go test -race -short
```

## Test Data

The tests create synthetic images with known patterns:
- **White background with red squares**
- **Colored circles and stripes**
- **Various sizes and positions**

This ensures predictable test results independent of external image files.
