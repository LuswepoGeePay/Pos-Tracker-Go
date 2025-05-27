package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCleanedFilePathUsers(filePath string) string {
	// Replace backslashes with forward slashes for uniform handling
	normalizedPath := strings.ReplaceAll(filePath, "\\", "/")

	// Split the path into parts
	parts := strings.Split(normalizedPath, "/")

	startIndex := -1
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "Files" && parts[i+1] == "Users" {
			startIndex = i
			break
		}
	}

	// If the "Files" and "Users" sequence is not found, return the original file path
	if startIndex == -1 {
		return normalizedPath
	}

	// Construct the cleaned file path by joining the parts from the index after "Users"
	cleanedParts := parts[startIndex+2:]
	cleanedFilePath := strings.Join(cleanedParts, "/")
	return cleanedFilePath
}

func GetCleanedFilePathUsers1(filePath string) string {

	parts := strings.Split(filePath, "/")

	startIndex := -1
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "Files" && parts[i+1] == "Users" {
			startIndex = i
			break
		}
	}

	// If the consecutive "Files" and "Companies" are not found, return the original file path
	if startIndex == -1 {
		return filePath
	}

	// Construct the cleaned file path by joining the parts from the index after "Companies"
	cleanedParts := parts[startIndex+2:]
	cleanedFilePath := strings.Join(cleanedParts, "/")
	return cleanedFilePath
}

func GetCleanedFilePathFiles(filePath string, folderName string) string {

	parts := strings.Split(filePath, "/")

	startIndex := -1
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "Files" && parts[i+1] == folderName {
			startIndex = i
			break
		}
	}

	// If the consecutive "Files" and "Companies" are not found, return the original file path
	if startIndex == -1 {
		return filePath
	}

	// Construct the cleaned file path by joining the parts from the index after "Companies"
	cleanedParts := parts[startIndex+2:]
	cleanedFilePath := strings.Join(cleanedParts, "/")
	return cleanedFilePath
}
func CapitalizeError(msg string) error {
	return errors.New(strings.ToUpper(msg[:1]) + msg[1:])
}

func ParseDateTime(dateStr, timeStr string) (time.Time, error) {
	// Combine the date and time without appending Z
	combined := fmt.Sprintf("%sT%s", dateStr, timeStr)
	layout := "2006-01-02T15:04:05" // RFC3339 without timezone

	// Parse the combined date and time as local time
	return time.ParseInLocation(layout, combined, time.Local)
}

func SavePhoto(baseDir, userId, fileName string, photo []byte) (string, error) {
	userDir := filepath.Join(baseDir, userId, "Images")
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.MkdirAll(userDir, 0755); err != nil {
			return "", CapitalizeError("failed to create directory for user images")
		}
	}

	photoExt := filepath.Ext(fileName)
	photoFileName := uuid.New().String() + "_" + time.Now().Format("20060102_150405") + photoExt
	photoPath := filepath.Join(userDir, photoFileName)

	if err := os.WriteFile(photoPath, photo, 0644); err != nil {
		return "", CapitalizeError("failed to save photo")
	}
	return photoPath, nil
}

func SaveFile(baseDir, userId, fileName string, photo []byte) (string, error) {
	userDir := filepath.Join(baseDir, userId, "Files")
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.MkdirAll(userDir, 0755); err != nil {
			return "", CapitalizeError("failed to create directory for files")
		}
	}

	photoExt := filepath.Ext(fileName)
	photoFileName := uuid.New().String() + "_" + time.Now().Format("20060102_150405") + photoExt
	photoPath := filepath.Join(userDir, photoFileName)

	if err := os.WriteFile(photoPath, photo, 0644); err != nil {
		return "", CapitalizeError("failed to save photo")
	}
	return photoPath, nil
}

func ParseAndFormatDate(dateStr string) (string, error) {
	inputFormats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"January 2, 2006 15:04",
	}

	for _, format := range inputFormats {
		if parsedTime, err := time.Parse(format, dateStr); err == nil {
			return parsedTime.Format("2006-01-02 15:04:05"), nil
		}
	}

	return "", fmt.Errorf("invalid date-time format: %s", dateStr)
}

