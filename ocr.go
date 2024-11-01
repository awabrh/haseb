package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/otiai10/gosseract/v2"
)

// OCRResult represents the structured output from OCR processing
type OCRResult struct {
	Text           string    `json:"text"`
	ProcessingTime float64   `json:"processing_time"`
	Timestamp      time.Time `json:"timestamp"`
	Error          string    `json:"error,omitempty"`
}

// OCRClient provides thread-safe OCR operations
type OCRClient struct {
	client *gosseract.Client
	mutex  sync.Mutex
}

// NewOCRClient creates a new OCR client
func NewOCRClient() (*OCRClient, error) {
	client := gosseract.NewClient()
	err := client.SetLanguage("eng", "ara")
	if err != nil {
		return nil, fmt.Errorf("failed to set language: %v", err)
	}

	return &OCRClient{
		client: client,
		mutex:  sync.Mutex{},
	}, nil
}

// ProcessImage performs OCR on the given image file
func (c *OCRClient) ProcessImage(imagePath string) (*OCRResult, error) {
	start := time.Now()

	// Lock to ensure thread safety
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Validate file exists and is an image
	if !isImageFile(imagePath) {
		return nil, fmt.Errorf("invalid image file: %s", imagePath)
	}

	// Set the image file
	if err := c.client.SetImage(imagePath); err != nil {
		return nil, fmt.Errorf("failed to set image: %v", err)
	}

	// Perform OCR
	text, err := c.client.Text()
	if err != nil {
		return nil, fmt.Errorf("OCR failed: %v", err)
	}

	processingTime := time.Since(start).Seconds()

	return &OCRResult{
		Text:           text,
		ProcessingTime: processingTime,
		Timestamp:      time.Now(),
	}, nil
}

// ProcessImageBytes performs OCR on image data in bytes
func (c *OCRClient) ProcessImageBytes(imageData []byte) (*OCRResult, error) {
	start := time.Now()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err := c.client.SetImageFromBytes(imageData); err != nil {
		return nil, fmt.Errorf("failed to set image from bytes: %v", err)
	}

	text, err := c.client.Text()
	if err != nil {
		return nil, fmt.Errorf("OCR failed: %v", err)
	}

	processingTime := time.Since(start).Seconds()

	return &OCRResult{
		Text:           text,
		ProcessingTime: processingTime,
		Timestamp:      time.Now(),
	}, nil
}

// isImageFile checks if the file is a valid image
func isImageFile(path string) bool {
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".bmp":  true,
		".tiff": true,
	}
	ext := filepath.Ext(path)
	return validExtensions[ext]
}

// Close frees up resources
func (c *OCRClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.client.Close()
}

// CheckTesseractInstallation verifies if tesseract is installed
func CheckTesseractInstallation() error {
	_, err := exec.LookPath("tesseract")
	if err != nil {
		return fmt.Errorf("tesseract is not installed: %v", err)
	}
	return nil
}
