#ifndef __DOCKZEN_H__
#define __DOCKZEN_H__

#include "dockzen_types.h"

int dockzen_get_containers_info(containers_info_s *c_info);
int dockzen_update_container(container_update_s *container_update, container_update_res_s *container_update_res, container_update_cb callback, void* user_data);

#endif
