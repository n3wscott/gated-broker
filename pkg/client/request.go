// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"crypto/rand"
	"fmt"
	"io"
)

type Request struct {
	ID   string
	hub  *Hub
	Send chan *Response
}

type Response struct {
	ID       string
	Approved bool
}

func NewRequest(hub *Hub) *Request {

	id, _ := newUUID()

	r := Request{
		ID:   id,
		hub:  hub,
		Send: make(chan *Response, 2),
	}
	r.hub.RegisterRequest <- &r

	r.hub.Broadcast <- []byte(id) // TODO this should be an object at some point. like {id: {body}}

	return &r
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
