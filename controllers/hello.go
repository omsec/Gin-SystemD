package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SayHello(c *gin.Context) {
	time.Sleep(10 * time.Second)
	c.JSON(http.StatusOK, "hello there")
}
