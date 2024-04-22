package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // This should be more restrictive in production.
		},
	}
	clients      = make(map[string]map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

func HandleConnections(c *gin.Context) {
	roomID := c.Query("room")
	join := c.Query("join")

	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	if join == "1" {
		if _, ok := clients[roomID]; !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room does not exist"})
			return
		}

		if GetClientCount(roomID) != 1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room does not exist or maybe full"})
			return
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	registerClient(roomID, conn)
	defer removeClient(roomID, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		broadcast(roomID, conn, msg)
	}
}

func registerClient(roomID string, conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	if _, ok := clients[roomID]; !ok {
		clients[roomID] = make(map[*websocket.Conn]bool)
	}

	clients[roomID][conn] = true
}

func removeClient(roomID string, conn *websocket.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	if _, ok := clients[roomID]; ok {
		delete(clients[roomID], conn)
	}
}

func broadcast(roomID string, sender *websocket.Conn, message []byte) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for conn := range clients[roomID] {
		if conn != sender {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				conn.Close()
				delete(clients[roomID], conn)
			}
		}
	}
}

// GetClientCount returns the number of clients in the specified room.
func GetClientCount(roomID string) int {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	if clientSet, exists := clients[roomID]; exists {
		return len(clientSet)
	}
	return 0
}
