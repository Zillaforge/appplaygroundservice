package grpc

import (
	"AppPlaygroundService/utility/querydecoder"
	"reflect"

	"google.golang.org/protobuf/types/known/emptypb"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

// A generic empty message that you can re-use to avoid defining duplicated
var EmptyPb = &emptypb.Empty{}

func WhereErrorParser(input error) error {
	return whereErrorParser(reflect.ValueOf(input))
}
func whereErrorParser(v reflect.Value) error {
	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			return whereErrorParser(reflect.ValueOf(v.MapIndex(key).Interface()))
		}
	case reflect.Struct:
		switch err := v.Interface().(type) {
		case querydecoder.UnknownKeyError:
			return tkErr.New(cCnt.GRPCWhereBindingErr).With("field", err.Key)
		case querydecoder.RegexError:
			return tkErr.New(cCnt.GRPCWhereBindingErr).WithInner(err)
		}
	}
	return nil
}
