package common

type DeployInput struct {
	ModuleID      string
	ApplicationID string
	ProjectID     string
	UserID        string
	AppName       string
	Namespace     string
	Config        map[string]interface{}
}

type InstanceInfo struct {
	Name        string
	ReferenceID string
	Extra       []byte
}

type DeployOutput struct {
	Data []InstanceInfo
}
