package utility

import "reflect"

func Interface2map(input interface{}) (m map[string]interface{}) {
	if m == nil {
		m = make(map[string]interface{})
	}
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() == reflect.Map {
		for _, e := range val.MapKeys() {
			if k, ok := e.Interface().(string); ok {
				m[k] = val.MapIndex(e).Interface()
			}
		}
	}
	return
}
