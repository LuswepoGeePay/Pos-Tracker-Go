package pocketbase

import (
	"encoding/json"
	"fmt"
	"pos-master/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func HandleUpload(c *gin.Context, token string) (string, error) {

	if token == "" {
		return "", utils.CapitalizeError("unauthorization header is required")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}

	file, err := c.FormFile("file")
	if err != nil {
		return "", utils.CapitalizeError("no file uploaded")
	}

	openedFile, err := file.Open()
	if err != nil {
		return "", utils.CapitalizeError("Can't open file")
	}
	defer openedFile.Close()

	collectionName := "pos_files" // Your collection name
	fileFieldName := "file"       // Your file field name in the collection schema

	client := resty.New()
	endpoint := fmt.Sprintf("http://89.250.72.76:8090/api/collections/%s/records", collectionName)

	resp, err := client.R().
		SetHeader("Authorization", token).
		SetFileReader(fileFieldName, file.Filename, openedFile).
		Post(endpoint)

	if err != nil || resp.IsError() {
		return "", utils.CapitalizeError("Failed to upload file")

	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", utils.CapitalizeError("failed to parse response")
	}

	// Construct file URL from response
	recordID := result["id"].(string)
	fileName := result[fileFieldName].(string) // file field value in record (filename/path)
	fileURL := fmt.Sprintf("http://89.250.72.76:8090/api/files/%s/%s/%s", collectionName, recordID, fileName)

	return fileURL, nil
}
