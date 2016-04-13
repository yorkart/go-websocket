package ws

import (
	"errors"
	"log"
	"sync"

	"wsbroadcaster/util"

	"github.com/gorilla/websocket"
)

var h *hub

func init() {
	h = newHub()
	go h.run()
	log.Println("new hub")
}

type msg struct {
	roomId string
	data   []byte
}

type hub struct {
	// Registered connections.
	rooms *util.Map

	// Inbound messages from the connections.
	broadcast chan *msg

	// Register requests from the connections.
	register chan *Peer

	// Unregister requests from connections.
	unregister chan *Peer

	lock *sync.Mutex
}

func newHub() *hub {
	return &hub{
		rooms:      util.NewMap(),
		broadcast:  make(chan *msg),
		register:   make(chan *Peer),
		unregister: make(chan *Peer),
		lock:       &sync.Mutex{},
	}
}

func (p *hub) enterRoom(roomId string, ws *websocket.Conn) error {
	h.rooms.Range(func(key util.Key, val util.Val) {
		fmt.Println(key, val)
	})

	if r, ok := h.rooms.Get(roomId); ok {
		p := NewPeer(roomId, ws)
		p.room = r.(*Room)

		fmt.Println("r +>")
		h.register <- p
		fmt.Println("r ->")

		p.listen()
		return errors.New("session finish")
	}

	return errors.New("no room enter")
}

func (p *hub) run() {
	for {
		select {
		case peer := <-h.register:
			if r, ok := p.rooms.Get(peer.roomId); ok {
				r.(*Room).add(peer)
			}

			fmt.Println("add peer to room ", peer.roomId)
		case peer := <-h.unregister:
			if r, ok := p.rooms.Get(peer.roomId); ok {
				r.(*Room).del(peer)
				close(peer.send)
			}
		case m := <-h.broadcast:
			if r, ok := p.rooms.Get(m.roomId); ok {
				for _, item := range r.(*Room).connections.List() {
					var peer = item.(*Peer)
					select {
					case item.(*Peer).send <- m.data:
					default:
						close(peer.send)
						r.(*Room).connections.Remove(peer)
					}
				}
			}
		}
	}
}

func Listen(server WSServer) *Room {
	roomId := server.GetRoomId()
	r := NewRoom(server)
	fmt.Println(h, r)
	h.rooms.Put(roomId, r)
	return r
}

func Register(roomId string, ws *websocket.Conn) (*Peer, bool) {
	h.rooms.Range(func(key util.Key, val util.Val) {
		fmt.Println(key, val)
	})

	if r, ok := h.rooms.Get(roomId); !ok {
		p := NewPeer(roomId, ws)
		p.room = r.(*Room)

		fmt.Println("r +>")
		h.register <- p
		fmt.Println("r ->")

		return p, true
	}

	return nil, false
}
