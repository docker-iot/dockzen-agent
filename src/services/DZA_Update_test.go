package services

import (
  "testing"
  "log"
  "unsafe"
  dockzen_h "include"
)

/**
 * @fn	Test__Updatecallback(t *testing.T)
 * @brief unit test function.
 *
 * @param	t, [in] testing structure
*/
func Test__Updatecallback(t *testing.T){
  log.Printf("[TEST] ========== Updatecallback test code ===========")
  var user_data update_userData
  var status dockzen_h.Container_update_cb_s
  user_data.Container_Name = "tizen_test"
  status.Container_name = "tizen_test"
  status.Image_name = "test_headless:v0.1"
  status.Status = "running"

  __Updatecallback(status ,unsafe.Pointer(&user_data))

}
