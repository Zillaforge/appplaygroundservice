package vps

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/opskresource/common"
	"AppPlaygroundService/utility"
	"context"
	"fmt"

	"go.uber.org/zap"
	"github.com/Zillaforge/toolkits/tracer"
	tkUtils "github.com/Zillaforge/toolkits/utilities"
	"github.com/Zillaforge/virtualplatformserviceclient/pb"
)

const (
	ActionAssociate    = "associate"
	ActionDisassociate = "disassociate"
)

func (h *Handler) UpdateFloatingIPStatus(ctx context.Context, input *common.UpdateFloatingIPStatusInput) (output *common.UpdateFloatingIPStatusOutput, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(tracer.Attributes{
		"input":  &input,
		"output": &output,
		"error":  &err,
	})

	switch input.Action {
	case ActionAssociate:
		associateFIPInput := &pb.FIPAssociateInput{
			Auth: &pb.AuthInfo{
				UserID:    input.IAMAuth.UserID,
				ProjectID: input.IAMAuth.ProjectID,
				Admin:     input.IAMAuth.IsAdmin,
			},
			ID: input.FloatingIPID,
			Device: &pb.FIPDeviceInput{
				Type:      "appaas", // 由 VPS 決定的名稱
				ID:        input.Device.ID,
				PortID:    input.Device.PortID,
				NetworkID: input.Device.NetworkID,
			},
		}
		if err = h.poolHandler.FloatingIP().Associate(associateFIPInput, ctx); err != nil {
			zap.L().With(
				zap.String(cnt.OpskResource, "h.poolHandler.FloatingIP().Associate(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", associateFIPInput),
			).Error(err.Error())
			return
		}
	case ActionDisassociate:
		disassociateFIPInput := &pb.IDInputWithAuth{
			Auth: &pb.AuthInfo{
				UserID:    input.IAMAuth.UserID,
				ProjectID: input.IAMAuth.ProjectID,
				Admin:     input.IAMAuth.IsAdmin,
			},
			ID: input.FloatingIPID,
		}
		if err = h.poolHandler.FloatingIP().Disassociate(disassociateFIPInput, ctx); err != nil {
			zap.L().With(
				zap.String(cnt.OpskResource, "h.poolHandler.FloatingIP().Disassociate(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", disassociateFIPInput),
			).Error(err.Error())
			return
		}
	default:
		err = fmt.Errorf("unsupported action: %s", input.Action)
		return
	}

	output = &common.UpdateFloatingIPStatusOutput{}
	return output, nil
}
