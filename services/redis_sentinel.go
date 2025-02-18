package services

import (
	"reflect"

	"github.com/go-redis/redis/v8"
)

type RedisSentinelInput struct {
	Name                                               string
	Hosts, MasterGroupName, Password, SentinelPassword interface{}
}

func UnmarshalRedisSentinel(s map[string]interface{}) {
	InitRedisSentinel(&RedisSentinelInput{
		Name:             s["name"].(string),
		Hosts:            s["hosts"],
		MasterGroupName:  s["master_group_name"],
		Password:         s["password"],
		SentinelPassword: s["sentinel_password"],
	})
}

func InitRedisSentinel(input *RedisSentinelInput) {
	conn := newRedisSentinel(&newRedisSentinelInput{
		hosts:            input.Hosts,
		masterGroupName:  input.MasterGroupName,
		password:         input.Password,
		sentinelPassword: input.SentinelPassword,
	})
	ServiceMap[input.Name] = &Service{
		Kind: _redisSentinelKind,
		Conn: conn,
	}
}

func NewRedisSentinel(input *RedisSentinelInput) (conn *redis.Client) {
	return newRedisSentinel(&newRedisSentinelInput{
		hosts:            input.Hosts,
		masterGroupName:  input.MasterGroupName,
		password:         input.Password,
		sentinelPassword: input.SentinelPassword,
	})
}

type newRedisSentinelInput struct {
	hosts, masterGroupName, password, sentinelPassword interface{}
}

func newRedisSentinel(input *newRedisSentinelInput) (conn *redis.Client) {
	return redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName: func(input interface{}) (output string) {
			if input != nil &&
				reflect.TypeOf(input).Kind() == reflect.String {
				return input.(string)
			}
			return ""
		}(input.masterGroupName),
		SentinelAddrs: func(input interface{}) (output []string) {
			if input != nil {
				for _, host := range input.([]interface{}) {
					output = append(output, host.(string))
				}
				return output
			}
			return []string{}
		}(input.hosts),
		Password: func(input interface{}) (output string) {
			if input != nil &&
				reflect.TypeOf(input).Kind() == reflect.String {
				return input.(string)
			}
			return ""
		}(input.password),
		SentinelPassword: func(input interface{}) (output string) {
			if input != nil &&
				reflect.TypeOf(input).Kind() == reflect.String {
				return input.(string)
			}
			return ""
		}(input.sentinelPassword),
	})
}
