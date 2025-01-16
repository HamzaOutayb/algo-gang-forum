package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Data_send struct {
	Message         string
	HistoryMessages []string
	List_online     []string
	List_offline    []string
}

var (
	conns = make(map[string]*websocket.Conn)
	mu    = &sync.Mutex{}
	data  Data_send
)

func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	to := r.URL.Query().Get("to")
	if user.Value == "" {
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	user_id, to_id, err := H.Service.Database.GetId(user.Value, to)
	if err != nil {
		// err message
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		fmt.Println(user.Value + " disconnected")
		mu.Lock()
		delete(conns, user.Value)
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	conns[user.Value] = conn
	data.List_online = append(data.List_online, user.Value)
	data.HistoryMessages = append(data.HistoryMessages, H.Service.GetHistory(user_id,to_id)...)
	datajson, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	if err = conn.WriteMessage(1, datajson); err != nil {
		log.Println(err)
		return
	}


	fmt.Println(user.Value + " connected")
	mu.Unlock()

	for {
		messageType, Message, err := conns[user.Value].ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		err = H.Service.Database.InsertChat(user_id, to_id, Message)
		if err != nil {
			// err message
		}
		mu.Unlock()
		for k, value := range conns {
			if k == to || k == user.Value {
				fmt.Println(k)
				if err := value.WriteMessage(messageType, []byte(user.Value+" : "+string(Message))); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
