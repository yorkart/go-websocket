package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WSHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	vars := mux.Vars(r)
	roomId := vars["roomId"]

	if err := h.enterRoom(roomId, c); err != nil {
		log.Println("register failure", roomId)
	}
}
