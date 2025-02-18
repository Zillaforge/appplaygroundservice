package opstkidentity

import (
	"AppPlaygroundService/authentication"
	authCom "AppPlaygroundService/authentication/common"
	"context"
	"fmt"

	"pegasus-cloud.com/aes/toolkits/flatten"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (o *OpstkIdentity) GetOpstkUID(ctx context.Context, userID string) (opstkUID string, err error) {
	var (
		funcName = tkUtils.NameOfFunction().String()
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"projectID": &userID,
		"opstkUID":  &opstkUID,
		"error":     &err,
	})
	authUserInput := &authCom.GetUserInput{ID: userID, Cacheable: true}
	authUserOutput, err := authentication.Use().GetUser(ctx, authUserInput)
	if err != nil {
		return "", err
	}
	userInfo, err := flatten.Flatten(authUserOutput.ToMap(), "", flatten.DotStyle)
	if err != nil {
		return "", err

	}
	opstkUID = fmt.Sprintf("%v", userInfo[o.UidSource])
	return
}
