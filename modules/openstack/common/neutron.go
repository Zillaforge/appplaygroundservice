package common

type (
	AssociateFloatingIpInput struct {
		FloatingIpID       string
		PortID             string //from instance's extra
	}   
)

type (
	DisassociateFloatingIpInput struct {
		FloatingIpID       string
	}   
)



