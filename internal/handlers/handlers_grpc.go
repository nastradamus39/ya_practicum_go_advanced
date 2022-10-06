package handlers

import (
	"fmt"
	"net/http"
)

// CreateShortURLGRPCHandler — создает короткий урл.
func CreateShortURLGRPCHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateShortURLGRPCHandler")
}

// GetShortURLGRPCHandler — возвращает полный урл по короткому.
func GetShortURLGRPCHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetShortURLGRPCHandler")
}
