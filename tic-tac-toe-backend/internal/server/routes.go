package server

import (
	"fmt"
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
			c.Redirect(http.StatusMovedPermanently, c.Request.URL.Path)
		}
	}
}

func (s *Server) RegisterRoutes() {
	fmt.Println("Configuring routes")
	s.Use(removeTrailingSlashMiddleware())

	// CORS middleware
	s.Use(func(c *gin.Context) {
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

	s.GET("/", s.HelloWorldHandler)
	s.GET("/ws", socket.HandleConnections)
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := map[string]string{"message": "Hello World"}
	c.JSON(http.StatusOK, resp)
}
