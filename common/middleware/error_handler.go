package middleware

import (
	"net/http"

	resterrors "github.com/imylam/delivery-test/common/rest_errors"

	"github.com/gin-gonic/gin"
)

func HandleRestError(c *gin.Context) {
	c.Next()

	if c.Errors == nil {
		return
	}

	err := c.Errors.Last().Err

	if restErr, ok := err.(resterrors.RestError); ok {
		c.Header("HTTP", restErr.HttpStatusCodeString())
		c.JSON(restErr.HttpStatusCode(), gin.H{"error": restErr.Error()})
		return
	} else {
		c.Header("HTTP", "500")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
}
