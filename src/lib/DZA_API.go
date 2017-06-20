package lib

import (
	"fmt"
	dockzen_h "include"
	"unsafe"
)

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. ${SRCDIR}/libdockzen.a ${SRCDIR}/libjson-c.a

#include <stdlib.h>
#include <dockzen.h>

typedef struct{
	void *callback;
	void* user_data;
}token_s;

void* _C_SetCallbackStruct( void* callback, void* user_data);
void _C_CallbackContainerUpdate(int status);

*/
import "C"

type ContainerUpdateCB func(dockzen_h.Container_update_cb_s, unsafe.Pointer)

func GetContainerListsInfo(containers_info *dockzen_h.Containers_info) int {

	fmt.Println(">>>>>>>>>> GetContainerListsInfo()...")

	var C_containers_info C.containers_info_s
	var ret = C.dockzen_get_containers_info(&C_containers_info)

	fmt.Println("ret = ", ret)

	if ret == 0 {
		containers_info.Count = int(C_containers_info.count)
		var container dockzen_h.Container
		for i := 0; i<containers_info.Count; i++ {
			container = dockzen_h.Container{
				ID: C.GoString(C_containers_info.container[i].id),
				Name: C.GoString(C_containers_info.container[i].name),
				ImageName: C.GoString(C_containers_info.container[i].image_name),
				Status: C.GoString(C_containers_info.container[i].status),
			}
			containers_info.Containerinfo = append(containers_info.Containerinfo, container)

			C.free(unsafe.Pointer(C_containers_info.container[i].id))
			C.free(unsafe.Pointer(C_containers_info.container[i].name))
			C.free(unsafe.Pointer(C_containers_info.container[i].image_name))
			C.free(unsafe.Pointer(C_containers_info.container[i].status))
		}
	}

	fmt.Println("container = ", containers_info)
	return int(ret)
}

//export _GO_CallbackContainerUpdate
func _GO_CallbackContainerUpdate(c_status_info unsafe.Pointer, userdata unsafe.Pointer) {
	fmt.Println("_GO_CallbackContainerUpdate > !!!")
	C_statusInfo := (*C.container_update_cb_s)(c_status_info)
	var update_status dockzen_h.Container_update_cb_s

	update_status.Container_name = C.GoString(C_statusInfo.container_name)
	update_status.Image_name = C.GoString(C_statusInfo.image_name)
	update_status.Status = C.GoString(C_statusInfo.status)

	defer C.free(unsafe.Pointer(C_statusInfo.container_name))
	defer C.free(unsafe.Pointer(C_statusInfo.image_name))
	defer C.free(unsafe.Pointer(C_statusInfo.status))

	CToken := (*C.token_s)(userdata)
	update_callback := *(*ContainerUpdateCB)(CToken.callback)
	update_callback(update_status, CToken.user_data)

	fmt.Println("CToken pointer = ", &CToken)
	defer C.free(unsafe.Pointer(CToken))
}

func UpdateContainer(container_update dockzen_h.ContainerUpdateInfo, update_res * dockzen_h.ContainerUpdateRes, callback ContainerUpdateCB, userdata unsafe.Pointer) int {
	fmt.Println(">>>>>>>>>> UpdateContainer()...")

	var C_update_info C.container_update_s
	C_update_info.container_name = C.CString(container_update.Container_Name)
	C_update_info.image_name = C.CString(container_update.Image_Name)

	defer C.free(unsafe.Pointer(C_update_info.container_name))
	defer C.free(unsafe.Pointer(C_update_info.image_name))

	var C_update_res C.container_update_res_s

	user_data := unsafe.Pointer(C._C_SetCallbackStruct(unsafe.Pointer(&callback), userdata))

	var ret = C.dockzen_update_container(&C_update_info, &C_update_res, (C.container_update_cb)(unsafe.Pointer(C._C_CallbackContainerUpdate)), unsafe.Pointer(user_data))

	update_res.Container_Name = C.GoString(C_update_res.container_name)
	update_res.Image_name_Prev = C.GoString(C_update_res.image_name_prev)
	update_res.Image_name_New = C.GoString(C_update_res.image_name_new)
	update_res.Status = C.GoString(C_update_res.status)

	defer C.free(unsafe.Pointer(C_update_res.container_name))
	defer C.free(unsafe.Pointer(C_update_res.image_name_prev))
	defer C.free(unsafe.Pointer(C_update_res.image_name_new))
	defer C.free(unsafe.Pointer(C_update_res.status))

	return int(ret)

}
