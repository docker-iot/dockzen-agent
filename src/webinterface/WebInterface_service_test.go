package webinterface

import (
  "log"
  "testing"
  dockzen_h "include"
)

/**
 * @fn	TestParseUpdateParam(t *testing.T)
 * @brief unit test function.
 *
 * @param	t, [in] testing structure
*/
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

/**
 * @fn	TestWS_GetContainerLists_Res(t *testing.T)
 * @brief unit test function.
 *
 * @param	t, [in] testing structure
*/
func TestWS_GetContainerLists_Res(t *testing.T){
  log.Printf("[TEST] ========== WS_GetContainerLists test code ===========")
  var send_info ws_ContainerList_info
  var containersInfo dockzen_h.Containers_info
  containersInfo.Count = 1
  containersInfo.Containerinfo = make([]dockzen_h.Container, 1)
  containersInfo.Containerinfo[0].ID = "1234"
  containersInfo.Containerinfo[0].Name = "tizen"
  containersInfo.Containerinfo[0].ImageName = "headless:v0.1"
  containersInfo.Containerinfo[0].Status = "running"

  err := WS_GetContainerLists_Res(&send_info, containersInfo)

  if err != nil {
    t.Errorf("[TEST] WS_GetContainerList_Res error")
  }else {
    log.Printf("[TEST] WS_GetContainerList_Res = ", send_info)
  }
}

/**
 * @fn	TestWS_UpdateImage_Res(t *testing.T)
 * @brief unit test function.
 *
 * @param	t, [in] testing structure
*/
func TestWS_UpdateImage_Res(t *testing.T){
  log.Printf("[TEST] ========== WS_UpdateImage test code ===========")
  var send_info ws_ContainerUpdateReturn
  var updateReturn dockzen_h.ContainerUpdateRes
  updateReturn.Container_Name = "tizen"
  updateReturn.Image_name_Prev = "headless:v0.1"
  updateReturn.Image_name_New = "headless:v0.2"
  updateReturn.Status = "Running"
  err := WS_UpdateImage_Res(&send_info, updateReturn)

  if err != nil {
    t.Errorf("[TEST] WS_UpdateImage_Res error")
  } else {
    log.Printf("[TEST] WS_UpdateImage_Res = ", send_info)
  }
}
