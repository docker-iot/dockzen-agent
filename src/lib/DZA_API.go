package lib

import (
	"log"
	dockzen_h "include"
	"unsafe"
)

/*
//#cgo CFLAGS: -I. ${SRCDIR}/install/include/dockzen/ ${SRCDIR}/install/include/json-c/
#cgo LDFLAGS: -L. ${SRCDIR}/install/arm/lib/libdockzen.a ${SRCDIR}/install/arm/lib/libjson-c.a

#include <stdlib.h>
#include "install/include/dockzen/dockzen.h"

typedef struct{
	void *callback;
	void* user_data;
}token_s;

void* _C_SetCallbackStruct( void* callback, void* user_data);
void _C_CallbackContainerUpdate(int status);

*/
import "C"

type ContainerUpdateCB func(dockzen_h.Container_update_cb_s, unsafe.Pointer)
var __FILE__ = "LIB"

func GetContainerListsInfo(containers_info *dockzen_h.Containers_info) int {

	log.Printf("[%s] >>>>>>>>>> GetContainerListsInfo()...", __FILE__)

	var C_containers_info C.containers_info_s
	var ret = C.dockzen_get_containers_info(&C_containers_info)

	log.Printf("[%s] ret = ", __FILE__, ret)

	if int(ret) == dockzen_h.DOCKZEN_ERROR_NONE {
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

			if unsafe.Pointer(C_containers_info.container[i].id) != nil {
				C.free(unsafe.Pointer(C_containers_info.container[i].id))
				C_containers_info.container[i].id = nil
			}
			if unsafe.Pointer(C_containers_info.container[i].name) != nil {
				C.free(unsafe.Pointer(C_containers_info.container[i].name))
				C_containers_info.container[i].name = nil
			}
			if unsafe.Pointer(C_containers_info.container[i].image_name) != nil {
				C.free(unsafe.Pointer(C_containers_info.container[i].image_name))
				C_containers_info.container[i].image_name = nil
			}
			if unsafe.Pointer(C_containers_info.container[i].status) != nil {
				C.free(unsafe.Pointer(C_containers_info.container[i].status))
				C_containers_info.container[i].status = nil
			}
		}
	}

	log.Printf("[%s] container = ", __FILE__, containers_info)
	return int(ret)
}

//export _GO_CallbackContainerUpdate
func _GO_CallbackContainerUpdate(c_status_info unsafe.Pointer, userdata unsafe.Pointer) {
	log.Printf("[%s] _GO_CallbackContainerUpdate > !!!", __FILE__)
	C_statusInfo := (*C.container_update_cb_s)(c_status_info)
	var update_status dockzen_h.Container_update_cb_s

	update_status.Container_name = C.GoString(C_statusInfo.container_name)
	update_status.Image_name = C.GoString(C_statusInfo.image_name)
	update_status.Status = C.GoString(C_statusInfo.status)

	CToken := (*C.token_s)(userdata)
	update_callback := *(*ContainerUpdateCB)(CToken.callback)
	update_callback(update_status, CToken.user_data)

	defer func() {
		if unsafe.Pointer(CToken) != nil {
			C.free(unsafe.Pointer(CToken))
			CToken = nil

			if unsafe.Pointer(C_statusInfo.container_name) != nil {
				C.free(unsafe.Pointer(C_statusInfo.container_name))
				C_statusInfo.container_name = nil
			}
			if unsafe.Pointer(C_statusInfo.image_name) != nil {
				C.free(unsafe.Pointer(C_statusInfo.image_name))
				C_statusInfo.image_name = nil
			}
			if unsafe.Pointer(C_statusInfo.status) != nil {
				C.free(unsafe.Pointer(C_statusInfo.status))
				C_statusInfo.status = nil
			}
		}
	}()

}

func UpdateContainer(container_update dockzen_h.ContainerUpdateInfo, update_res * dockzen_h.ContainerUpdateRes, callback ContainerUpdateCB, userdata unsafe.Pointer) int {
	log.Printf("[%s] >>>>>>>>>> UpdateContainer()...", __FILE__)

	var C_update_info C.container_update_s
	C_update_info.container_name = C.CString(container_update.Container_Name)
	C_update_info.image_name = C.CString(container_update.Image_Name)

	var C_update_res C.container_update_res_s

	user_data := unsafe.Pointer(C._C_SetCallbackStruct(unsafe.Pointer(&callback), userdata))

	var ret = C.dockzen_update_container(&C_update_info, &C_update_res, (C.container_update_cb)(unsafe.Pointer(C._C_CallbackContainerUpdate)), unsafe.Pointer(user_data))

	if int(ret) == dockzen_h.DOCKZEN_ERROR_NONE {
		update_res.Container_Name = C.GoString(C_update_res.container_name)
		update_res.Image_name_Prev = C.GoString(C_update_res.image_name_prev)
		update_res.Image_name_New = C.GoString(C_update_res.image_name_new)
		update_res.Status = C.GoString(C_update_res.status)
	}

	defer func(){
		if unsafe.Pointer(C_update_info.container_name) != nil {
			C.free(unsafe.Pointer(C_update_info.container_name))
			C_update_info.container_name = nil
		}
		if unsafe.Pointer(C_update_info.image_name) != nil {
			C.free(unsafe.Pointer(C_update_info.image_name))
			C_update_info.image_name = nil
		}
		if unsafe.Pointer(C_update_res.container_name) != nil {
			C.free(unsafe.Pointer(C_update_res.container_name))
			C_update_res.container_name = nil
		}
		if unsafe.Pointer(C_update_res.image_name_prev) != nil {
			C.free(unsafe.Pointer(C_update_res.image_name_prev))
			C_update_res.image_name_prev = nil
		}
		if unsafe.Pointer(C_update_res.image_name_new) != nil {
			C.free(unsafe.Pointer(C_update_res.image_name_new))
			C_update_res.image_name_new = nil
		}
		if unsafe.Pointer(C_update_res.status) != nil {
			C.free(unsafe.Pointer(C_update_res.status))
			C_update_res.status = nil
		}
	}()

	return int(ret)
}
