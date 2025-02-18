package vps

import (
	"AppPlaygroundService/services"

	"pegasus-cloud.com/aes/virtualplatformserviceclient/vps"
)

var (
	_serviceName string = ""
)

type Option func()

func WithServiceName(name string) Option {
	return func() {
		_serviceName = name
	}
}

// Input defines configuration which will be used by IAM solution
type Input struct {
	Hosts       []string
	TLSEnable   bool
	TLSCerPath  string
	ConnPerHost int
}

// Handler stores pool handler
type Handler struct {
	poolHandler *vps.PoolHandler
}

// New creates a pool handler and new connection with iam server in it
func New(opts ...Option) (handler *Handler) {
	handler = &Handler{}

	for _, opt := range opts {
		opt()
	}
	handler.poolHandler = services.ServiceMap[_serviceName].Conn.(*vps.PoolHandler)
	return
}
