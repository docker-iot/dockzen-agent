package services

import (
	"log"
	dockzen_h "include"
	dockzen_api "lib"
	"unsafe"
)


// Static __Updatecallback is callback function of updatecontainer command.
// param status is status information of updated container.
func __Updatecallback(status dockzen_h.Container_update_cb_s, user_data unsafe.Pointer) {
	log.Printf("[%s] service Callback OK!!!!", __FILE__)
	log.Printf("[%s] __Updatecallback> status.Container_name = ", __FILE__, status.Container_name)
	log.Printf("[%s] __Updatecallback> status.Image_name = ", __FILE__, status.Image_name)
	log.Printf("[%s] __Updatecallback> status.Status = ", __FILE__, status.Status)

	update_data := *(*update_userData)(user_data)
	log.Printf("[%s] Updatecallback> user_data.ContainerName = ", __FILE__, update_data)

}


// DZA_Update_Do calls the dockzen_api.UpdateContainer function.
// Param updateinfo is structure of container information.
// Param updateReturn is structure of update information.
 func DZA_Update_Do(updateinfo dockzen_h.ContainerUpdateInfo, updateReturn *dockzen_h.ContainerUpdateRes) int {
	log.Printf("[%s] UpdateImageRequest", __FILE__)

	var userdata update_userData
	userdata.Container_Name ="test"//= updateinfo.Container_Name

	log.Printf("[%s] userdata containerName =", __FILE__, userdata.Container_Name)

	var ret = dockzen_api.UpdateContainer(updateinfo, updateReturn, __Updatecallback, unsafe.Pointer(&userdata))

	log.Printf("[%s] updateReturn->status = ", __FILE__, updateReturn.Status)
	return ret
}
