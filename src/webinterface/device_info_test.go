package webinterface

import (
  "testing"
  "log"
)

func TestGetHardwareAddress(t *testing.T){
  log.Printf("[TEST] ========== GetHardwareAddress test code ===========")
  address, r := GetHardwareAddress()
  if r != nil {
    t.Errorf("[TEST] HardwareAddress error = ", r)
  }else {
    log.Printf("[TEST] HardwareAddress = ", address)
  }
}
