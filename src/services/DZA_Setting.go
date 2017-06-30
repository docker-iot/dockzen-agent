package services

import (
  "log"
  "io/ioutil"
  "encoding/json"
)

var SERVER_URL_FILE_PATH  = "data/server_url.json"

func GetServerURL(config_path string) string {

  if config_path == "" {
      config_path = SERVER_URL_FILE_PATH
  }
  data, err := ioutil.ReadFile(config_path)

  if err != nil {
    log.Printf("[%s] server_url.config file error!", __FILE__)
    return ""
  }
  var config server_config

  json.Unmarshal([]byte(data), &config)
  log.Printf("[%s] Server URL = %s", __FILE__, config.Server_URL)

  return config.Server_URL
}
