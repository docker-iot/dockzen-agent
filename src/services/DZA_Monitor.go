package services

import (
	dockzen_api "lib"
	dockzen_h "include"
	"log"
)

func DZA_Mon_GetContainersInfo(containersInfo *dockzen_h.Containers_info) int {

	log.Printf("[%s] DZA_Mon_GetContainersInfo call !! ", __FILE__)

	var ret = dockzen_api.GetContainerListsInfo(containersInfo)

	return ret
}
