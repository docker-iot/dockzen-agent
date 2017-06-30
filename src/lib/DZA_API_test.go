package lib

import (
  "testing"
  "log"
  dockzen_h "include"
  "unsafe"
)


func TestGetContainerListsInfo(t *testing.T){
  log.Printf("[TEST] ========== GetContainerListsInfo test code ===========")
  testGetContainerListsInfo(t)
}

func Updatecallback_test(status dockzen_h.Container_update_cb_s, user_data unsafe.Pointer) {
	log.Printf("[TEST] Updatecallback_test OK!!!!")
	log.Printf("[TEST] Updatecallback_test> status.Container_name = ", status.Container_name)
	log.Printf("[TEST] Updatecallback_test> status.Image_name = ", status.Image_name)
	log.Printf("[TEST] Updatecallback_test> status.Status = ", status.Status)
}

func TestUpdateContainer(t *testing.T){
  log.Printf("[TEST] ========== UpdateContainer test code ===========")
  testUpdateContainer(t, Updatecallback_test)
}
