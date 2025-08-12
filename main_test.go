package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// createTestImage creates a simple test image with a specific pattern
func createTestImage(width, height int, fillColor color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, fillColor)
		}
	}
	return img
}

// createTestImageWithSquare creates a test image with a colored square in the center
func createTestImageWithSquare(width, height, squareSize int, bgColor, squareColor color.RGBA) *image.RGBA {
	img := createTestImage(width, height, bgColor)

	startX := (width - squareSize) / 2
	startY := (height - squareSize) / 2

	for y := startY; y < startY+squareSize && y < height; y++ {
		for x := startX; x < startX+squareSize && x < width; x++ {
			img.Set(x, y, squareColor)
		}
	}
	return img
}

// saveTestImage saves an image to a file for testing
func saveTestImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func TestFindImageBasic(t *testing.T) {
	// Create a test haystack image (100x100 with a red square in the center)
	haystack := createTestImageWithSquare(100, 100, 20,
		color.RGBA{255, 255, 255, 255}, // white background
		color.RGBA{255, 0, 0, 255})     // red square

	// Create a test needle image (20x20 red square)
	needle := createTestImage(20, 20, color.RGBA{255, 0, 0, 255})

	// Test finding the image
	opts := Opts{
		k:       1,
		verbose: false,
	}

	matches := findImage(haystack, needle, opts)

	if len(matches) == 0 {
		t.Fatal("Expected to find at least one match")
	}

	// The match should be around the center of the image
	match := matches[0]
	expectedX := 40 // (100-20)/2
	expectedY := 40 // (100-20)/2

	tolerance := 5
	if abs(match.Bounds.Min.X-expectedX) > tolerance ||
		abs(match.Bounds.Min.Y-expectedY) > tolerance {
		t.Errorf("Match position (%d, %d) not near expected position (%d, %d)",
			match.Bounds.Min.X, match.Bounds.Min.Y, expectedX, expectedY)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestRandomSubimage(t *testing.T) {
	// Create a test image
	img := createTestImage(100, 100, color.RGBA{255, 0, 0, 255})

	// Test random subimage generation
	subimg := randomSubimage(img)

	if subimg == nil {
		t.Fatal("randomSubimage returned nil")
	}

	bounds := subimg.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Fatal("randomSubimage returned invalid dimensions")
	}
}

func TestMatchMethods(t *testing.T) {
	match := Match{
		Bounds:    image.Rect(10, 20, 30, 40),
		Confident: 0.95,
	}

	// Test CenterX
	expectedCenterX := 15 // (10+20)/2
	if match.CenterX() != expectedCenterX {
		t.Errorf("CenterX() = %d, expected %d", match.CenterX(), expectedCenterX)
	}

	// Test CenterY
	expectedCenterY := 20 // (20+20)/2
	if match.CenterY() != expectedCenterY {
		t.Errorf("CenterY() = %d, expected %d", match.CenterY(), expectedCenterY)
	}

	// Test Scale
	scaled := match.Scale(2.0)
	expectedScaledBounds := image.Rect(20, 40, 60, 80)
	if scaled.Bounds != expectedScaledBounds {
		t.Errorf("Scale(2.0) bounds = %v, expected %v", scaled.Bounds, expectedScaledBounds)
	}
}

func TestSumOfAbsDiff(t *testing.T) {
	// Create two identical small images
	img1 := createTestImage(10, 10, color.RGBA{100, 100, 100, 255})
	img2 := createTestImage(5, 5, color.RGBA{100, 100, 100, 255})

	// When images are identical at the given position, difference should be 0
	diff := sumOfAbsDiff(img1, 0, 0, img2)
	if diff != 0 {
		t.Errorf("sumOfAbsDiff for identical images = %d, expected 0", diff)
	}

	// Test with different colors
	img3 := createTestImage(5, 5, color.RGBA{200, 200, 200, 255})
	diff2 := sumOfAbsDiff(img1, 0, 0, img3)
	if diff2 == 0 {
		t.Error("sumOfAbsDiff for different images should not be 0")
	}
}

// Example test that demonstrates usage patterns
func ExampleMain() {
	// This example shows how to use the findImage function programmatically

	// Create a haystack image (larger image to search in)
	haystack := createTestImageWithSquare(200, 200, 40,
		color.RGBA{255, 255, 255, 255}, // white background
		color.RGBA{0, 255, 0, 255})     // green square

	// Create a needle image (pattern to find)
	needle := createTestImage(40, 40, color.RGBA{0, 255, 0, 255}) // green square

	// Set up options
	opts := Opts{
		k:           3, // Find top 3 matches
		verbose:     false,
		imgMinWidth: 0,
		imgMaxWidth: 0,
		subMinArea:  0,
		subMaxDiv:   0,
	}

	// Find matches
	matches := findImage(haystack, needle, opts)

	// Process results
	for i, match := range matches {
		_ = i               // avoid unused variable in example
		_ = match.Confident // confidence score
		_ = match.CenterX() // center X coordinate
		_ = match.CenterY() // center Y coordinate
		// In a real application, you would use these values
	}
}

// Benchmark tests
func BenchmarkFindImageSmall(b *testing.B) {
	haystack := createTestImageWithSquare(100, 100, 20,
		color.RGBA{255, 255, 255, 255},
		color.RGBA{255, 0, 0, 255})
	needle := createTestImage(20, 20, color.RGBA{255, 0, 0, 255})
	opts := Opts{k: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findImage(haystack, needle, opts)
	}
}

func BenchmarkFindImageLarge(b *testing.B) {
	haystack := createTestImageWithSquare(500, 500, 50,
		color.RGBA{255, 255, 255, 255},
		color.RGBA{0, 0, 255, 255})
	needle := createTestImage(50, 50, color.RGBA{0, 0, 255, 255})
	opts := Opts{k: 1}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		findImage(haystack, needle, opts)
	}
}

// Integration test using actual files
func TestIntegrationWithFiles(t *testing.T) {
	// Skip this test if the assets don't exist
	if _, err := os.Stat("assets/haystack.jpg"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: assets/haystack.jpg not found")
	}
	if _, err := os.Stat("assets/needle.jpg"); os.IsNotExist(err) {
		t.Skip("Skipping integration test: assets/needle.jpg not found")
	}

	// Open the actual test images
	haystack, err := openImage("assets/haystack.jpg")
	if err != nil {
		t.Fatalf("Failed to open haystack image: %v", err)
	}

	needle, err := openImage("assets/needle.jpg")
	if err != nil {
		t.Fatalf("Failed to open needle image: %v", err)
	}

	// Test finding the needle in the haystack
	opts := Opts{k: 5}
	matches := findImage(haystack, needle, opts)

	if len(matches) == 0 {
		t.Error("Expected to find matches in the integration test")
	}

	// Verify that matches have reasonable confidence scores
	for _, match := range matches {
		if match.Confident < 0 || match.Confident > 1 {
			t.Errorf("Match confidence %f is out of range [0, 1]", match.Confident)
		}
	}
}
