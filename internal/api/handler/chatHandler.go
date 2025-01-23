package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	utils "real-time-forum/pkg"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Message string
}

type Data_send struct {
	Sender string
	Message     string
}

var (
	conns = make(map[string]*websocket.Conn)
	mu    = &sync.Mutex{}
	data  Data_send
)

func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJson(w, 500, "no cookies")
		return
	}
	to := r.URL.Query().Get("to")
	if user.Value == "" {
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	user_name,to, user_id, to_id, err := H.Service.Database.GetId(user.Value, to)
	if err != nil {
		// err message
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		fmt.Println(user_name + " disconnected")
		mu.Lock()
		delete(conns, user.Value)
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	conns[user.Value] = conn
	/*data.List_online = append(data.List_online, user.Value)
	datajson, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}

	if err = conn.WriteMessage(1, datajson); err != nil {
		log.Println(err)
		return
	}*/

	fmt.Println(user_name + " connected")
	mu.Unlock()

	for {
		messageType, Message, err := conns[user.Value].ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		data := Data_send {
			Sender: user_name,
			Message: string(Message),
		}
		datajson, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
		err = H.Service.Database.InsertChat(user_id, to_id, Message)
		if err != nil {
			fmt.Println(err)
		}
		mu.Unlock()
		for k, value := range conns {
			if k == to || k == user.Value {
				fmt.Println(k)
				if err := value.WriteMessage(messageType, datajson); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func (H *Handler) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}
	to := Message{
		Message: "string",
	}
	user, err := r.Cookie("session_token")
	if err != nil {
		
		utils.WriteJson(w, 500, "no cookies")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&to)
	if err != nil {
		utils.WriteJson(w, 500, "err to")
		return
	}
	_,_,user_id, to_id, err := H.Service.Database.GetId(user.Value, to.Message)
	if err != nil {
		utils.WriteJson(w, 500, "err looking for ids")
		return
	}
	HistoryMessages, err := H.Service.Database.HistoryMessages(user_id, to_id)
	if err != nil {
		utils.WriteJson(w, 500, "err history")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(HistoryMessages)
	if err != nil {
		utils.WriteJson(w, 500, "err send")
		return
	}
}
