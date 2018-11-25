package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/sparcs-home-go/internal/app/configure"
	"github.com/sparcs-home-go/internal/app/service"
)

var (
	key   = []byte(configure.AppProperties.CookieSecretKey)
	store = sessions.NewCookieStore(key)
)

const ssoSession = "sso_session"

// SSOLogin : login thr sso
func SSOLogin(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, ssoSession)
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		nextURL := session.Values["next"].(string)
		if nextURL == "" {
			nextURL = configure.AppProperties.LoginRedirectURL
		}
		http.Redirect(w, r, nextURL, 301)
		return
	}
	loginURL, state := service.GetLoginParams()
	session.Values["ssoState"] = state
	session.Save(r, w)
	http.Redirect(w, r, loginURL, 301)
}

// SSOLoginCallback : callback from SSO
func SSOLoginCallback(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, ssoSession)
	prevState, ok := session.Values["ssoState"].(string)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden session"))
		return
	}
	ssoState := r.FormValue("state")
	if prevState != ssoState {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Forbidden session"))
		return
	}
	code := r.FormValue("code")
	sso, err := service.GetUserInfo(code)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Not registered in SSO"))
		return
	}
	session.Values["authenticated"] = true
	session.Values["sid"] = sso["sid"]
	sparcsID := sso["sparcs_id"]
	if sparcsID != "" {
		session.Values["sparcsID"] = sparcsID
		session.Values["isSPARCS"] = true
		// handle admin
	} else {
		session.Values["isSPARCS"] = false
	}
	redirectURL := session.Values["next"].(string)
	if redirectURL == "" {
		redirectURL = configure.AppProperties.LoginRedirectURL
	} else {
		delete(session.Values, "next")
	}
	session.Save(r, w)
	http.Redirect(w, r, redirectURL, 301)
}

// SSOLogout : logout from SSO
func SSOLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, ssoSession)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		log.Println("logout by unauthenticated user")
		http.Redirect(w, r, configure.AppProperties.LogoutRedirectURL, 301)
		return
	}
	sid := session.Values["sid"].(string)
	logoutURL := service.GetLogoutURL(sid)
	session.Options.MaxAge = -1
	session.Save(r, w)
	log.Println("Logout from sparcs.org, logoutURL ", logoutURL)
	http.Redirect(w, r, logoutURL, 301)
}
