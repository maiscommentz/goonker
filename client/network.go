package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"Goonker/common"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type NetworkClient struct {
	conn      *websocket.Conn
	ctx       context.Context
	ctxCancel context.CancelFunc
	sendMu    sync.Mutex

	// We buffer incoming packets so the Game Loop isn't blocked by network lag
	incomingMessages chan common.Packet
}

func NewNetworkClient() *NetworkClient {
	return &NetworkClient{
		incomingMessages: make(chan common.Packet, 100),
	}
}

// Connect dials the server (ws://localhost:8080/ws for local dev)
func (c *NetworkClient) Connect(url string, roomID string, isBot bool) error {
	// Thread-safe check if already connected
	c.sendMu.Lock()
	if c.conn != nil {
		c.sendMu.Unlock()
		log.Println("Already connected to server")
		return nil
	}
	c.sendMu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return err
	}

	//Assign the connection safely
	c.sendMu.Lock()
	c.conn = conn
	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.sendMu.Unlock()

	joinPayload := common.JoinPayload{
        RoomID: roomID,
        IsBot:  isBot,
    }
    
    data, _ := json.Marshal(joinPayload)
    packet := common.Packet{
        Type: common.MsgJoin,
        Data: data,
    }

	if err := wsjson.Write(ctx, c.conn, packet); err != nil {
        c.conn.Close(websocket.StatusInternalError, "failed to send join")
        return err
    }

	// Start listening immediately in a separate goroutine
	go c.listen()

	log.Println("Connected to server at", url)
	return nil
}

func (c *NetworkClient) listen() {
	defer func() {
		// Lock, Close, and set c.conn to nil so we can reconnect later
		c.sendMu.Lock()
		if c.conn != nil {
			c.conn.Close(websocket.StatusInternalError, "connection closed")
			c.conn = nil // Important: Reset so Connect() works again
		}
		c.sendMu.Unlock()
		
		if c.ctxCancel != nil {
			c.ctxCancel()
		}
	}()

	for {
		var packet common.Packet
		
		// Use the client context for reading
		err := wsjson.Read(c.ctx, c.conn, &packet)
		if err != nil {
			// If context was canceled or connection closed, just return
			if c.ctx.Err() != nil {
				return
			}
			log.Printf("Read error: %v", err)
			return
		}

		// Non-blocking send to channel
		select {
		case c.incomingMessages <- packet:
		default:
			log.Println("Network buffer full, dropping packet")
		}
	}
}

// SendPacket sends a packet to the server
func (c *NetworkClient) SendPacket(packet common.Packet) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if c.conn == nil {
		log.Println("Not connected to server")
		return nil 
	}

	ctx, cancel := context.WithTimeout(c.ctx, time.Second*5)
	defer cancel()

	return wsjson.Write(ctx, c.conn, packet)
}

// Poll gets the next packet from the queue (Non-blocking)
func (c *NetworkClient) Poll() *common.Packet {
	select {
	case msg := <-c.incomingMessages:
		return &msg
	default:
		return nil
	}
}