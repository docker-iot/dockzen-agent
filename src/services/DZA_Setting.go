package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var SERVER_URL_FILE_PATH = "data"
var SERVER_URL_FILE = "server_url.json"
var DEFAULT_SERVER_URL = "13.124.64.10:80/ws"

// SetServerURL set web server url.
// Param url is server address.
// This function returns result of function.
func SetServerURL(url string) int {
	var config server_config
	config.Server_URL = url

	if _, err := os.Stat(SERVER_URL_FILE_PATH); os.IsNotExist(err) {
		err = os.MkdirAll(SERVER_URL_FILE_PATH, 0755)
		if err != nil {
			log.Printf("[%s] SetServerURL %s folder create error!!", __FILE__, SERVER_URL_FILE_PATH)
			return -1
		}
	}

	f, err := os.Create(SERVER_URL_FILE_PATH + "/" + SERVER_URL_FILE)
	if err != nil {
		log.Printf("[%s] SetServerURL file create error!!!", __FILE__)
		return -1
	}

	server_url, err := json.Marshal(config)
	log.Printf("[%s] SetServerURL server_url=", __FILE__, string(server_url))

	_, err = f.Write(server_url)
	if err != nil {
		log.Printf("[%s] SetServerURL file write error", __FILE__)
		return -1
	}

	defer f.Close()

	return 0

}

// GetServerURL returns web server url.
func GetServerURL() string {

	if _, err := os.Stat(SERVER_URL_FILE_PATH + "/" + SERVER_URL_FILE); os.IsNotExist(err) {
		// data/server_url.json does not exist
		log.Printf("[%s] SetServerURL!!!", __FILE__)
		SetServerURL(DEFAULT_SERVER_URL)
	}

	data, err := ioutil.ReadFile(SERVER_URL_FILE_PATH + "/" + SERVER_URL_FILE)

	if err != nil {
		log.Printf("[%s] GetServerURL file error!", __FILE__)
		return ""
	}
	var config server_config

	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Printf("[%s] GetServerURL URL error!!!!", __FILE__)
		return ""
	}

	log.Printf("[%s] Server URL = %s", __FILE__, config.Server_URL)

	return config.Server_URL
}
