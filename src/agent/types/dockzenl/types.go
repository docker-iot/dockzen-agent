package dockzenl

type DeviceState struct {
	CurrentState string `json:"CurrentState"`
}

type Container struct {
	ContainerID     string `json:"ContainerID"`
	ContainerName   string `json:"ContainerName"`
	ImageName       string `json:"ImageName"`
	ContainerStatus string `json:"ContainerStatus"`
}

type ErrorReturn struct {
	Message string `json:"Message"`
}

type Cmd struct {
	Cmd string `json:"Cmd"`
}
type UpdateParam struct {
	ImageName     string `json:"ImageName"`
	ContainerName string `json:"ContainerName"`
}

type UpdateImageParameters struct {
	Cmd   string      `json:"Cmd"`
	Param UpdateParam `json:"UpdateParam"`
}

type GetContainersInfoReturn struct {
	Containers []Container `json:"Containers"`
}

type UpdateImageReturn struct {
	State DeviceState `json:"DeviceState"`
}
