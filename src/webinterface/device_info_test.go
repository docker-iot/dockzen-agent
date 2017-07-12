package webinterface

import (
  "testing"
  "log"
)

/**
 * @fn	TestgetHardwareAddress(t *testing.T)
 * @brief unit test function.
 *
 * @param	t, [in] testing structure
*/
func TestgetHardwareAddress(t *testing.T){
  log.Printf("[TEST] ========== getHardwareAddress test code ===========")
  address, r := getHardwareAddress()
  if r != nil {
    t.Errorf("[TEST] HardwareAddress error = ", r)
  }else {
    log.Printf("[TEST] HardwareAddress = ", address)
  }
}
