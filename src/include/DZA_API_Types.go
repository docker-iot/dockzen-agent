package include

type Container struct{
	ID string
	Name string
	ImageName string
	Status string
}

type Containers_info struct {
	Count int
	Containerinfo []Container
}

type ContainerUpdateInfo struct {
	Image_Name     string `json:"ImageName"`
	Container_Name string `json:"Name"`
}

type ContainerUpdateRes struct {
	Container_Name 	string
	Image_name_Prev	string
	Image_name_New 	string
	Status		string
}

type ContainerUpdate_cb struct {
	Container_Name string
	Status string
}

type Container_update_cb_s struct {
	Container_name	string
	Image_name	string
	Status		string
}

const (
	DOCKZEN_ERROR_NONE int = iota 
	DOCKZEN_ERROR_INVALID_PARAMETER
	DOCKZEN_ERROR_OUT_OF_MEMORY
	DOCKZEN_ERROR_PERMISSION_DENIED
	DOCKZEN_ERROR_NOT_SUPPORTED
)
