typedef struct {
	int Count;
	struct {
		char *ID;
		char *Name;
		char *ImageName;
		char *Status;
	}Container[10];
}ContainersInfo;

ContainersInfo capi_Dockzen_GetContainerListsInfo(void);
