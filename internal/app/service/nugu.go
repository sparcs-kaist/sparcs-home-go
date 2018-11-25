package service

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/sparcs-home-go/internal/app/configure"
)

// NuguUserInfo : nugu collected info
type NuguUserInfo struct {
	Name            string `json:"name"`
	IsPrivate       bool   `json:"is_private"`
	IsDeveloper     bool   `json:"is_developer"`
	IsDesigner      bool   `json:"is_designer"`
	IsUndergraduate bool   `json:"is_undergraduate"`
	EntYear         int    `json:"ent_year"`
	Org             string `json:"org"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Birth           string `json:"birth"`
	Dorm            string `json:"dorn"`
	Lab             string `json:"lab"`
	HomeAdd         string `json:"home_add"`
	GithubID        string `json:"github_id"`
	LinkedinURL     string `json:"linkedin_url"`
	BehanceURL      string `json:"behance_url"`
	FacebookID      string `json:"facebook_id"`
	TwitterID       string `json:"twitter_id"`
	BattlenetID     string `json:"battlenet_id"`
	Website         string `json:"website"`
	Blog            string `json:"blog"`
}

func requestData(requestType string, url string, jsonData string, auth bool) (string, error) {
	var body io.Reader
	if jsonData != "" {
		body = strings.NewReader(jsonData)
	}
	req, err := http.NewRequest(requestType, url, body)
	if err != nil {
		log.Println("Failed while creating new request ", err)
		return "", err
	}
	if auth {
		req.SetBasicAuth(configure.AppProperties.NuguID, configure.AppProperties.NuguPassword)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		log.Println("Failed to ", err)
		return "", err
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read data from request \n err: ", err)
		return "", err
	}
	return string(data), nil
}

// GetNuguPublicUsers : get public user list
func GetNuguPublicUsers() (string, error) {
	requestURL := configure.AppProperties.NuguServiceURL + "/public_users"
	res, err := requestData("GET", requestURL, "", false)
	if err != nil {
		log.Println("Failed to get public users ", err)
		return "", err
	}
	return res, nil
}

// GetNuguUsers : get nugu info of all users
func GetNuguUsers() (string, error) {
	requestURL := configure.AppProperties.NuguServiceURL + "/users"
	res, err := requestData("GET", requestURL, "", true)
	if err != nil {
		log.Println("Failed to get all users ", err)
		return "", err
	}
	return res, nil
}

// GetNuguUserInfo : get nugu info of user
func GetNuguUserInfo(userID string) (string, error) {
	requestURL := configure.AppProperties.NuguServiceURL + "/users/" + userID
	res, err := requestData("GET", requestURL, "", true)
	if err != nil {
		log.Println("Failed to get user info of "+userID+",\n err: ", err)
		return "", err
	}
	return res, nil
}

// UpdateNuguUserInfo : update nugu info
func UpdateNuguUserInfo(userID string, info NuguUserInfo) (string, error) {
	requestURL := configure.AppProperties.NuguServiceURL + "/users/" + userID
	jsonData, err := json.Marshal(info)
	res, err := requestData("PUT", requestURL, string(jsonData), true)
	if err != nil {
		log.Println("Failed to update user info of "+userID+",\n err: ", err)
		return "", err
	}
	return res, nil
}
