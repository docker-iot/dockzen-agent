// Package service consists of services provided by the agent.
package services

// Structure update_userData contains user data information.
type update_userData struct{
	Container_Name string `json:"Name"`
}

//Structure server_config contains web server url information.
type server_config struct{
	Server_URL	string `json:"URL"`
}

var __FILE__ = "SERVICES"
