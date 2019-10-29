// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	receive chan InboundMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	functionHandlers map[string]func(InboundMessage)
}

type OutboundMessage struct {
	TargetWorker string `json:"targetWorker"`
	MessageType string `json:"messageType"`
	Message map[string]interface{} `json:"message"`
}

type InboundMessage struct {
	WorkerId string `json:"workerId"`
	MessageType string `json:"messageType"`
	Message map[string]interface{} `json:"message"`
}

func NewHub() *Hub {
	return &Hub{
		receive:    	  make(chan InboundMessage),
		register:   	  make(chan *Client),
		unregister:       make(chan *Client),
		clients:    	  make(map[*Client]bool),
		functionHandlers: make(map[string]func(InboundMessage)),
	}
}

func (h *Hub) On(messageType string, handler func(InboundMessage)) {
	h.functionHandlers[messageType] = handler
}

func (h *Hub) Broadcast(message OutboundMessage) {
	logrus.Tracef("Broadcasting message %s", message)

	for client := range h.clients {
		if client.workerId == message.TargetWorker || message.TargetWorker == "" {
			bytes, _ := json.Marshal(message)

			select {
			case client.send <- bytes:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.receive:
			logrus.Tracef("Hub message received: %s", message)
			h.functionHandlers[message.MessageType](message)
		}
	}
}