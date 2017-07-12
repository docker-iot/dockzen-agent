package services

/**
 * @struct update_userData
 * @brief This structure contains user data information
 *
 * The containers struct encapsulate user data in the one data
*/
type update_userData struct{
	Container_Name string `json:"Name"`
}

/**
 * @struct server_config
 * @brief This structure contains web server url information
 *
 * The containers struct encapsulate server configuration in the one data
*/
type server_config struct{
	Server_URL	string `json:"URL"`
}

var __FILE__ = "SERVICES"
