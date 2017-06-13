typedef struct {
	char *ContainerID;
	char *ContainerName;
	char *ImageName;
	char *ContainerStatus;
}ContainerListsInfo;

ContainerListsInfo capi_Dockzen_GetContainerListsInfo(void);
