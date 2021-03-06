package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware woof
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200") // für DEV: "http://localhost:4200" (erlaubt zugriffe von...)
		//c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN")) // für DEV: "http://localhost:4200" (erlaubt zugriffe von...)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
