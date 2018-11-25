package service

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sparcs-home-go/internal/app/configure"
)

const (
	possible        = "0123456789abcdef"
	urlTokenRequire = "token/require/"
	urlTokenInfo    = "token/info/"
	urlLogout       = "logout/"
	urlUnregister   = "unregister/"
)

var (
	clientID string
)

// SSOConfig :
type SSOConfig struct {
	ClientID  string
	SecretKey string
}

// SSOProperties :
var SSOProperties SSOConfig

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// SetProperties :
func SetProperties(config SSOConfig) {
	SSOProperties = config
}

func tokenHex(size int) string {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = possible[rand.Intn(16)]
	}
	return string(buf)
}

func urlEncode(payload map[string]string) string {
	path := "?"
	for k, v := range payload {
		path += (k + "=" + url.QueryEscape(v) + "&")
	}
	return path
}

func postData(postURL string, form map[string][]string) (map[string]string, error) {
	resp, err := http.PostForm(postURL, url.Values(form))
	if err != nil {
		log.Println(err)
	}
	x := make(map[string]string)
	err = json.NewDecoder(resp.Body).Decode(&x)
	if err != nil {
		fmt.Printf("Error while decoding json: %s\n", err)
		return nil, err
	}
	return x, nil
}

func signPayload(payload []string) (string, string) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	payloadWithTimestamp := append(payload, timestamp)
	msg := strings.Join(payloadWithTimestamp, "")
	h := hmac.New(md5.New, []byte(SSOProperties.SecretKey))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil)), timestamp
}

// GetLoginParams : get login url, token
func GetLoginParams() (string, string) {
	state := tokenHex(10)
	payload := map[string]string{
		"client_id": SSOProperties.ClientID,
		"state":     state,
	}
	loginURL := urlTokenRequire + urlEncode(payload)
	return loginURL, state
}

// GetUserInfo : get info from token
func GetUserInfo(code string) (map[string]string, error) {
	sign, timestamp := signPayload([]string{code})
	params := map[string][]string{
		"client_id": {SSOProperties.ClientID},
		"code":      {code},
		"timestamp": {timestamp},
		"sign":      {sign},
	}
	res, err := postData(urlTokenInfo, params)
	if err != nil {
		log.Println("Failed to get token info", err)
		return nil, err
	}
	return res, nil
}

// GetLogoutURL : logout url from service id
func GetLogoutURL(sid string) string {
	sign, timestamp := signPayload([]string{sid, configure.AppProperties.LogoutRedirectURL})
	payload := map[string]string{
		"client_id":   SSOProperties.ClientID,
		"sid":         sid,
		"timestamp":   timestamp,
		"redirectUri": configure.AppProperties.LogoutRedirectURL,
		"sign":        sign,
	}
	logoutURL := urlTokenRequire + urlEncode(payload)
	return logoutURL
}
