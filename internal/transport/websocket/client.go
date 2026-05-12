package websocket

import (
	"time"

	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/so1icitx/StockWise/internal/application"
)

type client struct {
	hub  *Hub
	conn *gorillawebsocket.Conn
	send chan application.NotificationEvent
}

func (client *client) readPump() {
	defer func() {
		client.hub.unregister <- client
		_ = client.conn.Close()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		return client.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		if _, _, err := client.conn.NextReader(); err != nil {
			break
		}
	}
}

func (client *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()

	for {
		select {
		case event, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = client.conn.WriteMessage(gorillawebsocket.CloseMessage, []byte{})
				return
			}
			if err := client.conn.WriteJSON(event); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(gorillawebsocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
