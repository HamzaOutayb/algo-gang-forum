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
	User_name string
}

type Data_send struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
	Date    time.Time
	To      int `json:"to"`
	list    []int
}

var (
	conns = make(map[*websocket.Conn]int)
	mu    = &sync.Mutex{}
	data  Data_send
)


func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("session_token")
	if err != nil {
		
		utils.WriteJson(w, 500, "no cookies")
		return
	}
	
	if user.Value == "" {
		
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		
		log.Println(err)
		return
	}
	
	user_name,user_id, err := H.Service.Database.GetId(user.Value)
	if err != nil {
		
		if err == sqlite3.ErrLocked {
			http.Error(w, "data base locked", http.StatusLocked)
		}
		// err bad request theres no sender or no receiver
		// err db is locked
	}



	defer func() {
		fmt.Println(user_name + " disconnected")
		mu.Lock()
		delete(conns, conn)
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	conns[conn] = user_id
	if !checkisincluded(data.list, user_id) {
	data.list = append(data.list, user_id)
	}
	go broadcast(conns, data.list)
	fmt.Println(user_name + " connected")
	mu.Unlock()
	

	 

	fmt.Println(user.Value + " connected")

	for {
		messageType, Message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		UnmarshalData := Data_send{}
		json.Unmarshal(Message, &UnmarshalData)
		err = H.Service.Database.InsertChat(user_id, UnmarshalData.To, UnmarshalData.Message)
		if err != nil {
			fmt.Println(err)
		}
		mu.Unlock()
		for k, value := range conns {
			if value == UnmarshalData.To || value == user_id {
				if err := k.WriteMessage(messageType, Message); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func checkisincluded(list []int, id int) bool {
	for _, value := range list {
		if value == id {
			return true
		}
	}
	return false
}
func (H *Handler) GetHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}
	type Message struct{
		User_name string `json:"message"`
	}
	var to Message
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
	fmt.Println(to.User_name,"ttttttttttttttttttttttttttttttttttttttttttttttttttt")
	user_id, to_id, err := H.Service.Database.GetId2(user.Value, to.User_name)
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

func broadcast(conns map[*websocket.Conn]int, data any) {
	
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("broadcasting",	data, jsonData)
	for value := range conns {
		if err = value.WriteMessage(1, jsonData); err != nil {
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

	cookie, err := r.Cookie("session_token")
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
		page = 0
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

	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, "unauthorized user")
		return
	}
	usrid := 0
	if err != http.ErrNoCookie && H.Service.Database.CheckExpiredCookie(cookie.Value, time.Now()) {
		usrid, _ = H.Service.Database.GetUser(cookie.Value)
	}
	pagenm := r.URL.Query().Get("Page-num")
	page, err := strconv.Atoi(pagenm)
	if err != nil {
		page = 0
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
