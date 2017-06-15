package services

import (
	"fmt"
	dockzen_h "include"
	dockzen_api "lib"
)


func __updatecallback(in int) {
	fmt.Println("__updatecallback OK!!! in = ", in)
}

func DZA_Update_Do(updateinfo dockzen_h.ContainerUpdateInfo) int {
	fmt.Println("UpdateImageRequest")

	var ret = dockzen_api.UpdateContainer(updateinfo, __updatecallback)
	
	return ret
}
