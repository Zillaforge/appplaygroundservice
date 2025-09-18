package services

import (
	"AppPlaygroundService/utility"
	"reflect"

	vpsUtil "github.com/Zillaforge/virtualplatformserviceclient/utility"
	"github.com/Zillaforge/virtualplatformserviceclient/vps"
)

type VPSInput struct {
	Name                    string
	Hosts, TLS, ConnPerHost interface{}
}

func UnmarshalVPS(s map[string]interface{}) {
	if err := InitVPS(&VPSInput{
		Name:        s["name"].(string),
		Hosts:       s["hosts"],
		TLS:         s["tls"],
		ConnPerHost: s["connection_per_host"],
	}); err != nil {
		panic(err)
	}
}

func InitVPS(input *VPSInput) (err error) {
	const _defaultConn int = 3

	var _connection *vps.PoolHandler
	_connection, err = vps.New(vps.PoolProvider{
		Mode: vps.TCPMode,
		TCPProvider: vps.TCPProvider{
			Hosts: func(in interface{}) (output []string) {
				if in != nil {
					for _, host := range in.([]interface{}) {
						if reflect.TypeOf(host).Kind() == reflect.String {
							output = append(output, host.(string))
						}
					}
					return output
				}
				return []string{}
			}(input.Hosts),
			TLS: vps.TLSConfig{
				Enable: func(in map[string]interface{}) (output bool) {
					if in != nil &&
						in["enable"] != nil &&
						reflect.TypeOf(in["enable"]).Kind() != reflect.Bool {
						return in["enable"].(bool)
					}
					return false
				}(utility.Interface2map(input.TLS)),
				CertPath: func(in map[string]interface{}) (output string) {
					if in != nil &&
						in["cert_path"] != nil {
						return in["cert_path"].(string)
					}
					return ""
				}(utility.Interface2map(input.TLS)),
			},
			ConnPerHost: func(in interface{}) (output int) {
				if in != nil &&
					reflect.TypeOf(in).Kind() != reflect.Int {
					return in.(int)
				}
				return _defaultConn
			}(input.ConnPerHost),
		},
		RouteResponseType: vpsUtil.JSON,
	})
	if err != nil {
		return err
	}
	ServiceMap[input.Name] = &Service{
		Kind: _vpsKind,
		Conn: _connection,
	}
	return nil
}
