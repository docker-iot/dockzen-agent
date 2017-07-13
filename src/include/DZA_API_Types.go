// Package include is types of API for dockzen.
package include

// Container structure contains container information.
type Container struct{
	ID string
	Name string
	ImageName string
	Status string
}

// Containers_info structure contains container list information.
type Containers_info struct {
	Count int
	Containerinfo []Container
}

// ContainerUpdateInfo structure contains container update information.
type ContainerUpdateInfo struct {
	Image_Name     string `json:"ImageName"`
	Container_Name string `json:"Name"`
}

// ContainerUpdateRes structure contains response information for container update.
type ContainerUpdateRes struct {
	Container_Name 	string
	Image_name_Prev	string
	Image_name_New 	string
	Status		string
}

// ContainerUpdate_cb_s contains callback information for container update.
type Container_update_cb_s struct {
	Container_name	string
	Image_name	string
	Status		string
}

// This enum contains dockzen error information.
//The dockzen_error_e indicates what error is happened.
const (
	DOCKZEN_ERROR_NONE int = iota 
	DOCKZEN_ERROR_INVALID_PARAMETER
	DOCKZEN_ERROR_OUT_OF_MEMORY
	DOCKZEN_ERROR_PERMISSION_DENIED
	DOCKZEN_ERROR_NOT_SUPPORTED
)
