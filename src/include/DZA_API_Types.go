package include

/**
 * @struct Container
 * @brief This structure contains container information.
 *
 * The containers struct encapsulate container information in the one data
*/
type Container struct{
	ID string
	Name string
	ImageName string
	Status string
}

/**
 * @struct Containers_info
 * @brief This structure contains container list information.
 *
 * The containers struct encapsulate count and container information in the one data
*/
type Containers_info struct {
	Count int
	Containerinfo []Container
}

/**
 * @struct ContainerUpdateInfo
 * @brief This structure contains container update information.
 *
 * The containers struct encapsulate update information in the one data
*/
type ContainerUpdateInfo struct {
	Image_Name     string `json:"ImageName"`
	Container_Name string `json:"Name"`
}

/**
 * @struct ContainerUpdateRes
 * @brief This structure contains response information for container update
 *
 * The containers struct encapsulate update information in the one data
*/
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

/**
 * @struct ContainerUpdate_cb_s
 * @brief This structure contains callback information for container update
 *
 * The containers struct encapsulate update callback information in the one data
*/
type Container_update_cb_s struct {
	Container_name	string
	Image_name	string
	Status		string
}

/**
 * @brief  This enum contains dockzen error information
 *
 * The dockzen_error_e indicates what error is happened
 *
 */
const (
	DOCKZEN_ERROR_NONE int = iota 
	DOCKZEN_ERROR_INVALID_PARAMETER
	DOCKZEN_ERROR_OUT_OF_MEMORY
	DOCKZEN_ERROR_PERMISSION_DENIED
	DOCKZEN_ERROR_NOT_SUPPORTED
)
