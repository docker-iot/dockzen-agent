package services

import (
  "log"
)

func ExampleGetServerURL() {
  log.Printf("[%s] GetServerURL test code", __FILE__)
  log.Printf("[%s] server_URL = %s", __FILE__, GetServerURL("../../data/server_url.json"))
  // Output:
  // .
}
