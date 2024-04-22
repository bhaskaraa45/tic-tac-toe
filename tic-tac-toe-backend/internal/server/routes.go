package server

import (
	"log"
	"net/http"
	"strings"
	"tic-tac-toe-backend/internal/socket"

	"github.com/gin-gonic/gin"
)

func removeTrailingSlashMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/" && strings.HasSuffix(c.Request.URL.Path, "/") {
			c.Request.URL.Path = strings.TrimSuffix(c.Request.URL.Path, "/")
		}
		c.Next()
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(removeTrailingSlashMiddleware())

	r.Use(func(c *gin.Context) {
		log.Println("CORS middleware triggered")

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})

	r.GET("/", s.HelloWorldHandler)
	r.GET("/ws", socket.HandleConnections)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}
