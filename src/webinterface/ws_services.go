package webinterface

import (
  "fmt"
  dockzen_h "include"
  "services"
)

func wsGetContainerLists() (ws_ContainerList_info, int) {
	var containersInfo dockzen_h.Containers_info
  var send_info ws_ContainerList_info
	var ret = services.DZA_Mon_GetContainersInfo(&containersInfo)

	if ret != 0 {
		fmt.Println("GetContainersInfo error = ", ret)
	} else {

    var err error
		send_info.Cmd = "GetContainersInfo"
		send_info.ContainerCount = int(containersInfo.Count)
		send_info.DeviceID, err = GetHardwareAddress()
    if err != nil{
      ret = -1
      fmt.Println("HardwareAddress error = ", err)
    }

		fmt.Println("DevicedID = ", send_info.DeviceID)

		for i := 0; i < send_info.ContainerCount; i++ {
			send_info.Container = append(send_info.Container, containersInfo.Containerinfo[i]);
		}

		fmt.Println("ContainerInfo -> ", send_info)
	}

	return send_info, ret
}

func wsUpdateImage(data dockzen_h.ContainerUpdateInfo) (ws_ContainerUpdateReturn, int) {
	var updateReturn dockzen_h.ContainerUpdateRes
  var send_Return ws_ContainerUpdateReturn
	var ret = services.DZA_Update_Do(data, &updateReturn)

	fmt.Println("wsUpdateImage> updateReturn->status = ", updateReturn.Status)

	if ret != 0{
		fmt.Println("UpdateInfo error = ", ret)
	} else {

    var err error
		send_Return.Cmd = "UpdateImage"
		send_Return.DeviceID, err = GetHardwareAddress()
		send_Return.UpdateState = updateReturn.Status

    if err != nil {
      ret = -1
      fmt.Println("HardwareAddress error = ", err)
    }
		fmt.Println("wsUpdateImage> send_Return = ", send_Return)
	}

	return send_Return, ret
}
