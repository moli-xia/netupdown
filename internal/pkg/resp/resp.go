package resp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moli-xia/netupdown/internal/pkg/apperr"
)

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": data})
}
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "ok", "data": data})
}
func Fail(c *gin.Context, err error) {
	var ae *apperr.Error
	if errors.As(err, &ae) {
		c.JSON(ae.Status, gin.H{"code": ae.Code, "message": ae.Message, "data": nil})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": 10000, "message": "服务器内部错误", "data": nil})
}
