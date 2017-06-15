#include <stdio.h>
#include "dockzen.h"
#include "dockzen_types.h"

/**
 *  capi_GetContainersInfo interface
 */
API int dockzen_get_containers_info(containers_info_s *containers_info)
{
	containers_info->count = 2;

	containers_info->container[0].id = "container-id1";
	containers_info->container[0].name = "container-name1";
	containers_info->container[0].image_name = "image-name1";
	containers_info->container[0].status = "exited";

	containers_info->container[1].id = "container-id2";
	containers_info->container[1].name = "container-name2";
	containers_info->container[1].image_name = "image-name2";
	containers_info->container[1].status = "running";

	return DOCKZEN_ERROR_NONE;
}

API int dockzen_update_container(container_update_s *container_update, container_update_cb callback, void* user_data)
{

	//callback("running", user_data);
	callback(12345678, user_data);

	return DOCKZEN_ERROR_NONE;
}
