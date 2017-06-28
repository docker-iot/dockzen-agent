package lib

/*
#include <stdlib.h>
#include <stdio.h>
#include "dockzen.h"

typedef struct{
  void *callback;
  void *user_data;
}token_s;

void *_C_SetCallbackStruct(void * callback, void * user_data)
{
  token_s * token_ptr;
  token_ptr = (token_s*)malloc(sizeof(token_s));

  token_ptr->callback = callback;
  token_ptr->user_data = user_data;

  return (void*)token_ptr;
}

//////////////////////////// Private for each API
void _C_CallbackContainerUpdate(container_update_cb_s *status, void * user_data)
{
	void _GO_CallbackContainerUpdate(void *, void *);
	_GO_CallbackContainerUpdate((void*)status, user_data);
}
*/
import "C"
