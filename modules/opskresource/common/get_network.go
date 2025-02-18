package common

type GetNetworkInput struct {
	ID string
}

type GetNetworkOutput struct {
	ID        string
	Name      string
	ProjectID string
	Namespace string
	RouterID  string
	SubnetID  string
}
