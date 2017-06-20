package services

import (
	"fmt"
	dockzen_h "include"
	dockzen_api "lib"
	"unsafe"
)

func __updatecallback(status dockzen_h.Container_update_cb_s, user_data unsafe.Pointer) {
	fmt.Println("service Callback OK!!!!")
	fmt.Println("__updatecallback> status.Container_name = ", status.Container_name)
	fmt.Println("__updatecallback> status.Image_name = ", status.Image_name)
	fmt.Println("__updatecallback> status.Status = ", status.Status)

	update_data := *(*update_userData)(user_data)
	fmt.Println("__updatecallback> user_data.ContainerName = ", update_data)

}

func DZA_Update_Do(updateinfo dockzen_h.ContainerUpdateInfo, updateReturn *dockzen_h.ContainerUpdateRes) int {
	fmt.Println("UpdateImageRequest")

	var userdata update_userData
	userdata.Container_Name = updateinfo.Container_Name

	fmt.Println("DZA_Update_Do> userdata containerName =", userdata.Container_Name)

	var ret = dockzen_api.UpdateContainer(updateinfo, updateReturn, __updatecallback, unsafe.Pointer(&userdata))

	fmt.Println("DZA_Update_Do> updateReturn->status = ", updateReturn.Status)
	return ret
}
