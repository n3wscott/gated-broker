// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import "strings"

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	RegisterRequest chan *Request
	ReleaseRequest  chan *Request
	Requests        map[*Request]bool
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:       make(chan []byte),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Clients:         make(map[*Client]bool),
		RegisterRequest: make(chan *Request),
		ReleaseRequest:  make(chan *Request),
		Requests:        make(map[*Request]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case request := <-h.RegisterRequest:
			h.Requests[request] = true
		case request := <-h.ReleaseRequest:
			if _, ok := h.Requests[request]; ok {
				delete(h.Requests, request)
				close(request.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}

			id := string(message)
			if strings.HasPrefix(id, "> ") {
				id := strings.TrimPrefix(id, "> ")
				for request := range h.Requests {
					if request.ID == id {
						request.Send <- &Response{ID: id, Approved: true}
					} else {
						request.Send <- &Response{ID: id, Approved: false}
					}
					close(request.Send)
					delete(h.Requests, request)
				}
			}
		}
	}
}
