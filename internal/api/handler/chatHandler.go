/*package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"real-time-forum/internal/models"

	"github.com/gorilla/websocket"
	"github.com/mattn/go-sqlite3"
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
	HistoryMessages []models.Messagesbody
}

var (
	conns = make(map[string]*websocket.Conn)
	mu    = &sync.Mutex{}
	data  Data_send
)

func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	to := r.URL.Query().Get("to")
	if user.Value == "" {
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	user_id, to_id, err := H.Service.Database.GetId(user.Value, to)
	if err != nil {
		if err == sqlite3.ErrLocked {
			http.Error(w, "data base locked", http.StatusLocked)
		}
		// err bad request theres no sender or no receiver
		// err db is locked
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		mu.Lock()
		fmt.Println(user.Value + " disconnected")
		delete(conns, user.Value)
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	conns[user.Value] = conn

	var Online_users []string
	Online_users = append(Online_users, user.Value)
	mu.Unlock()

	go broadcast(conns, Online_users)

	fmt.Println(user.Value + " connected")
	HistoryMessages := H.Service.GetHistory(user_id, to_id)

	for {

		messageType, Message, err := conns[user.Value].ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		err = H.Service.Database.InsertChat(user_id, to_id, Message)
		if err != nil {
			// err message
		}

		for k, conn := range conns {
			if k == to || k == user.Value {
				fmt.Println(k)
				if err := conn.WriteMessage(messageType, []byte(user.Value+" : "+string(Message))); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func broadcast(conns map[string]*websocket.Conn, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	for _, conn := range conns {
		if err = conn.WriteMessage(1, jsonData); err != nil {
			log.Println(err)
			return
		}
	}
}*/

package handler