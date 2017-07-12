package lib

import (
	"log"
	dockzen_h "include"
	"unsafe"
	"testing"
)

/*
#cgo CFLAGS: -I${SRCDIR}/install/include/dockzen/
#cgo LDFLAGS: -ldockzen -ljson-c

#include <stdlib.h>
#include "dockzen.h"

typedef struct{
	void *callback;
	void* user_data;
}token_s;

void* _C_SetCallbackStruct( void* callback, void* user_data);
void _C_CallbackContainerUpdate(int status);

/////////////////////// unit test API
int test_dockzen_get_containers_info(containers_info_s *containers_info);
int test_dockzen_update_container(container_update_s *container_update, container_update_res_s *container_update_return, container_update_cb callback, void* user_data);
*/
import "C"

type ContainerUpdateCB func(dockzen_h.Container_update_cb_s, unsafe.Pointer)
var __FILE__ = "LIB"

/**
 * @fn	[Mandatory] getContainerListsInfo_Res(C_containers_info C.containers_info_s,
  																						containers_info *dockzen_h.Containers_info)
 * @brief [Mandatory] This function convert C structure data to go structure data
 *
 * @param	C_containers_info,		[in] container information structure of C structure.
 * @param containers_info,			[inout] container information structure
 * @return void
*/
func getContainerListsInfo_Res(C_containers_info C.containers_info_s, containers_info *dockzen_h.Containers_info){
	log.Printf("[%s] >>> API GetContainerListsInfo Request", __FILE__)

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

/**
 * @fn	[Mandatory] GetContainerListsInfo(containers_info *dockzen_h.Containers_info) int
 * @brief [Mandatory] This function calls the dockzen library function.
 *
 * @param	containers_info,		[inout] container information structure
 * @return int,								[out] dockzen library function return value.
*/
func GetContainerListsInfo(containers_info *dockzen_h.Containers_info) int {

	log.Printf("[%s] >>>>>>>>>> API GetContainerListsInfo()...", __FILE__)

	var C_containers_info C.containers_info_s
	var ret = C.dockzen_get_containers_info(&C_containers_info)
	if int(ret) == dockzen_h.DOCKZEN_ERROR_NONE {
		getContainerListsInfo_Res(C_containers_info, containers_info)
	}

	log.Printf("[%s] container = ", __FILE__, containers_info)
	return int(ret)
}


/**
 * @fn	[Mandatory] _GO_CallbackContainerUpdate(c_status_info unsafe.Pointer,
 																								userdata unsafe.Pointer)
 * @brief [Mandatory] This function is callback function for updatecontainer command
 *
 * @param	c_status_info,		[in] container status information structure
 * @param userdata,					[in] user data
 * @return void
*/
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

/**
 * @fn	[Mandatory] updateContainer_Res(C_update_res C.container_update_res_s,
  																			update_res * dockzen_h.ContainerUpdateRes)
 * @brief [Mandatory] This function convert C structure data to go structure data.
 *
 * @param	C_update_res,		[in] C update information struture
 * @param	update_res,			[inout] update information struture
 * @return void
*/
func updateContainer_Res(C_update_res C.container_update_res_s, update_res * dockzen_h.ContainerUpdateRes){

	update_res.Container_Name = C.GoString(C_update_res.container_name)
	update_res.Image_name_Prev = C.GoString(C_update_res.image_name_prev)
	update_res.Image_name_New = C.GoString(C_update_res.image_name_new)
	update_res.Status = C.GoString(C_update_res.status)

	defer func(){

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
}

/**
 * @fn	[Mandatory] UpdateContainer(container_update dockzen_h.ContainerUpdateInfo,
 																		update_res * dockzen_h.ContainerUpdateRes,
																		callback ContainerUpdateCB,
																		userdata unsafe.Pointer) int
 * @brief [Mandatory] This function calls the updatecontainer function in dockzen library
 *
 * @param	c_status_info,		[in] container information structure
 * @return int,							[out] dockzen library function return value.
*/
func UpdateContainer(container_update dockzen_h.ContainerUpdateInfo, update_res * dockzen_h.ContainerUpdateRes, callback ContainerUpdateCB, userdata unsafe.Pointer) int {
	log.Printf("[%s] >>>>>>>>>> UpdateContainer()...", __FILE__)

	var C_update_info C.container_update_s
	C_update_info.container_name = C.CString(container_update.Container_Name)
	C_update_info.image_name = C.CString(container_update.Image_Name)

	var C_update_res C.container_update_res_s

	user_data := unsafe.Pointer(C._C_SetCallbackStruct(unsafe.Pointer(&callback), userdata))

	var ret = C.dockzen_update_container(&C_update_info, &C_update_res, (C.container_update_cb)(unsafe.Pointer(C._C_CallbackContainerUpdate)), unsafe.Pointer(user_data))
	if int(ret) == dockzen_h.DOCKZEN_ERROR_NONE {
		updateContainer_Res(C_update_res, update_res)
	}

	return int(ret)
}

/**
 * @fn	[Mandatory] testGetContainerListsInfo(t *testing.T)
 * @brief [Mandatory] This function is unit test for GetContainerListsInfo
 *
 * @param	t,		[in] testing struture
 * @return void
*/
func testGetContainerListsInfo(t *testing.T){
	var C_containers_info C.containers_info_s
	var containers_info dockzen_h.Containers_info

	var ret = C.test_dockzen_get_containers_info(&C_containers_info)

	if int(ret) == dockzen_h.DOCKZEN_ERROR_NONE{
		getContainerListsInfo_Res(C_containers_info, &containers_info)
		log.Printf("[TEST] container = ", containers_info)
	} else {
		t.Errorf("[TEST] GetContainerInfo error")
	}
}

/**
 * @fn	[Mandatory] testUpdateContainer(t *testing.T,
  																			callback ContainerUpdateCB)
 * @brief [Mandatory] This function is unit test for UpdatateContainer
 *
 * @param	t,					[in] testing struture
 * @param	callback,		[in] callback function
 * @return void
*/
func testUpdateContainer(t *testing.T, callback ContainerUpdateCB){
	var C_update_info C.container_update_s
	var update_res  dockzen_h.ContainerUpdateRes
	C_update_info.container_name = C.CString("tizen")
	C_update_info.image_name = C.CString("tizen_headless:v0.1")

	var C_update_res C.container_update_res_s
	user_data := unsafe.Pointer(C._C_SetCallbackStruct(unsafe.Pointer(&callback), nil))

	var ret = C.test_dockzen_update_container(&C_update_info, &C_update_res,(C.container_update_cb)(unsafe.Pointer(C._C_CallbackContainerUpdate)), unsafe.Pointer(user_data))

	if int(ret) == dockzen_h.DOCKZEN_ERROR_NONE {
		updateContainer_Res(C_update_res, &update_res)
		log.Printf("[TEST] update_res = ", update_res)
	} else {
		t.Errorf("[TEST] updateContainer error")
	}

}
