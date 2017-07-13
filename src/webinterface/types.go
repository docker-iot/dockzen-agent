// Package webinterface implements a client for the WebSocket protocol.
package webinterface

import (
  dockzen_h "include"
  )


// ConnectReq structure contains request information for web server.
type ConnectReq struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}

// Ws_ContainerList_info structure contains container list information.
type ws_ContainerList_info struct {
	Cmd            string          			`json:"Cmd"`
	DeviceID       string          			`json:"DeviceID"`
	ContainerCount int             			`json:"ContainerCount"`
	Container      []dockzen_h.Container `json:"ContainerInfo"`
}

// Ws_ContainerUpdateReturn structure contains container update return information.
type ws_ContainerUpdateReturn struct {
    Cmd           string         `json:"Cmd"`
    DeviceID      string         `json:"DeviceID"`
    UpdateState   string         `json:"UpdateState"`
}

var __FILE__ = "WEBINTERFACE"
