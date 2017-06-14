package lib

import (
	"fmt"
	"log"
	"net"
	"os"
)

/*
#include <dockzen_api_types.h>
*/
import "C"

func GetHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	log.Printf("Hardware name : %s\n", string(name))

	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		log.Printf("No able to parse MAC address : %s\n", err)
		os.Exit(-1)
	}

	log.Printf("Physical hardware address : %s \n", hwAddr.String())

	return hwAddr.String(), nil
}

func GetContainerListsInfo() (ContainerLists, error ) {

	fmt.Println(">>>>>>>>>> GetContainerListsInfo()...")

	lists := C.capi_Dockzen_GetContainerListsInfo()

	var info ContainerLists

	info.Cmd = "GetContainersInfo"

	numOfList := int(lists.Count)
	info.ContainerCount = numOfList

	macaddress, err := GetHardwareAddress()

	if err == nil {
		fmt.Println("Mac address error!!!!")
	}

	info.DeviceID = macaddress
	fmt.Println("DevicedID = ", info.DeviceID)
	var containerValue ContainerInfo

	for i := 0; i < numOfList; i++ {
		containerValue = ContainerInfo{
			ContainerID: 			C.GoString(lists.Container[i].ID),
			ContainerName: 		C.GoString(lists.Container[i].Name),
			ImageName: 				C.GoString(lists.Container[i].ImageName),
			ContainerStatus: 	C.GoString(lists.Container[i].Status),
		}
		info.Container = append(info.Container, containerValue)
	}

	fmt.Println("ContainerInfo -> ", info)

	return info, nil

}
