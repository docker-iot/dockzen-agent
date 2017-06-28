package services

import (
  "testing"
  "log"
  dockzen_h "include"
)

func TestDZA_Update_Do(t *testing.T){
  log.Printf("[TEST] ========== DZA_Update_Do test code ===========")
  var updateReturn dockzen_h.ContainerUpdateRes
  var updateinfo dockzen_h.ContainerUpdateInfo
  updateinfo.Image_Name = "tizen_headless:v0.2"
  updateinfo.Container_Name = "tizen"


  var ret = DZA_Update_Do(updateinfo, &updateReturn)


  if ret != 0{
    t.Errorf("[TEST] UpdateContainer ret =", ret)
  }else{
    log.Printf("[TEST] updateReturn!!!!")
    log.Printf("[TEST] Container_name = ",updateReturn.Container_Name)
    log.Printf("[TEST] Image_name_prev = ", updateReturn.Image_name_Prev)
    log.Printf("[TEST] Image_name_new = ", updateReturn.Image_name_New)
    log.Printf("[TEST] status = ", updateReturn.Status)
  }
}
