package lib

type ContainerInfo struct {
	ContainerID string `json:"ContainerID"`
	ContainerName string `json:"ContainerName"`
	ImageName string `json:"ImageName"`
	ContainerStatus string `json:"ContainerStatus"`
}

type ContainerLists struct {
	Cmd            string          			`json:"Cmd"`
	DeviceID       string          			`json:"DeviceID"`
	ContainerCount int             			`json:"ContainerCount"`
	Container      []ContainerInfo `json:"ContainerInfo"`
}
