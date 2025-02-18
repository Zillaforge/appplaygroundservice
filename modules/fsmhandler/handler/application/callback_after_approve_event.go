package application

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/application"
	"AppPlaygroundService/modules/application/common"
	"AppPlaygroundService/modules/fsmhandler"
	fsmCom "AppPlaygroundService/modules/fsmhandler/common/application"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/storages/tables"
	"AppPlaygroundService/tasks"
	"AppPlaygroundService/utility"
	"AppPlaygroundService/utility/fsm"
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
)

func (h Handler) callBackAfterApproveEvent() func(context.Context, *fsm.Event) {
	return func(ctx context.Context, e *fsm.Event) {
		var (
			requestID = utility.MustGetContextRequestID(ctx)
		)
		zap.L().With(
			zap.String(cnt.FSM, "callBackAfterApproveEvent"),
			zap.String(cnt.RequestID, requestID),
		).Info("")

		for _, arg := range e.Args {
			// deploy application
			if val, ok := arg.(*common.DeployInput); ok {
				deployOutput, err := application.Use().Deploy(ctx, *val)
				if err != nil {
					if e, ok := tkErr.IsError(err); ok {
						switch e.Code() {
						}
					}
					zap.L().With(
						zap.String(cnt.FSM, "application.Use().Deploy(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.Any("input", val),
					).Error(err.Error())

					// get summary log
					getSummaryLogInput := common.GetSummaryLogInput{
						ApplicationID: val.ApplicationID,
						ProjectID:     val.ProjectID,
					}
					getSummaryLogOutput, getSummaryLogErr := application.Use().GetSummaryLog(ctx, getSummaryLogInput)
					if getSummaryLogErr != nil {
						zap.L().With(
							zap.String(cnt.FSM, "application.Use().GetSummaryLog(...)"),
							zap.String(cnt.RequestID, requestID),
							zap.Any("input", getSummaryLogInput),
						).Error(getSummaryLogErr.Error())
						// set default log if get summary failed
						getSummaryLogOutput.Log = common.DefaultErrorMsg
					}

					// save summary log to db
					mapExtra := map[string]string{
						"error": getSummaryLogOutput.Log,
					}
					byteExtra, err := json.Marshal(mapExtra)
					updateApplicationInput := &pb.UpdateApplicationInput{
						ID:    val.ApplicationID,
						Extra: byteExtra,
					}
					_, updateApplicationErr := aps.UpdateApplication(updateApplicationInput, ctx)
					if updateApplicationErr != nil {
						zap.L().With(
							zap.String(cnt.FSM, "aps.UpdateApplication(...)"),
							zap.String(cnt.RequestID, requestID),
							zap.Any("input", updateApplicationInput),
						).Error(updateApplicationErr.Error())
					}

					// switch application state to fail if error occurred.
					if err = fsmhandler.Application.Event(ctx, val.ApplicationID, fsmCom.FailEvent); err != nil {
						zap.L().With(
							zap.String(cnt.FSM, "fsmhandler.Application.Event(...)"),
							zap.String(cnt.RequestID, requestID),
							zap.String("id", val.ApplicationID),
							zap.String("event", fsmCom.FailEvent),
						).Error(err.Error())
						return
					}
				}

				for _, instanceInfo := range deployOutput.Data {
					createInstanceInput := &pb.InstanceInfo{
						Name:          instanceInfo.Name,
						ReferenceID:   instanceInfo.ReferenceID,
						ApplicationID: val.ApplicationID,
						ProjectID:     val.ProjectID,
						Extra:         instanceInfo.Extra,
					}
					_, err = aps.CreateInstance(createInstanceInput, ctx)
					if err != nil {
						zap.L().With(
							zap.String(cnt.FSM, "aps.CreateInstance(...)"),
							zap.String(cnt.RequestID, requestID),
							zap.Any("input", createInstanceInput),
						).Error(err.Error())
						return
					}
				}

				// update application shiftable to true if the are only 1 instance
				if len(deployOutput.Data) == 1 {
					shiftable := true
					updateApplicationInput := &pb.UpdateApplicationInput{
						ID:        val.ApplicationID,
						Shiftable: &shiftable,
					}
					_, err := aps.UpdateApplication(updateApplicationInput, ctx)
					if err != nil {
						zap.L().With(
							zap.String(cnt.FSM, "aps.UpdateApplication(...)"),
							zap.String(cnt.RequestID, requestID),
							zap.Any("input", updateApplicationInput),
						).Error(err.Error())
						return
					}
				}

				// update application state to ready when done.
				if err = fsmhandler.Application.Event(ctx, val.ApplicationID, fsmCom.FinishEvent); err != nil {
					zap.L().With(
						zap.String(cnt.FSM, "fsmhandler.Application.Event(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.String("id", val.ApplicationID),
						zap.String("event", fsmCom.FinishEvent),
					).Error(err.Error())
					return
				}

				// 利用 Application ID 取得 Application Info
				getApplicationInput := &storCom.GetApplicationInput{
					ID: val.ApplicationID,
				}
				getApplicationOutput, getApplicationErr := storages.Use().GetApplication(ctx, getApplicationInput)
				if getApplicationErr != nil {
					zap.L().With(
						zap.String(cnt.FSM, "storages.Use().GetApplication()"),
						zap.String(cnt.RequestID, requestID),
						zap.Any("input", getApplicationInput),
					).Error(getApplicationErr.Error())
					return
				}

				// 新增一筆 Metering Record
				createMeteringInput := &pb.MeteringInfo{
					ApplicationID: val.ApplicationID,
					Name:          getApplicationOutput.Application.Name,
					ProjectID:     getApplicationOutput.Application.ProjectID,
					Creator:       getApplicationOutput.Application.CreatorID,
					Instances: func(instanceInfo []common.InstanceInfo) (byteInstances []byte) {
						i := &tables.MeteringInstance{}
						var instances []tasks.InstanceInfo
						for _, instance := range instanceInfo {
							json.Unmarshal(instance.Extra, i)
							instances = append(instances, tasks.InstanceInfo{
								ID:       i.Instance.ID,
								FlavorID: i.Instance.FlavorID,
							})
						}
						byteInstances, _ = json.Marshal(instances)
						return byteInstances
					}(deployOutput.Data),
					CreatedAt: getApplicationOutput.Application.CreatedAt.Format(time.RFC3339),
				}
				if _, createMeteringErr := aps.CreateMetering(createMeteringInput, ctx); createMeteringErr != nil {
					zap.L().With(
						zap.String(cnt.FSM, "aps.CreateMetering(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.Any("input", createMeteringInput),
					).Error(createMeteringErr.Error())
					return
				}
			}
		}
	}
}
