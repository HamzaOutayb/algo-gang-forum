package models

type Data_send struct {
	Message         string
	HistoryMessages []Messagesbody
	List_online     []string
}

type Messagesbody struct {
	Sender  string
	To      string
	Message string
	Date    string
}

type Chat struct {
	Id        string `json:"id"`
	Nickname  string `json:"nickname"`
	FriendsId string `json:"friendid"`
}

type Conversations struct {
	Sender     string
	Content    string
	Created_at string
}
