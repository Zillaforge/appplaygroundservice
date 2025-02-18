package common

type GetFloatingIPInput struct {
	ID string
}

type GetFloatingIPOutput struct {
	ID         string
	UUID       string
	Name       string
	ProjectID  string
	UserID     string
	Namespace  string
	Status     string
	Reserved   bool
	DeviceType string
	DeviceID   string
	Address    string
}
