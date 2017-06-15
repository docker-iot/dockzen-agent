package include

type Container struct{
	ID string `json:"ContainerID"`
	Name string `json:"ContainerName"`
	ImageName string `json:"ImageName"`
	Status string `json:"ContainerStatus"`
}

type Containers_info struct {
	Count int `json:"ContainerCount"`
	Containerinfo [10]Container `json:"ContainerInfo"`
}

type ContainerUpdateInfo struct {
	ImageName     string `json:"ImageName"`
	ContainerName string `json:"ContainerName"`
}
