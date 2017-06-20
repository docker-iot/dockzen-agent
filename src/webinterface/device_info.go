package webinterface

import (
  "fmt"
  "net"
  "os"
)

func GetHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		fmt.Println(err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	fmt.Println("Hardware name : ", string(name))

	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		fmt.Println("No able to parse MAC address : ", err)
		os.Exit(-1)
	}

	fmt.Println("Physical hardware address : ", hwAddr.String())

	return hwAddr.String(), nil
}
