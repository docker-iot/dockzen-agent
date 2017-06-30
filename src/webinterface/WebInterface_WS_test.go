package webinterface

import (
  "log"
  "testing"
)

func TestWS_Server_Connect(t *testing.T){
  log.Printf("[TEST] ========== WS_Server_Connect test code ===========")
  ws, err := WS_Server_Connect("10.113.62.204:4000")

  if ws == nil || err != nil {
    t.Errorf("[TEST] WS_Server connection error")
  } else {
    log.Printf("[TEST] Server Connection !!")
    ws.Close()
  }
}
