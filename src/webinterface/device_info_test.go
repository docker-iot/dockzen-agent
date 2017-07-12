package webinterface

import (
  "testing"
  "log"
)

/**
 * @fn	TestGetHardwareAddress(t *testing.T)
 * @brief unit test function.
 *
 * @param	t, [in] testing structure
*/
func TestGetHardwareAddress(t *testing.T){
  log.Printf("[TEST] ========== GetHardwareAddress test code ===========")
  address, r := GetHardwareAddress()
  if r != nil {
    t.Errorf("[TEST] HardwareAddress error = ", r)
  }else {
    log.Printf("[TEST] HardwareAddress = ", address)
  }
}
