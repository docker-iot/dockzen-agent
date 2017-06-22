package webinterface

import (
	"log"
)

type ConnectedResp struct {
	Cmd       string `json:"cmd"`
	Token     string `json:"token"`
	Clinetnum int    `json:"clientnum"`
}

func ExampleWebService_Update(){
	//var containersInfo dockzen_h.Containers_info
	//var ret = services.DZA_Mon_GetContainersInfo(&containersInfo)

	//fmt.Println("containerInfo = ", containersInfo)
	//fmt.Println("ret = ", ret)
	 //send_info, ret := wsGetContainerLists()

	 //fmt.Println("containerinfo ws_test = ", send_info)

	var data dockzen_h.ContainerUpdateInfo

	data.Container_Name = "tizen_ksy"
	data.Image_Name = "10.113.62.204:443/headless:v0.2"

	send_update, ret := wsUpdateImage(data)
	log.Printf("update ws_test = ", send_update)
	log.Printf("ret =", ret)
}
