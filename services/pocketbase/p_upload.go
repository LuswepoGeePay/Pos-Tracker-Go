package pocketbase

import (
	"encoding/json"
	"fmt"
	"os"
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

	collectionName := os.Getenv("POCKETBASE_COLLECTION")
	if collectionName == "" {
		collectionName = "pos_master_files"
	}
	fileFieldName := "file"

	client := resty.New()
	endpoint := fmt.Sprintf("%s/api/collections/%s/records", BaseURL(), collectionName)

	resp, err := client.R().
		SetHeader("Authorization", token).
		SetFileReader(fileFieldName, file.Filename, openedFile).
		Post(endpoint)

	if err != nil {
		utils.Error("pocketbase upload request failed", "error", err.Error(), "endpoint", endpoint)
		return "", utils.CapitalizeError(fmt.Sprintf("Failed to upload file: %v", err))
	}

	if resp.IsError() {
		utils.Error("pocketbase upload rejected",
			"status", resp.StatusCode(),
			"body", resp.String(),
			"filename", file.Filename,
		)
		return "", utils.CapitalizeError(fmt.Sprintf("Failed to upload file: %s", resp.String()))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		utils.Error("failed to parse pocketbase upload response", "error", err.Error(), "body", resp.String())
		return "", utils.CapitalizeError("failed to parse response")
	}

	recordID, _ := result["id"].(string)
	fileName, _ := result[fileFieldName].(string)
	if recordID == "" || fileName == "" {
		utils.Error("pocketbase upload response missing file fields", "body", resp.String())
		return "", utils.CapitalizeError("failed to parse response")
	}

	fileURL := fmt.Sprintf("%s/api/files/%s/%s/%s", BaseURL(), collectionName, recordID, fileName)
	return fileURL, nil
}
