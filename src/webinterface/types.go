package webinterface

import (
  dockzen_h "include"
  )

type ConnectReq struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}

type ws_ContainerList_info struct {
	Cmd            string          			`json:"Cmd"`
	DeviceID       string          			`json:"DeviceID"`
	ContainerCount int             			`json:"ContainerCount"`
	Container      []dockzen_h.Container `json:"ContainerInfo"`
}

type ws_ContainerUpdateReturn struct {
    Cmd           string         `json:"Cmd"`
    DeviceID      string         `json:"DeviceID"`
    UpdateState   string         `json:"UpdateState"`
}
