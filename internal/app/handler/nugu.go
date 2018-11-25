package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sparcs-home-go/internal/app/configure"
	"github.com/sparcs-home-go/internal/app/service"
)

// GetNuguPublicUsers : get nugu public user list
func GetNuguPublicUsers(w http.ResponseWriter, r *http.Request) {
	res, err := service.GetNuguPublicUsers()
	if err != nil {
		log.Println("Failed to get nugu public users \n err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

// GetNuguUsers : get nugu user list, auth required
func GetNuguUsers(w http.ResponseWriter, r *http.Request) {
	session := configure.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("logout by unauthenticated user")
		http.Redirect(w, r, configure.AppProperties.LogoutRedirectURL, 301)
		return
	}
	res, err := service.GetNuguUsers()
	if err != nil {
		log.Println("Failed to get nugu users \n err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

// GetNuguUserInfo : get nugu user info, req param userID, auth required
func GetNuguUserInfo(w http.ResponseWriter, r *http.Request) {
	session := configure.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("logout by unauthenticated user")
		http.Redirect(w, r, configure.AppProperties.LogoutRedirectURL, 301)
		return
	}
	req := mux.Vars(r)
	res, err := service.GetNuguUserInfo(req["userID"])
	if err != nil {
		log.Println("Failed to get nugu users \n err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

// UpdateNuguUserInfo : update nugu user info, req param userID, auth required
func UpdateNuguUserInfo(w http.ResponseWriter, r *http.Request) {
	session := configure.GetSession(r)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("logout by unauthenticated user")
		http.Redirect(w, r, configure.AppProperties.LogoutRedirectURL, 301)
		return
	}
	reqParams := mux.Vars(r)
	var reqBody service.NuguUserInfo
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	res, err := service.UpdateNuguUserInfo(reqParams["userID"], reqBody)
	if err != nil {
		log.Println("Failed to update nugu user info \n err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error occurred"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
