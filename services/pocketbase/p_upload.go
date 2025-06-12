package pocketbase

import (
	"encoding/json"
	"fmt"
	"log"
	"pos-master/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func HandleUpload(c *gin.Context, token string, formKey string) (string, error) {

	if token == "" {
		return "", utils.CapitalizeError("unauthorization header is required")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	file, err := c.FormFile(formKey)
	if err != nil {
		return "", utils.CapitalizeError("no file uploaded")
	}

	openedFile, err := file.Open()
	if err != nil {
		return "", utils.CapitalizeError("Can't open file")
	}
	defer openedFile.Close()

	collectionName := "pos_master_files" // Your collection name
	fileFieldName := "file"              // Your file field name in the collection schema

	client := resty.New()
	endpoint := fmt.Sprintf("http://102.23.120.239:8090/api/collections/%s/records", collectionName)

	resp, err := client.R().
		SetHeader("Authorization", token).
		SetFileReader(fileFieldName, file.Filename, openedFile).
		Post(endpoint)

	if err != nil {
		log.Printf("Upload error (request issue): %v\n", err)
		return "", utils.CapitalizeError(fmt.Sprintf("Failed to upload file: %s", fmt.Sprintf("error: %v", err)))
	}

	if resp.IsError() {
		log.Printf("Upload error (response error): %s\n", resp.String()) // Log full response body
		return "", utils.CapitalizeError(fmt.Sprintf("Failed to upload file: %s", resp.String()))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", utils.CapitalizeError("failed to parse response")
	}

	// Construct file URL from response
	recordID := result["id"].(string)
	fileName := result[fileFieldName].(string) // file field value in record (filename/path)
	fileURL := fmt.Sprintf("http://102.23.120.239:8090/api/files/%s/%s/%s", collectionName, recordID, fileName)

	return fileURL, nil
}
