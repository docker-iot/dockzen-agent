package lib

import (
	"fmt"
	dockzen_h "include"
)

/*
#include <dockzen.h>

void _C_CallbackContainerUpdate(int status, void* user_data);

typedef struct{
	void * cb_fcn;
}update_cb;

*/
import "C"

type ContainerUpdateCB func(int)

type UserDataCB struct {
	fn ContainerUpdateCB
}

var fn_cb ContainerUpdateCB

func GetContainerListsInfo(containers_info *dockzen_h.Containers_info) int {

	fmt.Println(">>>>>>>>>> GetContainerListsInfo()...")

	var C_containers_info C.containers_info_s
	var ret = C.dockzen_get_containers_info(&C_containers_info)

	fmt.Println("ret = ", ret)

	if ret == 0 {
		containers_info.Count = int(C_containers_info.count)

		for i := 0; i<containers_info.Count; i++ {
			containers_info.Containerinfo[i].ID = C.GoString(C_containers_info.container[i].id)
			containers_info.Containerinfo[i].Name = C.GoString(C_containers_info.container[i].name)
			containers_info.Containerinfo[i].ImageName = C.GoString(C_containers_info.container[i].image_name)
			containers_info.Containerinfo[i].Status = C.GoString(C_containers_info.container[i].status)
		}
	}

	fmt.Println("container = ", containers_info)
	return int(ret)
}

//export _GO_CallbackContainerUpdate
func _GO_CallbackContainerUpdate(status C.int, user_data C.int) {
	fmt.Println("updateCallback status = ", status)

	fn_cb(int(status))
}

func UpdateContainer(container_update dockzen_h.ContainerUpdateInfo, callback ContainerUpdateCB) int {
	fmt.Println(">>>>>>>>>> UpdateContainer()...")

	//var C_update_info C.container_update_s
	//C_update_info.id = C.CString("container_update.ID")

	//var user_data UserDataCB
	//user_data.fn = callback

	//fn_cb = callback

	//var ret = C.dockzen_update_container(&C_update_info, ((C.container_update_cb)(unsafe.Pointer(C._C_CallbackContainerUpdate))), unsafe.Pointer(&user_data))

	//fmt.Println("ret = ", ret)
	var ret int

	return int(ret)

}
