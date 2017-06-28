package webinterface

import (
  "log"
  dockzen_h "include"
  "services"
  "encoding/json"
)

func WS_GetContainerLists() (ws_ContainerList_info, int) {
	var containersInfo dockzen_h.Containers_info
  var send_info ws_ContainerList_info
	var ret = services.DZA_Mon_GetContainersInfo(&containersInfo)

	if ret != 0 {
		log.Printf("[%s] GetContainersInfo error = ", __FILE__, ret)
	} else {

    var err error
		send_info.Cmd = "GetContainersInfo"
		send_info.ContainerCount = int(containersInfo.Count)
		send_info.DeviceID, err = GetHardwareAddress()
    if err != nil{
      ret = -1
      log.Printf("[%s] HardwareAddress error = ", __FILE__, err)
    }

		log.Printf("[%s] DevicedID = ", __FILE__, send_info.DeviceID)

		for i := 0; i < send_info.ContainerCount; i++ {
			send_info.Container = append(send_info.Container, containersInfo.Containerinfo[i]);
		}

		log.Printf("[%s] ContainerInfo -> ", __FILE__, send_info)
	}

	return send_info, ret
}

func ParseUpdateParam(msg string) (dockzen_h.ContainerUpdateInfo, error) {
	send := dockzen_h.ContainerUpdateInfo{}
	r := json.Unmarshal([]byte(msg), &send)
  if r == nil {
	   log.Printf("[%s] parsing ContainerName: ", __FILE__, send.Container_Name)
	 log.Printf("[%s] parsing ImageName: ", __FILE__, send.Image_Name)
  }

	return send, r
}

func WS_UpdateImage(data dockzen_h.ContainerUpdateInfo) (ws_ContainerUpdateReturn, int) {
	var updateReturn dockzen_h.ContainerUpdateRes
  var send_Return ws_ContainerUpdateReturn
	var ret = services.DZA_Update_Do(data, &updateReturn)

	log.Printf("[%s] updateReturn->status = ", __FILE__, updateReturn.Status)

	if ret != 0{
		log.Printf("[%s] UpdateInfo error = ", __FILE__, ret)
	} else {
    var err error
		send_Return.Cmd = "UpdateImage"
		send_Return.DeviceID, err = GetHardwareAddress()
		send_Return.UpdateState = updateReturn.Status

    if err != nil {
      ret = -1
      log.Printf("[%s] HardwareAddress error = ", __FILE__, err)
    }
		log.Printf("[%s] wsUpdateImage> send_Return = ", __FILE__, send_Return)
	}

	return send_Return, ret
}
