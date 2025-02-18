package services

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

type Service struct {
	Kind string
	Conn interface{}
}

const (
	_iamKind           string = "iam"
	_redisSentinelKind string = "redis_sentinel"
	_vpsKind           string = "vps"
	_openstackKind     string = "openstack"
)

var ServiceMap = make(map[string]*Service)

func InitServices() (err error) {
	zap.L().Info("initial upstream services")
	for _, service := range mviper.Get("services").([]interface{}) {
		if service == nil {
			continue
		}
		s := utility.Interface2map(service)
		for _, key := range []string{"name", "kind"} {
			if s[key] == nil {
				return tkErr.New(cnt.ServiceNameIsRequiredErr)
			}
			if key == "name" {
				if reflect.TypeOf(s[key]).Kind() != reflect.String {
					return tkErr.New(cnt.ServiceNameMustBeAStringErr)
				}
				if ServiceMap[s["name"].(string)] != nil {
					return fmt.Errorf(cnt.ServiceNameIsRepeatedErrMsg)
				}
			}
		}
		switch s["kind"] {
		case _iamKind:
			UnmarshalIAM(s)
		case _redisSentinelKind:
			UnmarshalRedisSentinel(s)
		case _vpsKind:
			UnmarshalVPS(s)
		case _openstackKind:
			UnmarshalOpenstack(s)
		}
	}
	return nil
}
