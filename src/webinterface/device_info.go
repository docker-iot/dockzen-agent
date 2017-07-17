package webinterface

import (
  "log"
  "net"
  "io/ioutil"
  "encoding/json"
  "os"
)

var UNIQUE_ID_FILE_PATH = "data"
var UNIQUE_ID_FILE  = "device_uuid.json"
const DEFAULT_UNIQUE_ID = "default device uuid"

// Static getHardwareAddress get the unique HW id to distinguish in web dash board.
// In temporarily, return mac address althouth it is not proper in bridge network mode.
func getHardwareAddress() (string, error) {

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


// set the unique id to distinguish in web dash board.
func setUniqueID(data string) int {
	var uniqueid string
	uniqueid = data

	if _, err := os.Stat(UNIQUE_ID_FILE_PATH); os.IsNotExist(err) {
		err = os.MkdirAll(UNIQUE_ID_FILE_PATH, 0755)
		if err != nil {
			log.Printf("[%s] setUniqueID %s folder create error!!", __FILE__, UNIQUE_ID_FILE_PATH)
			return -1
		}
	}

	f, err := os.Create(UNIQUE_ID_FILE_PATH + "/" + UNIQUE_ID_FILE)
	if err != nil {
		log.Printf("[%s] setUniqueID file create error!!!", __FILE__)
		return -1
	}

	uniqueid_json, err := json.Marshal(uniqueid)
	log.Printf("[%s] setUniqueID uniqueid_json=", __FILE__, string(uniqueid_json))

	_, err = f.Write(uniqueid_json)
	if err != nil{
		log.Printf("[%s] setUniqueID file write error", __FILE__)
		return -1
	}

	defer f.Close()

	return 0
}

// get the unique id from unique_id.json file.
// if this file is not exist, then create default one.
func getUniqueID() string {

	if _, err := os.Stat(UNIQUE_ID_FILE_PATH + "/" + UNIQUE_ID_FILE); os.IsNotExist(err) {
		// initial value for uniqueid is set hardware address
		log.Printf("[%s] setUniqueID!!!", __FILE__)
		var hw_addr string
		hw_addr, err := getHardwareAddress()
		if err != nil{
			log.Printf("[%s] HardwareAddress error = ", __FILE__, err)
			return DEFAULT_UNIQUE_ID
		}
		if setUniqueID(hw_addr) != 0 {
			return DEFAULT_UNIQUE_ID
		}
	}

	data, err := ioutil.ReadFile(UNIQUE_ID_FILE_PATH + "/" + UNIQUE_ID_FILE)

	if err != nil {
		log.Printf("[%s] getUniqueID file error!", __FILE__)
		return DEFAULT_UNIQUE_ID
	}
	var uniqueid string

	err = json.Unmarshal([]byte(data), &uniqueid)
	if err != nil {
		log.Printf("[%s] getUniqueID data error!!!!", __FILE__)
		return DEFAULT_UNIQUE_ID
	}

	log.Printf("[%s] getUniqueID = %s", __FILE__, uniqueid)

	return uniqueid
}
