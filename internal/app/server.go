package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sparcs-home-go/internal/app/configure"
	"github.com/sparcs-home-go/internal/app/handler"
)

// ServerInfo : host settings
type ServerInfo struct {
	Host string
	Port int32
}

// Serve : serve http requests
func Serve(t ServerInfo) {
	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json; charset=UTF-8")
	r.HandleFunc("/getAlbum", handler.GetAlbum).Methods("GET")
	r.HandleFunc("/uploadPhoto", handler.UploadPhoto).Methods("POST")
	r.HandleFunc("/createAlbum", handler.CreateAlbum).Methods("POST")
	r.HandleFunc("/deleteAlbum/{albumID}", handler.DeleteAlbum).Methods("DELETE")
	r.HandleFunc("/deletePhoto/{photoID}", handler.DeletePhoto).Methods("DELETE")
	fs := http.FileServer(http.Dir(configure.AppProperties.StaticFilePath))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	url := fmt.Sprintf("%s:%d", t.Host, t.Port)
	if err := http.ListenAndServe(url, r); err != nil {
		log.Println("Failed to serve at ", url)
	}
	log.Println("Serving at ", url)
}
