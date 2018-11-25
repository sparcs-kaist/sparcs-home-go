package configure

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	key = []byte(AppProperties.CookieSecretKey)
	// Store : session store
	Store = sessions.NewCookieStore(key)
)

const ssoSession = "sso_session"

// GetSession : get session from key store
func GetSession(r *http.Request) *sessions.Session {
	session, err := Store.Get(r, ssoSession)
	if err != nil {
		log.Println("Failed to get session")
	}
	return session
}
