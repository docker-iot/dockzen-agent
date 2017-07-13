package services

import (
	dockzen_api "lib"
	dockzen_h "include"
	"log"
)


// DZA_Mon_GetContainersInfo calls the dockzen_api.GetContainerListsInfo function.
// Param containersInfo is container information array.
// This function returns result of dockzen_api.GetContainerListsInfo.
func DZA_Mon_GetContainersInfo(containersInfo *dockzen_h.Containers_info) int {

	log.Printf("[%s] DZA_Mon_GetContainersInfo call !! ", __FILE__)

	var ret = dockzen_api.GetContainerListsInfo(containersInfo)

	return ret
}
