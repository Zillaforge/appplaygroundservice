package services

import (
	opstkExec "AppPlaygroundService/modules/openstack/execution"
	"reflect"
)

type OpenstackInput struct {
	Name string
	IdentityEndpoint, AdminUsername, AdminPassword, DomainName,
	AllowReauth, AdminProject, PidSource, UsernameSource, PasswordOTPSecret interface{}
}

func UnmarshalOpenstack(s map[string]interface{}) {
	if err := InitOpenstack(&OpenstackInput{
		Name:              s["name"].(string),
		IdentityEndpoint:  s["identity_endpoint"],
		AdminUsername:     s["admin_username"],
		AdminPassword:     s["admin_password"],
		DomainName:        s["domain_name"],
		AllowReauth:       s["allow_reauth"],
		AdminProject:      s["admin_project"],
		PidSource:         s["pid_source"],
		UsernameSource:    s["username_source"],
		PasswordOTPSecret: s["pwd_otp_secret"],
	}); err != nil {
		panic(err)
	}
}

func InitOpenstack(input *OpenstackInput) (err error) {
	_connection, err := opstkExec.New(&opstkExec.Connection{
		IdentityEndpoint:  checkStringValue(input.IdentityEndpoint),
		AdminUsername:     checkStringValue(input.AdminUsername),
		AdminPassword:     checkStringValue(input.AdminPassword),
		UsernameSource:    checkStringValue(input.UsernameSource),
		DomainName:        checkStringValue(input.DomainName),
		AllowReauth:       checkBoolValue(input.AllowReauth),
		AdminProject:      checkStringValue(input.AdminProject),
		PidSource:         checkStringValue(input.PidSource),
		PasswordOTPSecret: checkStringValue(input.PasswordOTPSecret),
	})
	if err != nil {
		return err
	}

	ServiceMap[input.Name] = &Service{
		Kind: _openstackKind,
		Conn: _connection,
	}
	return
}

func checkStringValue(input interface{}) (output string) {
	if input != nil &&
		reflect.TypeOf(input).Kind() == reflect.String {
		return input.(string)
	}
	return ""
}

func checkBoolValue(input interface{}) (output bool) {
	if input != nil &&
		reflect.TypeOf(input).Kind() == reflect.Bool {
		return input.(bool)
	}
	return false
}
