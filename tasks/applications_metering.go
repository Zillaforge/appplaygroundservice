package tasks

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/opskresource"
	"AppPlaygroundService/modules/opskresource/common"
	"encoding/json"
	"time"

	"go.uber.org/zap"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkMetering "pegasus-cloud.com/aes/meteringtoolkits/metering"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

type InstanceInfo struct {
	ID       string                 `json:"id"`
	FlavorID string                 `json:"flavor_id"`
	Flavor   common.GetFlavorOutput `json:"flavor"`
	_        struct{}
}

type extraInstance struct {
	ID     string      `json:"id"`
	Flavor extraFlavor `json:"flavor"`
	_      struct{}
}

type extraFlavor struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	AZ     string   `json:"az"`
	Public bool     `json:"public"`
	Vcpu   int32    `json:"vcpu"`
	Memory int32    `json:"memory"`
	Disk   int32    `json:"disk"`
	GPU    extraGPU `json:"gpu"`
	_      struct{}
}

type extraGPU struct {
	Model  string `json:"model"`
	Count  int32  `json:"count"`
	IsVgpu bool   `json:"is_vgpu"`
	_      struct{}
}

func ApplicationsMetering() (err error) {
	var (
		funcName  = tkUtils.NameOfFunction().Name()
		requestID = tracer.EmptyRequestID
	)
	ctx, f := tracer.StartWithContext(tracer.StartEntryContext(requestID), funcName)
	defer f(tracer.Attributes{
		"err": &err,
	})

	// 列出所有 Metering Records
	listMeteringsInput := &pb.ListInput{
		Limit:  -1,
		Offset: 0,
	}
	listMeteringOutput, listMeteringErr := aps.ListMeterings(listMeteringsInput, ctx)
	if listMeteringErr != nil {
		zap.L().With(
			zap.String(cnt.Task, "aps.ListMeterings(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", listMeteringsInput),
		).Error(listMeteringErr.Error())
		return tkErr.New(cnt.StorageInternalServerErr)
	}

	for _, metering := range listMeteringOutput.Data {
		current := time.Now()
		rfc3339Current := current.Format(time.RFC3339)

		// 用 Application ID 取得 ProjectID 與 UserID
		getMeteringInput := &pb.GetInput{
			ID: metering.ApplicationID,
		}
		getMeteringOutput, getMeteringErr := aps.GetMetering(getMeteringInput, ctx)
		if getMeteringErr != nil {
			zap.L().With(
				zap.String(cnt.Task, "aps.GetMetering(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getMeteringInput),
			).Error(getMeteringErr.Error())
			continue
		}

		// 使用 Metering Toolkits 建立 Publish Struct
		publishInput := &tkMetering.PublishInput{
			Exchange:   mviper.GetString("app_playground_service_scheduler.tasks.applications_metering.metering_service.exchange"),
			RoutingKey: mviper.GetString("app_playground_service_scheduler.tasks.applications_metering.metering_service.routing_key"),
			Body: tkMetering.PublishBody{
				AvailabilityDistrict: mviper.GetString("app_playground_service.scopes.availability_district"),
				Service:              cnt.Name,
				ProjectID:            getMeteringOutput.ProjectID,
				UserID:               getMeteringOutput.Creator,
				ResourceType:         "application",
				ResourceID:           metering.ApplicationID,
				ResourceName:         getMeteringOutput.Name,
				RecordTime:           time.Now(),
				StartTime: func(input *pb.MeteringInfo) (output time.Time) {
					if input.LastPublishedAt == nil {
						parseOutput, parseErr := time.Parse(time.RFC3339, input.CreatedAt)
						if parseErr != nil {
							return time.Now()
						}
						return parseOutput
					}
					parseOutput, parseErr := time.Parse(time.RFC3339, *input.LastPublishedAt)
					if parseErr != nil {
						return time.Now()
					}
					return parseOutput
				}(metering),
				EndTime: current,
				Extra: map[string]interface{}{
					"instances": func(input []byte) (output interface{}) {
						instances := []InstanceInfo{}
						if unmarshalErr := json.Unmarshal(input, &instances); unmarshalErr != nil {
							zap.L().With(
								zap.String(cnt.Task, "json.Unmarshal(...)"),
								zap.String(cnt.RequestID, requestID),
								zap.String("data", string(input)),
							).Error(unmarshalErr.Error())
							return nil
						}
						extraInstances := []extraInstance{}
						for _, instance := range instances {
							getFlavorInput := &common.GetFlavorInput{
								ID: instance.FlavorID,
							}
							getFlavorOutput, getFlavorErr := opskresource.Use().GetFlavor(ctx, getFlavorInput)
							if getFlavorErr != nil {
								zap.L().With(
									zap.String(cnt.Task, "opskresource.Use().GetFlavor(...)"),
									zap.String(cnt.RequestID, requestID),
									zap.Any("input", getFlavorInput),
								).Error(getFlavorErr.Error())
								continue
							}

							extraInstances = append(extraInstances, extraInstance{
								ID: instance.ID,
								Flavor: extraFlavor{
									ID:     getFlavorOutput.ID,
									Name:   getFlavorOutput.Name,
									AZ:     getFlavorOutput.AZ,
									Public: getFlavorOutput.Public,
									Vcpu:   getFlavorOutput.Vcpu,
									Memory: getFlavorOutput.Memory,
									Disk:   getFlavorOutput.Disk,
									GPU: extraGPU{
										Model:  getFlavorOutput.Gpu.Model,
										Count:  getFlavorOutput.Gpu.Count,
										IsVgpu: getFlavorOutput.Gpu.IsVgpu,
									},
								},
							})
						}
						return extraInstances
					}(metering.Instances),
				},
			},
		}

		var last time.Time

		// 若 metering.EndedAt == nil 則表示 application 尚未被刪除
		if metering.EndedAt == nil {
			// 若 metering.LastPublishedAt == nil 則表示 application 剛被建立，尚未發送任何計量資訊
			if metering.LastPublishedAt == nil {
				// 取得 Application 建立時間
				createdAtParseOutput, parseErr := time.Parse(time.RFC3339, metering.CreatedAt)
				if parseErr != nil {
					zap.L().With(
						zap.String(cnt.Task, "time.Parse(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.String("value", metering.CreatedAt),
					).Error(parseErr.Error())
					continue
				}
				last = createdAtParseOutput
			} else {
				// 取得上次更新時間
				lastedPublishedAtOutput, lastedPublishedAtErr := time.Parse(time.RFC3339, *metering.LastPublishedAt)
				if lastedPublishedAtErr != nil {
					zap.L().With(
						zap.String(cnt.Task, "time.Parse(...)"),
						zap.String(cnt.RequestID, requestID),
						zap.String("value", *metering.LastPublishedAt),
					).Error(lastedPublishedAtErr.Error())
					continue
				}
				last = lastedPublishedAtOutput
			}
			// 更新 Application 的 LastPublishedAt
			updateMeteringInput := &pb.UpdateMeteringInput{
				ApplicationID:   metering.ApplicationID,
				LastPublishedAt: &rfc3339Current,
			}
			if _, updateMeteringErr := aps.UpdateMetering(updateMeteringInput, ctx); updateMeteringErr != nil {
				zap.L().With(
					zap.String(cnt.Task, "aps.UpdateMetering(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", updateMeteringInput),
				).Error(updateMeteringErr.Error())
				continue
			}
		} else {
			// 取得結束時間
			endedAtOutput, endedAtErr := time.Parse(time.RFC3339, *metering.EndedAt)
			if endedAtErr != nil {
				zap.L().With(
					zap.String(cnt.Task, "time.Parse(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.String("value", *metering.EndedAt),
				).Error(endedAtErr.Error())
				continue
			}
			last = endedAtOutput
			deleteMeteringInput := &pb.DeleteMeteringInput{
				ApplicationID: metering.ApplicationID,
			}
			if _, deleteMeteringErr := aps.DeleteMetering(deleteMeteringInput, ctx); deleteMeteringErr != nil {
				zap.L().With(
					zap.String(cnt.Task, "aps.DeleteMetering(...)"),
					zap.String(cnt.RequestID, requestID),
					zap.Any("input", deleteMeteringInput),
				).Error(deleteMeteringErr.Error())
				continue
			}
		}

		// 使用上次更新時間與當前時間計算差值
		duration := current.Sub(last)
		if duration < time.Second {
			publishInput.Body.Value = uint64(duration.Milliseconds())
			publishInput.Body.Unit = tkMetering.MillisecondUnit
		} else {
			publishInput.Body.Value = uint64(duration.Seconds())
			publishInput.Body.Unit = tkMetering.SecondUnit
		}

		// 使用 metering toolkits 呼叫 Publish function 發送計量資訊
		if publishErr := tkMetering.Publish(ctx, publishInput); publishErr != nil {
			zap.L().With(
				zap.String(cnt.Task, "tkMetering.Publish(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", publishInput),
			).Error(publishErr.Error())
			continue
		}
	}
	return nil
}
