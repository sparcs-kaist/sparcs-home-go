package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sparcs-home-go/internal/app/service"
)

// GetSeminar : GET seminar
func GetSeminar(w http.ResponseWriter, r *http.Request) {
	seminarList, err := service.GetSeminar()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to get seminar"))
		return
	}
	payload, err := json.Marshal(seminarList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred on Marshal"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// UploadSeminar : POST upload seminar
func UploadSeminar(w http.ResponseWriter, r *http.Request) {
	var req service.UploadSeminarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	if err := service.UploadSeminar(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to upload seminar"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
