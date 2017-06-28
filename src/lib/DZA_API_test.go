package lib

import (
  "testing"
  "log"
  "unsafe"
  dockzen_h "include"
)

type userData_test struct{
	Container_Name string `json:"Name"`
}

func TestGetContainerListsInfo(t *testing.T){
  log.Printf("[TEST] ========== GetContainersInfo test code ===========")
  var containers dockzen_h.Containers_info

  r := GetContainerListsInfo(&containers)

  if r != 0 {
    t.Errorf("[TEST] GetContainerListsInfo r = ", r)
  } else {
    log.Printf("[TEST] GetContainerListsInfo result= ", containers)
  }
}

func updatecallback_test(status dockzen_h.Container_update_cb_s, user_data unsafe.Pointer){
  log.Printf("[TEST] service Callback OK!!!!")
  log.Printf("[TEST] __updatecallback> status.Container_name = ",status.Container_name)
  log.Printf("[TEST] __updatecallback> status.Image_name = ", status.Image_name)
  log.Printf("[TEST] __updatecallback> status.Status = ", status.Status)

  update_data := *(*userData_test)(user_data)
  log.Printf("[TEST] __updatecallback> user_data.ContainerName = ", update_data)

}

func TestUpdateContainer(t *testing.T){
  log.Printf("[TEST] ========== UpdateContainer test code ===========")
  var updateReturn dockzen_h.ContainerUpdateRes
  var userdata userData_test
  userdata.Container_Name = "tizen"

  var updateinfo dockzen_h.ContainerUpdateInfo
  updateinfo.Image_Name = "tizen_headless:v0.2"
  updateinfo.Container_Name = "tizen"


  var ret = UpdateContainer(updateinfo, &updateReturn, updatecallback_test, unsafe.Pointer(&userdata))


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