func CombineDateAndTime(dateStr, timeStr string) (string, error) {
	// Debugging: Input values
	fmt.Println("Debugging: Input dateStr:", dateStr)
	fmt.Println("Debugging: Input timeStr:", timeStr)

	// Case 1: Try parsing `dateStr` as a combined date-time string
	combinedFormats := []string{
		"January 02, 2006 15:04", // Example: January 09, 2025 07:53
		"2006-01-02 15:04:05",    // Example: 2025-01-01 05:53:00
	}

	for _, format := range combinedFormats {
		parsedDateTime, err := time.Parse(format, dateStr)
		if err == nil {
			fmt.Println("Debugging: Matched combined date-time format:", format)
			// Subtract 2 hours
			parsedDateTime = parsedDateTime.Add(-2 * time.Hour)
			return parsedDateTime.Format("2006-01-02 15:04:05"), nil
		}
		fmt.Println("Debugging: Failed to match combined format", format, "error:", err)
	}

	// Case 2: Parse `dateStr` as a standalone date
	dateOnlyFormats := []string{
		"January 02, 2006", // Example: January 09, 2025
		"2006-01-02",       // Example: 2025-01-01
	}

	var parsedDate time.Time
	var err error
	for _, format := range dateOnlyFormats {
		parsedDate, err = time.Parse(format, dateStr)
		if err == nil {
			fmt.Println("Debugging: Matched date-only format:", format)
			break
		}
		fmt.Println("Debugging: Failed to match date-only format", format, "error:", err)
	}

	if err != nil {
		return "", fmt.Errorf("failed to parse date: %v", err)
	}

	// Case 3: Combine with `timeStr` if provided
	if timeStr != "" {
		// Parse the time string
		parsedTime, err := time.Parse("15:04", timeStr)
		if err != nil {
			fmt.Println("Debugging: Failed to parse time:", err)
			return "", fmt.Errorf("failed to parse time: %v", err)
		}

		// Combine date and time
		combinedDateTime := time.Date(
			parsedDate.Year(),
			parsedDate.Month(),
			parsedDate.Day(),
			parsedTime.Hour(),
			parsedTime.Minute(),
			0, // seconds
			0, // nanoseconds
			parsedDate.Location(),
		)

		// Subtract 2 hours
		combinedDateTime = combinedDateTime.Add(-2 * time.Hour)

		fmt.Println("Debugging: Combined date and time (after subtracting 2 hours):", combinedDateTime)
		return combinedDateTime.Format("2006-01-02 15:04:05"), nil
	}

	// Default to midnight if no `timeStr` is provided
	fmt.Println("Debugging: Using default time (midnight)")
	parsedDate = parsedDate.Add(-2 * time.Hour) // Subtract 2 hours
	return parsedDate.Format("2006-01-02 15:04:05"), nil
}

func handleImageUpload(c *gin.Context, photoKey string) ([]byte, string, error) {

	photo, photoErr := c.FormFile(photoKey)

	if photoErr != nil {
		fmt.Printf("%s error: %v\n", photoKey, photoErr)
		return nil, "", nil
	}

	file, err := photo.Open()

	if err != nil {
		return nil, "", err
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		fmt.Printf("Error reading this: %s, filename: %s\n", photoKey, photo.Filename)
		return nil, "", err
	}

	return fileBytes, photo.Filename, nil
}

func ProcessPhotos(c *gin.Context) ([]byte, []byte, []byte, string, string, string) {

	photo1Bytes, photo1name, _ := handleImageUpload(c, "photo1")
	photo2Bytes, photo2name, _ := handleImageUpload(c, "photo2")
	photo3Bytes, photo3name, _ := handleImageUpload(c, "photo3")

	return photo1Bytes, photo2Bytes, photo3Bytes, photo1name, photo2name, photo3name
}

func ProcessPhoto(c *gin.Context, photoKey string) ([]byte, string) {
	photoBytes, photoName, _ := handleImageUpload(c, photoKey)
	return photoBytes, photoName
}

func GenerateSixDigitCode() string {

	rand.Seed(time.Now().UnixNano())

	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%d", code)
}
