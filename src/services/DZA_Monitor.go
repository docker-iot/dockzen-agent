package services

import (
	dockzen_api "lib"
	dockzen_h "include"
	"log"
)

/**
 * @fn	DZA_Mon_GetContainersInfo
 * @brief This function calls the dockzen_api.GetContainerListsInfo function.
 *
 * @param	containersInfo, [inout] containers information array.
 * @return int,						[out] dockzen_api.GetContainerListsInfo return value
*/
func DZA_Mon_GetContainersInfo(containersInfo *dockzen_h.Containers_info) int {

	log.Printf("[%s] DZA_Mon_GetContainersInfo call !! ", __FILE__)

	var ret = dockzen_api.GetContainerListsInfo(containersInfo)

	return ret
}
