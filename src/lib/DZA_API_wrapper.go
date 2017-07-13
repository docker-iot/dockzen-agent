package lib

/*
#include <stdlib.h>
#include <stdio.h>
#include "dockzen.h"


// @struct token_s
// @brief This structure contains user data and callback information
//
// The containers struct encapsulate user data and callback in the one data
typedef struct{
  void *callback;
  void *user_data;
}token_s;

// @fn	_C_SetCallbackStruct(void * callback, void * user_data)
// @brief This function set callback structure
//
// @param	callback		[in] callback function
// @param user_data			[in] user data
// @return void *			[out] address for callback structure
void *_C_SetCallbackStruct(void * callback, void * user_data)
{
  token_s * token_ptr;
  token_ptr = (token_s*)malloc(sizeof(token_s));

  token_ptr->callback = callback;
  token_ptr->user_data = user_data;

  return (void*)token_ptr;
}

//////////////////////////// Private for each API

// @fn	_C_SetCallbackStruct(void * callback, void * user_data)
// @brief This function calls _GO_CallbackContainerUpdate
//
// @param	status		[in] status for update container
// @param user_data			[in] user data
void _C_CallbackContainerUpdate(container_update_cb_s *status, void * user_data)
{
	void _GO_CallbackContainerUpdate(void *, void *);
	_GO_CallbackContainerUpdate((void*)status, user_data);
}
*/
import "C"
