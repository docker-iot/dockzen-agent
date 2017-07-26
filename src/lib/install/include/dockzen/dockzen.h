/**
 * @file        dockzen.h
 * @brief       library for providing dockzen API

 * Copyright (c) 2017 Samsung Electronics Co., Ltd.
 * This software is the confidential and proprietary information
 * of Samsung Electronics, Inc. ("Confidential Information"). You
 * shall not disclose such Confidential Information and shall use
 * it only in accordance with the terms of the license agreement
 * you entered into with Samsung.
 */

#ifndef __DOCKZEN_H__
#define __DOCKZEN_H__

#ifdef __cplusplus
extern "C" {
#endif

#include "dockzen_types.h"

#ifndef API
#define API __attribute__ ((visibility("default")))
#endif

/**
 * @fn        int dockzen_get_containers_info(containers_info_s *c_info)
 * @brief     this function to call 'container info'
 * @param     *c_info		[out] fill containers infomation
 * @return   int 			return of function
 */
API int dockzen_get_containers_info(containers_info_s *c_info);
/**
 * @fn        int dockzen_update_container(container_update_s *container_update, container_update_res_s *container_update_res, container_update_cb callback, void* user_data)
 * @brief     this function to call 'container info'
 * @param     *container_update			[in] update parameters
 * @param     *container_update_res		[out] update response
 * @param     *callback					[inout] callback function to get update status
 * @param	  *user_data				[in] user token
 * @return   int 			return of function
 */
API int dockzen_update_container(container_update_s *container_update, container_update_res_s *container_update_res, container_update_cb callback, void* user_data);

#ifdef __cplusplus
}
#endif

#endif /* __DOCKZEN_H__ */
