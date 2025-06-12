package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"math/rand"
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

func FormatError(msg string, err error) string {
	return fmt.Sprintf("%s:%v", msg, err)
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

func GenerateSixDigitCode() string {

	rand.Seed(time.Now().UnixNano())

	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%d", code)
}
