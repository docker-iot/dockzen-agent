#ifndef __DOCKZEN_TYPES_H__
#define __DOCKZEN_TYPES_H__

#define MAX_CONTAINER_NUM	(10)

typedef enum {
	DOCKZEN_ERROR_NONE = 0,				/**< Successful */
	DOCKZEN_ERROR_INVALID_PARAMETER,	/**< Invalid parameter */
	DOCKZEN_ERROR_OUT_OF_MEMORY,		/**< Out of memory */
	DOCKZEN_ERROR_PERMISSION_DENIED,	/**< Permission denied */
	DOCKZEN_ERROR_NOT_SUPPORTED,		/**< Not supported  */
} dockzen_error_e;

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
 *  ContainersInfo type definition.
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
 *  Containers update info type definition.
 */
typedef struct{
	char * container_name;
	char * image_name;
}container_update_s;

/**
 *  Containers update response type definition.
 */
typedef struct{
	char * container_name;
	char * image_name_prev;
	char * image_name_new;
	char * status;
}container_update_res_s;

/**
 *  typedef container_update_cb
 */
typedef struct{
	char * container_name;
	char * image_name;
	char * status;
}container_update_cb_s;

typedef void (*container_update_cb) (container_update_cb_s *status, void * user_data);

#endif
