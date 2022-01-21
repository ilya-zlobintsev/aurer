package web

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ilyazzz/aurer/internal"
)

var upgrader = websocket.Upgrader{}

func (web *Web) statusSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	defer conn.Close()

	c := make(chan internal.StatusMsg)

	web.c.StatusChans = append(web.c.StatusChans, c)

	defer func(web *Web, c chan internal.StatusMsg) {
		chans := &web.c.StatusChans

		for i, other := range *chans {
			if other == c {
				(*chans)[i] = (*chans)[len(*chans)-1]

				*chans = (*chans)[:len(*chans)-1]
			}
		}
	}(web, c)

	msg := web.c.GetStatus()

	for {
		err := conn.WriteJSON(msg)

		if err != nil {
			log.Printf("Failed responding with json")
			break
		}

		msg = <-c
	}
}
