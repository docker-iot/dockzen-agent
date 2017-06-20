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
	Image_Name     string
	Container_Name string
}

type ContainerUpdateRes struct {
	Container_Name 	string
	Image_name_Prev	string
	Image_name_New 	string
	Status					string
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
