package services

import (
  "testing"
  "log"
)

func TestSetServerURL(t *testing.T){
  var SERVER_URL = "10.113.62.204:4000"
  r := SetServerURL(SERVER_URL)

  if r == 0 {
    log.Printf("[TEST] Set ServerURL OK!!")
  } else {
    t.Errorf("[TEST] SetServerURL FAIL!!")
  }
}

func TestGetServerURL(t *testing.T) {
  log.Printf("[TEST] ========== GetServerURL test code ===========")
  r := GetServerURL()
  if r == "" {
    t.Errorf("[TEST] server URL ret = ", r)
  } else {
    log.Printf("[TEST] server_URL = %s", r)
  }
}
