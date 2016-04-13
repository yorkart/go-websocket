package ws

import (
	"encoding/json"
	"log"
)

type WSServer interface {
	GetRoomId() string
	OnConnection(*Peer)
	OnReceive(*Peer, []byte)
	OnDisConnection(*Peer)
}

type TextPackage struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

type Action func(*Peer, string)

type WSEventTextServer struct {
	roomId string
	events map[string]Action
}

func NewWSEventTextServer(roomId string) *WSEventTextServer {
	return &WSEventTextServer{
		roomId: roomId,
		events: map[string]Action{},
	}
}

func (p *WSEventTextServer) On(event string, action Action) *WSEventTextServer {
	p.events[event] = action
	return p
}

func (p *WSEventTextServer) Send(peer *Peer, event string, message string) error {
	textPackage := &TextPackage{Event: event, Message: message}
	data, err := json.Marshal(textPackage)
	if err != nil {
		return err
	}
	peer.Send(data)
	return nil
}

func (p *WSEventTextServer) Broadcast(peer *Peer, event string, message string) error {
	textPackage := &TextPackage{Event: event, Message: message}
	data, err := json.Marshal(textPackage)
	if err != nil {
		return err
	}
	peer.Broadcast(data)
	return nil
}

func (p *WSEventTextServer) GetRoomId() string {
	return p.roomId
}

func (p *WSEventTextServer) OnConnection(peer *Peer) {
	//    fmt.Println("OnConnection")
	if action, ok := p.events["connection"]; ok {
		action(peer, "")
	}
}

func (p *WSEventTextServer) OnReceive(peer *Peer, data []byte) {
	log.Println("OnReceive <= ", string(data))
	//    peer.Send(data)

	textPackage := &TextPackage{}
	err := json.Unmarshal(data, textPackage)
	if err != nil {
		if action, ok := p.events["error"]; ok {
			action(peer, err.Error())
		}
		return
	}
	log.Println(textPackage)
	if action, ok := p.events[textPackage.Event]; ok {
		action(peer, textPackage.Message)
	}
	//    if action, ok := p.events["connection"]; ok {
	//        action(peer, nil)
	//    }
}

func (p *WSEventTextServer) OnDisConnection(peer *Peer) {
	//    fmt.Println("OnDisConnection")
	if action, ok := p.events["disconnection"]; ok {
		action(peer, "")
	}
}
