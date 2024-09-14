package response

import "github.com/gin-gonic/gin"

func SendSuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, map[string]interface{}{
		"message": "ok",
		"data":    data,
	})
}

func SendErrorResponse(c *gin.Context, message string, statusCode, errorCode int) {
	c.JSON(statusCode, map[string]interface{}{
		"error_code": errorCode,
		"message":    message,
	})
}
