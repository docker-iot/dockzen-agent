#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "../src/lib/dockzen.h"
#include "../src/lib/dockzen_types.h"

/**
 *  capi_GetContainersInfo interface
 */
API int dockzen_get_containers_info(containers_info_s *containers_info)
{
	containers_info->count = 1;

	containers_info->container[0].id = (char *)malloc(strlen("container-id1") +1);
	memset(containers_info->container[0].id, 0x00, strlen("container-id1")+1);
	sprintf(containers_info->container[0].id, "container-id1");

	containers_info->container[0].name = (char *)malloc(strlen("container-name1") +1);
	memset(containers_info->container[0].name, 0x00, strlen("container-name1")+1);
	sprintf(containers_info->container[0].name, "container-name1");

	containers_info->container[0].image_name = (char *)malloc(strlen("image-name1") +1);
	memset(containers_info->container[0].image_name, 0x00, strlen("image-name1")+1);
	sprintf(containers_info->container[0].image_name, "image-name1");

	containers_info->container[0].status = (char *)malloc(strlen("exited") +1);
	memset(containers_info->container[0].status, 0x00, strlen("exited")+1);
	sprintf(containers_info->container[0].status, "exited");


	return DOCKZEN_ERROR_NONE;
}

API int dockzen_update_container(container_update_s *container_update, container_update_res_s *container_update_return, container_update_cb callback, void* user_data)
{
	container_update_cb_s * updateinfo;

	updateinfo = (container_update_cb_s *)malloc(sizeof(container_update_cb_s));


	updateinfo->container_name = (char *)malloc(strlen(container_update->container_name) +1);

	memset(updateinfo->container_name, 0x00, strlen(container_update->container_name)+1);
	sprintf(updateinfo->container_name, "%s", container_update->container_name);

	updateinfo->image_name = (char *)malloc(strlen(container_update->image_name) +1);
	memset(updateinfo->image_name, 0x00, strlen(container_update->image_name)+1);
	sprintf(updateinfo->image_name, "%s", container_update->image_name);

	updateinfo->status = (char *)malloc(strlen("completed") +1);
	memset(updateinfo->status, 0x00, strlen("completed")+1);
	sprintf(updateinfo->status, "completed");

	////////////////////
	container_update_return->container_name = (char *)malloc(strlen(container_update->container_name) +1);
	memset(container_update_return->container_name, 0x00, strlen(container_update->container_name)+1);
	sprintf(container_update_return->container_name, "%s", container_update->container_name);

	container_update_return->image_name_prev = (char *)malloc(strlen("tizen2") +1);
	memset(container_update_return->image_name_prev, 0x00, strlen("tizen2")+1);
	sprintf(container_update_return->image_name_prev, "tizen2");

	container_update_return->image_name_new = (char *)malloc(strlen(container_update->image_name) +1);
	memset(container_update_return->image_name_new, 0x00, strlen(container_update->image_name)+1);
	sprintf(container_update_return->image_name_new, "%s", container_update->image_name);

	container_update_return->status = (char *)malloc(strlen("running") +1);
	memset(container_update_return->status, 0x00, strlen("running")+1);
	sprintf(container_update_return->status, "running");

	callback(updateinfo, user_data);

	return DOCKZEN_ERROR_NONE;
}
