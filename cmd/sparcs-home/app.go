package main

import (
	"encoding/json"
	"log"

	"github.com/sparcs-home-go/internal/app"
	"github.com/sparcs-home-go/internal/app/configure"
	"github.com/sparcs-home-go/internal/utils"
)

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database configure.DatabaseInfo `json:"Database"`
	App      configure.AppConfig    `json:"App"`
	Server   app.ServerInfo         `json:"Server"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

func init() {
	utils.Load("config.json", config)
	log.Println("config: ", config.Database.User)
}

func main() {
	// Load the configuration file
	configure.ConnectDB(config.Database)
	configure.SetProperties(config.App)
	app.Serve(config.Server)
}
