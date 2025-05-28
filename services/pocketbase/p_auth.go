package pocketbase

import (
	"net/http"
	"pos-master/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func HandleAuth(c *gin.Context) {
	var creds models.PocketBaseCredentials

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid input",
			"detail": err,
		})
		return
	}

	client := resty.New()

	var authResp models.PocketBaseAuthResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(creds).
		SetResult(&authResp).
		Post("http://89.250.72.76:8090/api/collections/_superusers/auth-with-password")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "request failed",
			"detail": err.Error(),
		})
		return
	}

	if resp.IsError() {
		c.JSON(resp.StatusCode(), gin.H{
			"error":  "authentication failed",
			"detail": resp.String(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": authResp.Token,
	})

}
