package common

import (
	"AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/utility"
	"context"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type (
	Pagination struct {
		Limit  int `json:"limit" form:"limit,default=100" binding:"max=100"`
		Offset int `json:"offset" form:"offset,default=0" binding:"min=0"`
		_      struct{}
	}
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}

func (data *User) Fill(ctx context.Context) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
		err       error
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"data":  &data,
			"error": &err,
		},
	)

	authUserInput := &authCom.GetUserInput{ID: data.ID, Cacheable: true}
	authUserOutput, err := authentication.Use().GetUser(ctx, authUserInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.GRPC, "authentication.Use().GetUser(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", authUserInput),
		).Warn(err.Error())
	} else {
		data.Name = authUserOutput.DisplayName
		data.Account = authUserOutput.Account
	}
}

type Resource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
