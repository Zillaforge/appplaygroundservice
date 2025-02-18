package services

import (
	"AppPlaygroundService/utility"
	"reflect"

	"pegasus-cloud.com/aes/pegasusiamclient/iam"
	iamUtil "pegasus-cloud.com/aes/pegasusiamclient/utility"
)

type IAMInput struct {
	Name                    string
	Hosts, TLS, ConnPerHost interface{}
}

func UnmarshalIAM(s map[string]interface{}) {
	if err := InitIAM(&IAMInput{
		Name:        s["name"].(string),
		Hosts:       s["hosts"],
		TLS:         s["tls"],
		ConnPerHost: s["connection_per_host"],
	}); err != nil {
		panic(err)
	}
}

func InitIAM(input *IAMInput) (err error) {
	const _defaultConn int = 3
	var _connection *iam.PoolHandler
	_connection, err = iam.New(iam.PoolProvider{
		Mode: iam.TCPMode,
		TCPProvider: iam.TCPProvider{
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
			TLS: iam.TLSConfig{
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
		RouteResponseType: iamUtil.JSON,
	})
	if err != nil {
		return err
	}
	ServiceMap[input.Name] = &Service{
		Kind: _iamKind,
		Conn: _connection,
	}
	return nil
}
