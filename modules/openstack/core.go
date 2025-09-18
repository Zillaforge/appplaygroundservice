package openstack

import (
	"AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	cnt "AppPlaygroundService/constants"
	opstkExec "AppPlaygroundService/modules/openstack/execution"
	"AppPlaygroundService/modules/openstack/keystone"
	"AppPlaygroundService/modules/openstack/neutron"
	"AppPlaygroundService/services"
	"AppPlaygroundService/utility"
	"fmt"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/flatten"
	"github.com/Zillaforge/toolkits/tracer"
)

const _openstackKind string = "openstack"

var provider map[string]Resource = map[string]Resource{}

type Resource interface {
	Admin() opstkExec.AdminResource
	Neutron(projectID string, userID string) *neutron.Neutron
	Keystone(projectID string, userID string) *keystone.Keystone
}

func Init(config interface{}) (err error) {
	for _, service := range cast.ToSlice(config) {
		if service == nil {
			continue
		}
		cfg := viper.New()
		cfg.MergeConfigMap(cast.ToStringMap(service))

		if !cfg.IsSet("namespace") || !cfg.IsSet("service") {
			return tkErr.New(cnt.OpenstackNamespaceAndServiceIsRequiredErr)
		}

		namespace := cfg.GetString("namespace")
		service := cfg.GetString("service")

		if _, exist := provider[namespace]; exist {
			return tkErr.New(cnt.OpenstackNamespaceIsRepeatedErr)
		}

		switch services.ServiceMap[service].Kind {
		case _openstackKind:
			if value, ok := services.ServiceMap[service].Conn.(*opstkExec.Connection); ok {
				p := opstkExec.Connection{
					IdentityEndpoint:  value.IdentityEndpoint,
					AdminUsername:     value.AdminUsername,
					AdminPassword:     value.AdminPassword,
					AdminProject:      value.AdminProject,
					DomainName:        value.DomainName,
					AllowReauth:       value.AllowReauth,
					PidSource:         value.PidSource,
					Pid:               GetOpstkPID(value.PidSource),
					UsernameSource:    value.UsernameSource,
					PasswordOTPSecret: value.PasswordOTPSecret,
					Username:          GetOpstkUserName(value.UsernameSource, value.AdminUsername),
					Password:          GetOpstkUserPassword(value.UsernameSource, value.PasswordOTPSecret, value.AdminUsername, value.AdminPassword),
				}
				p.SetNamespace(namespace)
				provider[namespace] = &p
			}
		default:
			return tkErr.New(cnt.OpenstackTypeIsNotSupportedErr)
		}
	}
	return nil
}

func Namespace(namespace string) Resource {
	return provider[namespace]
}

func NamespaceIsLegal(namespace string) bool {
	_, ok := provider[namespace]
	return ok
}

func ListNamespaces() []string {
	namespaces := []string{}
	for ns := range provider {
		namespaces = append(namespaces, ns)
	}
	return namespaces
}

// Replace replaces global provider by p
func Replace(namespace string, p opstkExec.Connection) {
	provider[namespace] = &p
}

func GetOpstkPID(pidSource string) func(string) string {
	ctx := tracer.StartEntryContext(tracer.EmptyRequestID)

	return func(projectID string) string {
		authProjectInput := &authCom.GetProjectInput{ID: projectID, Cacheable: true}
		authProjectOutput, _ := authentication.Use().GetProject(ctx, authProjectInput)
		projectInfo, _ := flatten.Flatten(authProjectOutput.ToMap(), "", flatten.DotStyle)
		return fmt.Sprintf("%v", projectInfo[pidSource])
	}
}

func GetOpstkUserName(usernameSource string, defaultName string) func(string) string {
	ctx := tracer.StartEntryContext(tracer.EmptyRequestID)

	return func(userID string) string {
		if userID == "" {
			return defaultName
		}
		authUserInput := &authCom.GetUserInput{ID: userID, Cacheable: true}
		authUserOutput, _ := authentication.Use().GetUser(ctx, authUserInput)
		userInfo, _ := flatten.Flatten(authUserOutput.ToMap(), "", flatten.DotStyle)
		return fmt.Sprintf("%v", userInfo[usernameSource])
	}
}

func GetOpstkUserPassword(usernameSource string, pwdOTPSecret string, defaultName string, defaultPassword string) func(string) string {
	return func(userID string) string {
		if userID == "" {
			return defaultPassword
		}
		username := GetOpstkUserName(usernameSource, defaultName)(userID)
		passCode, _ := utility.GetPassCode(username, pwdOTPSecret)
		return passCode
	}
}
