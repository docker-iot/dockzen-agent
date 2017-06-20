package services

import (
	dockzen_api "lib"
	dockzen_h "include"
	"fmt"
)

func DZA_Mon_GetContainersInfo(containersInfo *dockzen_h.Containers_info) int {

	fmt.Println("DZA_Mon_GetContainersInfo call !! ")

	var ret = dockzen_api.GetContainerListsInfo(containersInfo)

	return ret
}
