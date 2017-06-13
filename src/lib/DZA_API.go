package lib

import "fmt"

/*
#cgo CFLAGS: -I../lib/
#cgo LDFLAGS: ../lib/libapi.a
#include <dockzen_api_types.h>
*/

import "C"

func GetContainerListsInfo() {

	fmt.Println(">>>>>>>>>> GetContainerListsInfo()...")

	value := C.capi_Dockzen_GetContainerListsInfo()

/*
    var info ContainerListsInfo

    info.ContainerID = C.GoString(value.ContainerID)
    info.ContainerName = C.GoString(value.ContainerName)
    info.ImageName = C.GoString(value.ImageName)
    info.ContainerStatus = C.GoString(value.ContainerStatus)
*/
}
