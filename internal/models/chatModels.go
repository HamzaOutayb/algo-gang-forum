package models

type Data_send struct {
	Message         string
	HistoryMessages []Messagesbody
	List_online     []string
}

type Messagesbody struct {
	Sender  string
	Message string
	Date    string
}
