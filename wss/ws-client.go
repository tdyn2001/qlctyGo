package wss

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	// the websocket connection
	connection *websocket.Conn
	manager    *WsServer
	egress     chan []byte
	user       string
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, wss *WsServer, userName string) *WsClient {
	return &WsClient{
		connection: conn,
		manager:    wss,
		egress:     make(chan []byte),
		user:       userName,
	}
}

func (c *WsClient) readMessages() {
	defer func() {
		log.Printf("Client disconnected.")
		c.manager.removeClient(c)
	}()
	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		log.Println("user: ", c.user)
		log.Println("Payload: ", string(payload))

		for wsClient := range c.manager.clients {
			if wsClient != c {
				wsClient.egress <- payload
			}
		}
	}
}

func (c *WsClient) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			sendMessage := fmt.Sprintf("%s : %s", c.user, message)
			if err := c.connection.WriteMessage(websocket.TextMessage, []byte(sendMessage)); err != nil {
				log.Println(err)
			}
			log.Println("sent message")
		}

	}
}
