package services

import (
  "testing"
  "log"
)

func TestServerURL(t *testing.T) {
  log.Printf("[TEST] ========== GetServerURL test code ===========")
  r := GetServerURL("../../data/server_url.json")

  if r == "" {
    t.Errorf("[TEST] server URL ret = ", r)
  } else {
    log.Printf("[TEST] server_URL = %s", r)
  }
}
