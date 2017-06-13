package lib

type ContainerListsInfo struct {
	ContainerID string `json:"ContainerID"`
	ContainerName string `json:"ContainerName"`
	ImageName string `json:"ImageName"`
	ContainerStatus string `json:"ContainerStatus"`
}
