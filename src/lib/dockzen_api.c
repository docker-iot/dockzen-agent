#include <stdio.h>
#include "dockzen_api_types.h"

ContainersInfo capi_Dockzen_GetContainerListsInfo(void)
{
	ContainersInfo stContainerListsInfo={2,{"container-id", "container1", "tizen1", "running", "container-id", "container2", "tizen2", "running"}};

	printf("~~~api_GetContainerListsInfo~~~");
	return stContainerListsInfo;
}
