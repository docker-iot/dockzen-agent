package webinterface

import (
  dockzen_h "include"
  )

/**
 * @struct ConnectReq
 * @brief This structure contains request information for web server.
 *
 * The containers struct encapsulate connection information in the one data
*/
type ConnectReq struct {
	Cmd  string `json:"cmd"`
	Name string `json:"name"`
}

/**
 * @struct ws_ContainerList_info
 * @brief This structure contains container list information
 *
 * The containers struct encapsulate deviceid, count, command name and container information in the one data
*/
type ws_ContainerList_info struct {
	Cmd            string          			`json:"Cmd"`
	DeviceID       string          			`json:"DeviceID"`
	ContainerCount int             			`json:"ContainerCount"`
	Container      []dockzen_h.Container `json:"ContainerInfo"`
}

/**
 * @struct ws_ContainerUpdateReturn
 * @brief This structure contains container update return information.
 *
 * The containers struct encapsulate deviceid, command name and update information in the one data
*/
type ws_ContainerUpdateReturn struct {
    Cmd           string         `json:"Cmd"`
    DeviceID      string         `json:"DeviceID"`
    UpdateState   string         `json:"UpdateState"`
}

var __FILE__ = "WEBINTERFACE"
