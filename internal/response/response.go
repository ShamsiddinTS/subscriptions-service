package response

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, ErrorResponse{
		Error: msg,
	})
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
