package webinterface

import (
  "log"
  dockzen_h "include"
  "services"
  "encoding/json"
)

/**
 * @fn	ws_GetContainerLists_Res(send_info * ws_ContainerList_info,
                                containersInfo dockzen_h.Containers_info) (err error)
 * @brief This function set unique device information.
 *
 * @param send_info,      [inout] It is a structure that contains unique device information.
 * @param containersInfo  [in] container information structure
 * @return err,          [out] error value (if the value is null, it is not an error.)
*/
func ws_GetContainerLists_Res(send_info * ws_ContainerList_info, containersInfo dockzen_h.Containers_info) (err error){
  send_info.Cmd = "GetContainersInfo"
  send_info.ContainerCount = int(containersInfo.Count)
  send_info.DeviceID, err = getHardwareAddress()

  if err != nil{
    log.Printf("[%s] HardwareAddress error = ", __FILE__, err)
    return err
  }

  for i := 0; i < send_info.ContainerCount; i++ {
    send_info.Container = append(send_info.Container, containersInfo.Containerinfo[i]);
  }

  log.Printf("[%s] ContainerInfo -> ", __FILE__, send_info)

  return nil
}

/**
 * @fn	ws_GetContainerLists(container_ch Containers_Channel)
 * @brief This function calls sevices.DZA_Mon_GetContainersInfo function..
 *
 * @param container_ch,   [in] It is a channel to communicate with ws_SendMsg function
*/
func ws_GetContainerLists(container_ch Containers_Channel) {
  for{
    msg := <-container_ch.receive
    if msg == true {
    	var containersInfo dockzen_h.Containers_info
      var send_info ws_ContainerList_info
    	var ret = services.DZA_Mon_GetContainersInfo(&containersInfo)
    	if ret != 0 {
    		log.Printf("[%s] GetContainersInfo error = ", __FILE__, ret)
    	} else {
        err := ws_GetContainerLists_Res(&send_info, containersInfo)
        if err == nil {
          container_ch.send <-send_info
        }
      }
    }
  }
}

/**
 * @fn	parseUpdateParam(msg string) (dockzen_h.ContainerUpdateInfo, error)
 * @brief This function convert json data to ContainerUpdateInfo structure.
 *
 * @param msg,                            [in] json data
 * @return dockzen_h.ContainerUpdateInfo  [out] containerUpdateInfo structure
 * @return error                          [out] error value (if the value is null, it is not an error.)
*/
func parseUpdateParam(msg string) (dockzen_h.ContainerUpdateInfo, error) {
	send := dockzen_h.ContainerUpdateInfo{}
	r := json.Unmarshal([]byte(msg), &send)
  if r == nil {
	   log.Printf("[%s] parsing ContainerName: ", __FILE__, send.Container_Name)
	 log.Printf("[%s] parsing ImageName: ", __FILE__, send.Image_Name)
  }

	return send, r
}

/**
 * @fn	ws_UpdateImage_Res(send_Return *ws_ContainerUpdateReturn,
                          updateReturn dockzen_h.ContainerUpdateRes) (err error)
 * @brief This function set unique device information.
 *
 * @param send_Return,   [inout] It is a update information structure that contains unique device information.
 * @param updateReturn  [out] update information
 * @return error        [out] error value (if the value is null, it is not an error.)
*/
func ws_UpdateImage_Res(send_Return *ws_ContainerUpdateReturn, updateReturn dockzen_h.ContainerUpdateRes) (err error){
  send_Return.Cmd = "UpdateImage"
  send_Return.DeviceID, err = getHardwareAddress()
  send_Return.UpdateState = updateReturn.Status

  if err != nil {
    log.Printf("[%s] HardwareAddress error = ", __FILE__, err)
    return err
  }
  log.Printf("[%s] wsUpdateImage> send_Return = ", __FILE__, send_Return)

  return nil
}

/**
 * @fn	ws_UpdateImage(update_ch Update_Channel)
 * @brief This function calls services.DZA_Update_Do function.
 *
 * @param update_ch,   [in] It is a channel to communicate with ws_SendMsg function
*/
func ws_UpdateImage(update_ch Update_Channel){
  for{
    msg := <-update_ch.receive

  	var updateReturn dockzen_h.ContainerUpdateRes
    var send_Return ws_ContainerUpdateReturn
  	var ret = services.DZA_Update_Do(msg, &updateReturn)

  	log.Printf("[%s] updateReturn->status = ", __FILE__, updateReturn.Status)

  	if ret != 0{
  		log.Printf("[%s] UpdateInfo error = ", __FILE__, ret)
  	} else {
      err := ws_UpdateImage_Res(&send_Return, updateReturn)
      if err == nil {
        update_ch.send <- send_Return
      }
  	}
  }
}
