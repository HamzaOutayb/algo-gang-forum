package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"real-time-forum/internal/models"
	utils "real-time-forum/pkg"

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

type Message struct {
	Message string
}

type Data_send struct {
	Sender  string
	Message string
	Date    time.Time
	To      string
}

var (
	conns = make(map[int]*websocket.Conn)
	mu    = &sync.Mutex{}
	data  Data_send
)

func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJson(w, 500, "no cookies")
		return
	}
	To := r.URL.Query().Get("to")
	if user.Value == "" {
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	user_name, _, user_id, to_id, err := H.Service.Database.GetId(user.Value, To)
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
		fmt.Println(user_name + " disconnected")
		mu.Lock()
		delete(conns, user_id)
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	conns[user_id] = conn
	/*data.List_online = append(data.List_online, user_id)
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

	// go broadcast(conns, Online_users)

	fmt.Println(user.Value + " connected")
	// HistoryMessages := H.Service.GetHistory(user_id, to_id)

	for {
		messageType, Message, err := conns[user_id].ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		data := Data_send{
			Sender:  user_name,
			Message: string(Message),
			Date:    time.Now(),
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
			if k == to_id || k == user_id {
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
	_, _, user_id, to_id, err := H.Service.Database.GetId(user.Value, to.Message)
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
}

func (H *Handler) Lastconversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJson(w, http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "bad request"})
		return
	}

	cookie, err := r.Cookie("Session_token")
	usrid := 0
	if err != http.ErrNoCookie && H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		usrid, _ = H.Service.Database.GetUser(cookie.Value)
	} else {
		utils.WriteJson(w, http.StatusUnauthorized, "unauthorized user")
		return
	}
	pagenm := r.URL.Query().Get("Page-num")
	page, err := strconv.Atoi(pagenm)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
		return
	}

	chat, err := H.Service.GetLastconversations(page, usrid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.WriteJson(w, http.StatusOK, []models.Chat{})
			return
		case sqlite3.ErrLocked:
			utils.WriteJson(w, http.StatusLocked, struct {
				Error string `json:"error"`
			}{Error: "Database Locked"})
			return
		}
	}
	utils.WriteJson(w, http.StatusOK, chat)
}

func (H *Handler) Conversations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJson(w, http.StatusBadRequest, struct {
			Error string `json:"error"`
		}{Error: "bad request"})
		return
	}

	cookie, err := r.Cookie("Session_token")
	usrid := 0
	if err != http.ErrNoCookie && H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		usrid, _ = H.Service.Database.GetUser(cookie.Value)
	} else {
		utils.WriteJson(w, http.StatusUnauthorized, "unauthorized user")
		return
	}
	pagenm := r.URL.Query().Get("Page-num")
	page, err := strconv.Atoi(pagenm)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
		return
	}

	chat, err := H.Service.Getconversations(page, usrid)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.WriteJson(w, http.StatusOK, []models.Chat{})
			return
		case sqlite3.ErrLocked:
			utils.WriteJson(w, http.StatusLocked, struct {
				Error string `json:"error"`
			}{Error: "Database Locked"})
			return
		}
	}
	utils.WriteJson(w, http.StatusOK, chat)
}
