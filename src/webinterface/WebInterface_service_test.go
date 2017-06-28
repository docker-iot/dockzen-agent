package webinterface

import (
  "log"
  "testing"
  dockzen_h "include"
  services "services"
)

func TestParseUpdateParam(t *testing.T){
  log.Printf("[TEST] ========== ParseUpdateParam test code ===========")

  msg := `{"ImageName":"tizen","Name":"test"}`
  updateinfo, r := ParseUpdateParam(msg)
  if r != nil {
    t.Errorf("[TEST] ParseUpdateParam message error =", r)
  } else {
    log.Printf("[TEST] updateinfo =", updateinfo)
  }
}

func TestWS_GetContainerLists(t *testing.T){
  log.Printf("[TEST] ========== WS_GetContainerLists test code ===========")
  var containersInfo dockzen_h.Containers_info
	var ret = services.DZA_Mon_GetContainersInfo(&containersInfo)

  if ret != 0 {
    t.Errorf("[TEST] GetContainerLists error = ", ret)
  } else{
    log.Printf("[TEST] GetContainerLists = ", containersInfo)
  }
}

func TestWS_UpdateImage(t *testing.T){
  log.Printf("[TEST] ========== WS_UpdateImage test code ===========")
  var data dockzen_h.ContainerUpdateInfo
  data.Image_Name = "headless:v1.0"
  data.Container_Name = "tizen"
  update_Return, ret := WS_UpdateImage(data)

  if ret != 0 {
    t.Errorf("[TEST] UpdateImage error = ", ret)
  } else {
    log.Printf("[TEST] UpdateImage return = ", update_Return)
  }

}
