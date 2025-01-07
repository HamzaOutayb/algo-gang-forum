package utils

import (
	"encoding/json"
	"net/http"
)


func WriteJson(w http.ResponseWriter, statuscode int, Data any) error {
	w.WriteHeader(statuscode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Data)
	if err != nil {
		return err
	}
	return nil
}
/*
func ReadJson(r *http.Request) any {

}*/