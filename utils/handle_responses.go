package utils

import "github.com/gin-gonic/gin"

const (
	ResponseErrorKey       = "response_error"
	ResponseErrorDetailKey = "response_error_detail"
)

func RespondWithError(c *gin.Context, statusCode int, errorMessage string, details ...string) {
	c.Set(ResponseErrorKey, errorMessage)
	if len(details) > 0 {
		c.Set(ResponseErrorDetailKey, details[0])
	}

	response := gin.H{
		"status": "failure",
		"error":  errorMessage,
	}

	if len(details) > 0 {
		response["detail"] = details[0]
	}

	c.JSON(statusCode, response)
}

func RespondWithSuccess(c *gin.Context, message string, data ...gin.H) {
	response := gin.H{
		"status":  "success",
		"message": message,
	}

	if len(data) > 0 {
		for key, value := range data[0] {
			response[key] = value
		}
	}

	c.JSON(200, response)
}
