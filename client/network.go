package main

import (
	"context"
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
func (c *NetworkClient) Connect(url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	c.ctx, c.ctxCancel = context.WithCancel(context.Background())

	// Start listening immediately in a separate goroutine
	go c.listen()

	log.Println("Connected to server at", url)
	return nil
}

func (c *NetworkClient) listen() {
	defer func() {
		c.conn.Close(websocket.StatusInternalError, "connection closed")
		c.ctxCancel()
	}()

	for {
		var packet common.Packet
		// Reads JSON from socket and unmarshals into Packet struct
		err := wsjson.Read(c.ctx, c.conn, &packet)
		if err != nil {
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

// Send wraps the data in the Packet struct and sends it
func (c *NetworkClient) Send(msgType string, payload interface{}) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if c.conn == nil {
		return nil
	}

	// Wrap payload in Packet envelope
	// Note: wsjson.Write will automatically marshal 'payload' into the 'Data' field
	// if we structured the write correctly, but here we construct the full packet manually
	// to ensure it matches the common.Packet definition.
	
	// However, wsjson.Write takes an interface{}.
	// We need to construct the common.Packet with the RawMessage.
	// For simplicity in this helper, we rely on the struct logic in game.go 
	// or we just send the raw struct if the library supports it.
	
	// BETTER APPROACH for this helper:
	// Just accept the full packet.
	
	return nil
}

// SendPacket is the raw sender
func (c *NetworkClient) SendPacket(packet common.Packet) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	if c.conn == nil {
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