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
	conns   = make(map[int]*websocket.Conn)
	counter = make(map[int]int)
	mu      = &sync.Mutex{}
)

func (H *Handler) ChatService(w http.ResponseWriter, r *http.Request) {
	user, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJson(w, http.StatusUnauthorized, "unothorized")
		return
	}
	receiver := r.URL.Query().Get("to")
	if user.Value == "" {
		http.Error(w, "User not specified", http.StatusBadRequest)
		return
	}

	user_name, _, user_id, receiver_id, err := H.Service.Database.GetId(user.Value, receiver)
	if err != nil {
		if err == sqlite3.ErrLocked {
			http.Error(w, "data base locked", http.StatusLocked)
			return
		}
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, "bad request")
			return
		}
		if err == sqlite3.ErrLocked {
			utils.WriteJson(w, http.StatusLocked, "data base locked")
			return
		}
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		fmt.Println(user_name + " disconnected")
		mu.Lock()
		counter[user_id]--
		if counter[user_id] == 0 {
			// update and handle user status: online or offline
			delete(conns, user_id)
		}
		mu.Unlock()
		conn.Close()
	}()

	mu.Lock()
	conns[user_id] = conn
	counter[user_id]++
	fmt.Println(user_name + " connected")
	mu.Unlock()

	go readloop(user_name, user_id, receiver_id, H.Service.Database.Db)
}

func readloop(sendername string, userid int, receiverid int, db *sql.DB) {
	conn := conns[userid]
	for {
		messageType, Message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		mu.Lock()
		data := Data_send{
			Sender:  sendername,
			Message: string(Message),
			Date:    time.Now(),
		}

		datajson, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			return
		}

		err = InsertChat(userid, receiverid, Message, db)
		if err != nil {
			fmt.Println(err)
		}

		mu.Unlock()
		// message

		// status

		// typing

		for k, value := range conns {
			if k == receiverid || k == userid {
				if err := value.WriteMessage(messageType, datajson); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func InsertChat(From, To int, Message []byte, Db *sql.DB) error {
	var Conversations_ID int64
	err := Db.QueryRow("SELECT id FROM conversations WHERE (user_one = ? AND user_two = ?) OR (user_one = ? AND user_two = ?)", From, To, To, From).Scan(&Conversations_ID)
	if err != nil {
		if err == sql.ErrNoRows {
			Insertchat, err := Db.Exec("INSERT INTO conversations (user_one, user_two, created_at) VALUES (?, ?)", From, To, time.Now())
			if err != nil {
				return err
			}

			Conversations_ID, err = Insertchat.LastInsertId()
			if err != nil {
				return err
			}
		}
	} else {
		_, err = Db.Exec("UPDATE INTO conversations (created_at) VALUES (?)", time.Now())
		if err != nil {
			return err
		}
	}
	_, err = Db.Exec("INSERT INTO messages (sender_id, content, conversation_id) VALUES (?, ?, ?)", From, string(Message), Conversations_ID)
	if err != nil {
		return err
	}
	return nil
}

func (H *Handler) GethistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJson(w, http.StatusMethodNotAllowed, "MethodNotAllowed")
		return
	}
	to := Message{
		Message: "string",
	}
	user, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJson(w, 500, "unothorized")
		return
	}
	err = json.NewDecoder(r.Body).Decode(&to)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
		return
	}
	_, _, user_id, to_id, err := H.Service.Database.GetId(user.Value, to.Message)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, "bad request")
		return
	}
	HistoryMessages, err := H.Service.Database.HistoryMessages(user_id, to_id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteJson(w, http.StatusBadRequest, "bad request")
			return
		}
		if err == sqlite3.ErrLocked {
			utils.WriteJson(w, http.StatusLocked, "bad request")
			return
		}
		utils.WriteJson(w, http.StatusInternalServerError, "internal server err")
		return
	}

	utils.WriteJson(w, http.StatusOK, HistoryMessages)
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
		fmt.Println("get last conversation", err)
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
