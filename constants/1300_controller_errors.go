package constants

import tkErr "pegasus-cloud.com/aes/toolkits/errors"

const (
	// 1300xxxx: Controller

	ControllerInternalServerErrCode    = 13010000
	ControllerInternalServerErrMsg     = "internal server error"
	ControllerWhereQueryInvalidErrCode = 13000001
	ControllerWhereQueryInvalidErrMsg  = "where query invalid"
)

var (
	// 1300xxxx: Controller

	// 13000000(internal server error)
	ControllerInternalServerErr = tkErr.Error(ControllerInternalServerErrCode, ControllerInternalServerErrMsg)
	// 13000001(where query invalid)
	ControllerWhereQueryInvalidErr = tkErr.Error(ControllerWhereQueryInvalidErrCode, ControllerWhereQueryInvalidErrMsg)
)
