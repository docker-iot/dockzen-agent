/**
 * @file        dockzen_types.h
 * @brief       Types of API for dockzen

 * Copyright (c) 2017 Samsung Electronics Co., Ltd.
 * This software is the confidential and proprietary information
 * of Samsung Electronics, Inc. ("Confidential Information"). You
 * shall not disclose such Confidential Information and shall use
 * it only in accordance with the terms of the license agreement
 * you entered into with Samsung.
 */

#ifndef __DOCKZEN_TYPES_H__
#define __DOCKZEN_TYPES_H__

#define MAX_CONTAINER_NUM	(10)

/**
 * @brief  This enum contains dockzen error information
 *
 * The dockzen_error_e indicates what error is happened
 *
 */
typedef enum {
	DOCKZEN_ERROR_NONE = 0,				/**< Successful */
	DOCKZEN_ERROR_INVALID_PARAMETER,	/**< Invalid parameter */
	DOCKZEN_ERROR_OUT_OF_MEMORY,		/**< Out of memory */
	DOCKZEN_ERROR_PERMISSION_DENIED,	/**< Permission denied */
	DOCKZEN_ERROR_NOT_SUPPORTED,		/**< Not supported  */
} dockzen_error_e;

/**
 * @brief  This enum contains dockzen update state and error information
 *
 * The dockzen_update_e indicates what is dockzen update state.
 *
 */
typedef enum {
	DOCKZEN_UPDATE_STATE_STARTED = 0,
	DOCKZEN_UPDATE_STATE_DOWNLOADING,
	DOCKZEN_UPDATE_STATE_UPDATING,
	DOCKZEN_UPDATE_STATE_DONE,
	DOCKZEN_UPDATE_ERROR_UNKNOWN,
	DOCKZEN_UPDATE_ERROR_NOTSTART,
	DOCKZEN_UPDATE_ERROR_DOWNLOAD,
	DOCKZEN_UPDATE_ERROR_UPDATE
} dockzen_update_e;


/**
 * @struct containers_info_s
 * @brief  This struct contains containers information
 *
 * The containers_info_s struct encapsulate count and container information in the one data
 *
 */
typedef struct{
	int count;		/**< the counts of containers info */
	struct {
		char * id;			/**< Numeric ID which is managed by swam */
		char * name;		/**< Container name from service */
		char * image_name;	/**< Docker image name */
		char * status;		/**< Container Status :  */
	}container[MAX_CONTAINER_NUM];			/**< Max Count constraint */
}containers_info_s;

/**
 * @struct container_update_s
 * @brief  This struct contains update parameters
 *
 * The container_update_s struct encapsulate container_name and image_name information in the one data
 *
 */
typedef struct{
	char * container_name;
	char * image_name;
}container_update_s;

/**
 * @struct container_update_res_s
 * @brief  This struct contains update response
 *
 * The container_update_res_s struct encapsulate container_name, image_name_pre, image_name_new and status in the one data
 *
 */
typedef struct{
	char * container_name;
	char * image_name_prev;
	char * image_name_new;
	char * status;
}container_update_res_s;

/**
 * @struct container_update_cb_s
 * @brief  This struct contains update callback parameters
 *
 * The container_update_cb_s struct encapsulate conainer_name, image_name and status in the one data
 *
 */
typedef struct{
	char * container_name;
	char * image_name;
	char * status;
}container_update_cb_s;

typedef void (*container_update_cb) (container_update_cb_s *status, void * user_data);

#endif /* __DOCKZEN_TYPES_H__ */
