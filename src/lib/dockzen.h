#ifndef __DOCKZEN_H
#define __DOCKZEN_H

typedef enum {
	DOCKZEN_ERROR_NONE = 0, /**< Successful */
	DOCKZEN_ERROR_INVALID_PARAMETER, /**< Invalid parameter */
	DOCKZEN_ERROR_OUT_OF_MEMORY, /**< Out of memory */
	DOCKZEN_ERROR_PERMISSION_DENIED, /**< Permission denied */
	DOCKZEN_ERROR_NOT_SUPPORTED /**< Not supported  */
} dockzen_error_e;


/**
 *  ContainersInfo type definition.
 */
typedef struct{
	int count;		/**< the counts of containers info */
	struct {
		char * id;			/**< Numeric ID which is managed by swam */
		char * name;		/**< Container name from service */
		char * image_name;	/**< Docker image name */
		char * status;		/**< Container Status :  */
	}container[10];			/**< Max Count constraint */
}containers_info_s;


/**
 *  ContainersInfo type definition.
 */
typedef struct{
	char * id;
}container_update_s;


//typedef void (*container_update_cb) (char* status, void* user_data);
typedef void (*container_update_cb) (int status, int user_data);


/**
 *  capi_GetContainersInfo interface.
 *	return created & serviced containers info which is contolled by dockzen-launcher
 */
int dockzen_get_containers_info(containers_info_s *c_info);


int dockzen_update_container(container_update_s *container_update, container_update_cb callback, void* user_data);

#endif
