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
	Status    map[int]bool
}

var (
	conns = make(map[int][]*websocket.Conn)
	mu    = &sync.Mutex{}
	statusmap = make(map[int]bool)
	data  Data_send
)


func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	fmt.Println("chat")
	user, err := r.Cookie("session_token")
	if err != nil {
		
		utils.WriteJson(w, 500, "no cookies")
		return
	}
	fmt.Println(user.Value)
	if user.Value == "" {
		
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}
	fmt.Println("chat2")
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
fmt.Println(user_name,user_id)


	defer func() {
		mu.Lock()
		
		statusmap[user_id] = false
		logout := Data_send{
			Status:    statusmap,
		}
		Indexconss := H.Service.LookingForIndexconns(conns[user_id],conn)
		conns[user_id] = append(conns[user_id][:Indexconss], conns[user_id][Indexconss+1:]...)
		go broadcast(conns, logout)
		conn.Close()
		mu.Unlock()
		fmt.Println(user_name + " disconnected")
	}()
		
//	mu.Lock()
	conns[user_id] = append(conns[user_id], conn)
	statusmap[user_id] = true
		login := Data_send{
			Status:    statusmap,
		}
	go broadcast(conns, login )
	
//	mu.Unlock()
	fmt.Println("chat3")

	for {
		var UnmarshalData Data_send
	  err := conn.ReadJSON( &UnmarshalData)
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		fmt.Println("UnmarshalData", UnmarshalData)
		err = H.Service.Database.InsertChat(user_id, UnmarshalData.To, UnmarshalData.Message)
		if err != nil {
			fmt.Println(err)
		}
		mu.Unlock()
		for _, value := range conns[user_id] {
				if err := value.WriteJSON(UnmarshalData); err != nil {
					log.Println(err)
					return
				}
		}
		for _, value := range conns[user_id] {
				if err := value.WriteJSON(UnmarshalData); err != nil {
					log.Println(err)
					return
				}
		}
	}
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
fmt.Println("to", to.User_name )
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
	
	utils.WriteJson(w, http.StatusOK, HistoryMessages)
}

func broadcast(conns map[int][]*websocket.Conn, data any) {
	fmt.Println("broadcast")
	fmt.Println(data)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return
	}
	for key,_ := range conns {
		for _,value := range conns[key]{
			if err = value.WriteMessage(1, jsonData); err != nil {
			log.Println(err)
			return
		}
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
