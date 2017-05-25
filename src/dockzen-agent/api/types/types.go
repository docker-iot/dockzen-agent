package types

import "golang.org/x/net/websocket"

const (
	ContainerServiceSocket string = "/var/run/container_service.sock"
)

type DockzenHndl struct {
	WS *websocket.Conn
}

type Cmd struct {
	Cmd string `json:"Cmd"`
}

type ContainerInfo struct {
	ImageName       string `json:"ImageName"`
	ContainerName   string `json:"ContainerName"`
	ContainerStatus string `json:"ContainerStatus"`
}
type ContainerLists struct {
	Cmd            string          `json:"Cmd"`
	DeviceID       string          `json:"DeviceID"`
	ContainerCount int             `json:"ContainerCount"`
	Container      []ContainerInfo `json:"ContainerInfo"`
}

type UpdateImageInfo struct {
	DeviceID  string `json:"DeviceID"`
	ImageName string `json:"ImageName"`
}

type UpdateImageParams struct {
	ImageName     string `json:"ImageName"`
	ContainerName string `json:"ContainerName"`
}

type UpdateImageReturn struct {
	Cmd         string `json:"Cmd"`
	DeviceID    string `json:"DeviceID"`
	UpdateState string `json:"UpdateState"`
}
