package main

import (
	s3 "app/shared/aws"
	"encoding/json"
	"log"
	"os"
	"runtime"

	"app/route"
	"app/shared/database"
	"app/shared/email"
	"app/shared/jsonconfig"
	"app/shared/server"
	"app/shared/session"
	"app/shared/view"
	"app/shared/view/plugin"
)

// *****************************************************************************
// Application Logic
// *****************************************************************************

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)

	// Connect to database
	database.Connect(config.Database)

	// Email config
	s3.Configure(config.AWS)

	// Email config
	email.Configure(config.Email)

	// Setup the views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime())

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database database.PostGresSQL `json:"Database"`
	AWS      s3.Info              `json:"AWS"`
	Email    email.SMTPInfo       `json:"Email"`
	Server   server.Server        `json:"Server"`
	Session  session.Session      `json:"Session"`
	Template view.Template        `json:"Template"`
	View     view.View            `json:"View"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
