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

func (o *OpstkIdentity) GetOpstkPID(ctx context.Context, projectID string) (opstkPID string, err error) {
	var (
		funcName = tkUtils.NameOfFunction().String()
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"projectID": &projectID,
		"opstkPid":  &opstkPID,
		"error":     &err,
	})
	authProjectInput := &authCom.GetProjectInput{ID: projectID, Cacheable: true}
	authProjectOutput, err := authentication.Use().GetProject(ctx, authProjectInput)
	if err != nil {
		return "", err
	}
	projectInfo, err := flatten.Flatten(authProjectOutput.ToMap(), "", flatten.DotStyle)
	if err != nil {
		return "", err

	}
	opstkPID = fmt.Sprintf("%v", projectInfo[o.PidSource])
	return
}
