package wss

import (
	"log"
	"sync"
	gRedis "v2/external-service/redis"
	"v2/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	/**
	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type ClientList map[*WsClient]bool

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type WsServer struct {
	sync.RWMutex
	clients     ClientList
	Broadcast   chan []byte
	redisClient *gRedis.RedisClient
}

// NewManager is used to initalize all the values inside the manager
func NewWebsocketServer(redisClient *gRedis.RedisClient) *WsServer {
	wsServer := &WsServer{
		clients:   make(ClientList),
		Broadcast: make(chan []byte),
	}
	wsServer.redisClient = redisClient
	// initKafka(wsServer)
	go wsServer.broadcasting()
	return wsServer
}

func (m *WsServer) addClient(client *WsClient) {
	// Lock so we can manipulate
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

func (m *WsServer) removeClient(client *WsClient) {
	m.Lock()
	defer m.Unlock()

	// Check if Client exists, then delete it
	if _, ok := m.clients[client]; ok {
		// close connection
		client.connection.Close()
		// remove
		delete(m.clients, client)
	}
}

func (m *WsServer) broadcasting() {
	for message := range m.Broadcast {
		for wsClient := range m.clients {
			wsClient.egress <- message
		}
	}
}

// serveWS is a HTTP Handler that the has the Manager that allows connections
func (ws *WsServer) SetupWSS(ctx *gin.Context) {
	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user := ctx.MustGet("currentUser").(models.User)
	wsClient := NewClient(conn, ws, user.Name)
	ws.addClient(wsClient)
	go wsClient.readMessages()
	go wsClient.writeMessages()
}

// func initKafka(ws *WsServer) {
// 	go kafkas.Consume(context.Background(), "notification", "general", func(msg []byte) {
// 		ws.broadcast <- msg
// 	})
// }
