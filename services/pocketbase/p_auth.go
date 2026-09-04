package pocketbase

import (
	"fmt"
	"os"
	"pos-master/models"
	"pos-master/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func BaseURL() string {
	baseURL := strings.TrimRight(os.Getenv("POCKETBASE_URL"), "/")
	if baseURL == "" {
		baseURL = "https://file-server.mygeepay.com"
	}
	return baseURL
}

func HandlePocketBaseAuth(c *gin.Context) (string, error) {
	identity := os.Getenv("POCKETBASE_IDENTITY")
	password := os.Getenv("POCKETBASE_PASSWORD")
	if identity == "" || password == "" {
		utils.Error("pocketbase credentials are not configured",
			"has_identity", identity != "",
			"has_password", password != "",
		)
		return "", utils.CapitalizeError("file storage is not configured")
	}

	body := map[string]string{
		"identity": identity,
		"password": password,
	}

	client := resty.New()

	var authResp models.PocketBaseAuthResponse
	endpoint := BaseURL() + "/api/collections/_superusers/auth-with-password"

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&authResp).
		Post(endpoint)

	if err != nil {
		utils.Error("pocketbase auth request failed", "error", err.Error(), "endpoint", endpoint)
		return "", utils.CapitalizeError(fmt.Sprintf("Unable to login: %v", err))
	}

	if resp.IsError() {
		utils.Error("pocketbase auth rejected",
			"status", resp.StatusCode(),
			"body", resp.String(),
		)
		return "", utils.CapitalizeError("Unable to login")
	}

	return authResp.Token, nil
}
