package common

type GetFlavorInput struct {
	ID string
}

type GetFlavorOutput struct {
	ID     string
	Name   string
	AZ     string
	Public bool
	Vcpu   int32 `json:"vcpu"`
	Memory int32
	Disk   int32
	Gpu    GpuInfo
}

type GpuInfo struct {
	Model  string
	Count  int32
	IsVgpu bool `json:"is_vgpu"`
}
