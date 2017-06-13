#include <stdio.h>
#include "dockzen_api_types.h"

ContainerListsInfo capi_Dockzen_GetContainerListsInfo(void)
{
	ContainerListsInfo stContainerListsInfo = {"container-id", "container-name", "image-name", "111"};

	printf("~~~api_GetContainerListsInfo~~~");
	return stContainerListsInfo;
}
