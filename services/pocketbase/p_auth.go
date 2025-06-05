package pocketbase

import (
	"fmt"
	"pos-master/models"
	"pos-master/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func HandlePocketBaseAuth(c *gin.Context) (string, error) {

	var body = map[string]string{
		"identity": "luswepo@geepay.co.zm",
		"password": "0147green",
	}

	client := resty.New()

	var authResp models.PocketBaseAuthResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&authResp).
		Post("http://102.23.120.239:8090/api/collections/_superusers/auth-with-password")

	if err != nil {
		return "", utils.CapitalizeError(fmt.Sprintf("Unable to login: %s", err.Error()))
	}

	if resp.IsError() {
		return "", utils.CapitalizeError("Unable to login")
	}

	return authResp.Token, nil

}
