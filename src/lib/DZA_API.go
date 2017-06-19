package lib

import (
	"fmt"
	dockzen_h "include"
	"unsafe"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. ${SRCDIR}/libdockzen.a ${SRCDIR}/libjson-c.a
#include <dockzen.h>

void _C_CallbackContainerUpdate(int status);

*/
import "C"

type ContainerUpdateCB func(int)

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
func _GO_CallbackContainerUpdate(c_status C.int, callback unsafe.Pointer) {
	update_callback := *(*func(int))(callback)
	update_callback(int(c_status))
}

func UpdateContainer(container_update dockzen_h.ContainerUpdateInfo, callback ContainerUpdateCB) int {
	fmt.Println(">>>>>>>>>> UpdateContainer()...")
	var C_update_info C.container_update_s
	C_update_info.id = C.CString("container_update.ImageName")

	var ret = C.dockzen_update_container(&C_update_info, (C.container_update_cb)(unsafe.Pointer(C._C_CallbackContainerUpdate)), (unsafe.Pointer(&callback)))

	return int(ret)

}
