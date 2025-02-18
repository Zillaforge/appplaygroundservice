package common

type UpdateFloatingIPStatusInput struct {
	Action       string
	FloatingIPID string
	IAMAuth      IAMAuthInfo
	Device       FIPDeviceInput
}

type FIPDeviceInput struct {
	Type      string
	ID        string
	PortID    *string
	NetworkID *string
}

type UpdateFloatingIPStatusOutput struct {
}
