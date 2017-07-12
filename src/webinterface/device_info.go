package webinterface

import (
  "log"
  "net"
)

/**
 * @fn	GetHardwareAddress() (string, error)
 * @brief This function get the unique HW id to distinguish in web dash board.
 *   In temporarily, return mac address althouth it is not proper in bridge network mode.
 *
 * @return string,		[out] hardware address
 * @return error,     [out] error value (if the value is null, it is not an error.)
*/
func GetHardwareAddress() (string, error) {

	var netInterface net.Interface

	// get all Interfaces
	netInterfaceLists, err := net.Interfaces()
	
	if err != nil {
		// if failed to get Interfaces
		log.Printf("[%s] err = ",__FILE__, err)
		return "Error-Hardware-Address", nil
	} else {
		// search valid Interface in lists. validation means not-empty Name and Mac
		for _,v := range netInterfaceLists {
			if len(v.Name) > 0 && len(v.HardwareAddr) > 0 {
				netInterface = v;
				break;
			}
		}
		
		if netInterface.Index ==0 {
			// if failed to search valid Interfaces
			log.Printf("[%s] err = ",__FILE__, "can't find valid interface")
			for i, v := range netInterfaceLists {
				log.Printf("[%v] = [%v]", i, v)	
			}
			return "Invalid-Hardware-Address", nil
		} else {
			log.Printf("[%s] Interface.Name         : %v", __FILE__, netInterface.Name)
			log.Printf("[%s] Interface.HardwareAddr : %v", __FILE__, netInterface.HardwareAddr.String())

			return netInterface.HardwareAddr.String(), nil
		}
	}
}
