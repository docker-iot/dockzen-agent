package services

import (
  "testing"
  "log"
  dockzen_h "include"
)

func TestDZA_Mon_GetContainersInfo(t *testing.T){
  log.Printf("[TEST] ========== DZA_Mon_GetContainersInfo test code ===========")
  var containersInfo dockzen_h.Containers_info

  r := DZA_Mon_GetContainersInfo(&containersInfo)
  if r != 0 {
    t.Errorf("[TEST] GetContainerListsInfo r = ", r)
  } else {
    log.Printf("[TEST] GetContainerListsInfo result= ", containersInfo)
  }
}
