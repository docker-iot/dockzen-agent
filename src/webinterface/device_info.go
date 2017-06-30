package webinterface

import (
  "log"
  "net"
  "os"
)

func GetHardwareAddress() (string, error) {

	currentNetworkHardwareName := "eth0"
	netInterface, err := net.InterfaceByName(currentNetworkHardwareName)

	if err != nil {
		log.Printf("[%s] err = ",__FILE__, err)
	}

	name := netInterface.Name
	macAddress := netInterface.HardwareAddr

	log.Printf("[%s] Hardware name : ", __FILE__, string(name))

	hwAddr, err := net.ParseMAC(macAddress.String())

	if err != nil {
		log.Printf("[%s] No able to parse MAC address : ", __FILE__, err)
		os.Exit(-1)
	}

	log.Printf("[%s] Physical hardware address : ", __FILE__, hwAddr.String())

	return hwAddr.String(), nil
}
