package webinterface

import (
  dockzen_h "include"
  )

type ContainerList_info struct {
	Cmd            string          			`json:"Cmd"`
	DeviceID       string          			`json:"DeviceID"`
	ContainerCount int             			`json:"ContainerCount"`
	Container      [10]dockzen_h.Container `json:"ContainerInfo"`
}

type ConnectReq struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}
