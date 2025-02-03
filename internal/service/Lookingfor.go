package service

import "github.com/gorilla/websocket"

func (S *Service) LookingForIndexconns(slice []*websocket.Conn, conn *websocket.Conn) int {
	for i := 0; i < len(slice); i++ {
		if conn == slice[i] {
			return i
		}
	}
	return 0
}
