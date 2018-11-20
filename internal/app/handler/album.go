package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sparcs-home-go/internal/app/service"
)

// GetAlbumResponse : GetAlbum response body
type GetAlbumResponse struct {
	YearList  []service.YearSchema  `json:"yearList"`
	AlbumList []service.AlbumSchema `json:"albumList"`
}

// GetAlbum : response = {yearList, albumList}
func GetAlbum(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAlbum called")
	yearList, albumList, err := service.ListAlbum()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred on ListAlbum"))
		return
	}
	result := GetAlbumResponse{yearList, albumList}
	payload, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred on Marshal"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

// UploadPhoto : upload new photo in album
func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	var req service.UploadPhotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	if err := service.UploadPhoto(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to upload photo"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CreateAlbum : create new album
func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	if err := service.CreateAlbum(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to create album"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteAlbum : delete album
func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	albumID, err := strconv.Atoi(req["albumID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("AlbumID is not integer"))
		return
	}
	if err := service.DeleteAlbum(albumID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to delete album"))
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeletePhoto : delete photo from DB only, not from storage
func DeletePhoto(w http.ResponseWriter, r *http.Request) {
	req := mux.Vars(r)
	photoID, err := strconv.Atoi(req["photoID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("PhotoID is not integer"))
		return
	}
	if err := service.DeletePhoto(photoID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to delte photo"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
