// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import "go-websocket/util"

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Room struct {
	roomId string
	// Registered connections.
	connections *util.Set

	server WSServer
}

func NewRoom(server WSServer) *Room {
	return &Room{
		roomId:      server.GetRoomId(),
		connections: util.NewSet(),
		server:      server,
	}
}

func (p *Room) Broadcast(data []byte) {
	h.broadcast <- &msg{roomId: p.roomId, data: data}
}

func (p *Room) add(peer *Peer) {
	peer.room = p
	p.connections.Add(peer)
}

func (p *Room) del(peer *Peer) {
	p.connections.Remove(peer)
}
