package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 0
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type Peer struct {
	roomId string
	room   *Room
	props  map[string]interface{}
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func NewPeer(roomId string, ws *websocket.Conn) *Peer {
	return &Peer{
		roomId: roomId,
		ws:     ws,
		send:   make(chan []byte, 256),
	}
}

func (p *Peer) SetProp(key string, value interface{}) {
	p.props[key] = value
}

func (p *Peer) GetProp(key string) (interface{}, bool) {
	v, ok := p.props[key]
	return v, ok
}

func (p *Peer) listen() {
	go p.writePump()
	p.readPump()
}

// readPump pumps messages from the websocket connection to the hub.
func (peer *Peer) readPump() {
	defer func() {
		h.unregister <- peer
		peer.ws.Close()
		peer.room.server.OnDisConnection(peer)
	}()

	peer.ws.SetReadLimit(maxMessageSize)
	peer.ws.SetReadDeadline(time.Now().Add(pongWait))
	peer.ws.SetPongHandler(func(string) error { peer.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	//    fmt.Println("r server ", peer)
	//    time.Sleep(time.Second)
	peer.room.server.OnConnection(peer)
	for {
		_, message, err := peer.ws.ReadMessage()
		if err != nil {
			log.Println("read error")
			break
		}

		peer.room.server.OnReceive(peer, message)
	}
}

// write writes a message with the given message type and payload.
func (c *Peer) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (p *Peer) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.ws.Close()
	}()
	for {
		select {
		case message, ok := <-p.send:
			if !ok {
				p.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := p.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := p.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (p *Peer) Send(data []byte) {
	p.send <- data
}

func (p *Peer) Broadcast(data []byte) {
	p.room.Broadcast(data)
}
