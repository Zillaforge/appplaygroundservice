package instance

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/modules/openstack"
	openstackCom "AppPlaygroundService/modules/openstack/common"
	"AppPlaygroundService/modules/opskresource"
	opskCom "AppPlaygroundService/modules/opskresource/common"
	"AppPlaygroundService/modules/opskresource/vps"
	"AppPlaygroundService/storages"
	storCom "AppPlaygroundService/storages/common"
	"AppPlaygroundService/utility"
	"context"
	"encoding/json"

	"go.uber.org/zap"
	cCnt "pegasus-cloud.com/aes/appplaygroundserviceclient/constants"
	"pegasus-cloud.com/aes/appplaygroundserviceclient/pb"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/tracer"
	tkUtils "pegasus-cloud.com/aes/toolkits/utilities"
)

func (m *Method) AssociateFloatingIP(ctx context.Context, input *pb.UpdateFIPInput) (output *pb.InstanceDetail, err error) {
	var (
		funcName  = tkUtils.NameOfFunction().String()
		requestID = utility.MustGetContextRequestID(ctx)
	)

	ctx, f := tracer.StartWithContext(ctx, funcName)
	defer f(
		tracer.Attributes{
			"input":  &input,
			"output": &output,
			"error":  &err,
		},
	)

	// get instance
	getInput := &storCom.GetInstanceInput{
		ID: input.ID,
	}
	getOutput, err := storages.Use().GetInstance(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageInstanceNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCInstanceNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetInstance()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// check whether there is a FIP already attached on this instance
	if getOutput.Instance.FloatingIPID != "" {
		err = tkErr.New(cCnt.GRPCInstanceAlreadyHasFIPErr)
		zap.L().With(
			zap.String(cnt.GRPC, "getOutput.Instance.FloatingIPID != \"\""),
			zap.String(cnt.RequestID, requestID),
			zap.Any("FloatingIPID", getOutput.Instance.FloatingIPID),
		).Error(err.Error())
		return
	}

	// get floating ip
	opskGetFipInput := &opskCom.GetFloatingIPInput{
		ID: input.FloatingIPID,
	}
	opskGetFipOutput, err := opskresource.Use().GetFloatingIP(ctx, opskGetFipInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.OpskResourceRecordNotFoundErrCode:
				err = tkErr.New(cCnt.GRPCFloatingIPNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "opskresource.Use().GetFloatingIP(...)"),
			zap.Any("input", opskGetFipInput),
		).Error(err.Error())
		return
	}

	application := getOutput.Instance.Application
	// check resource permission
	if opskGetFipOutput.ProjectID != application.ProjectID {
		err = tkErr.New(cCnt.GRPCFloatingIPNotFoundErr)
		return
	}

	// check whether FIP can be use
	if opskGetFipOutput.Status != "DOWN" {
		err = tkErr.New(cCnt.GRPCFIPCannotBeUseErr)
		zap.L().With(
			zap.String(cnt.GRPC, "opskGetFipOutput.Status != \"DOWN\""),
			zap.String(cnt.RequestID, requestID),
			zap.Any("fipOutput", opskGetFipOutput),
		).Error(err.Error())
		return
	}

	extraMap := map[string]interface{}{}
	if getOutput.Instance.Extra != nil {
		if err = json.Unmarshal(getOutput.Instance.Extra, &extraMap); err != nil {
			zap.L().With(
				zap.String(cnt.GRPC, "json.Unmarshal(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", getOutput.Instance.Extra),
			).Error(err.Error())
			return
		}
	}

	// get portID from Extra
	portID, ok := extraMap["instance"].(map[string]interface{})["port_id"].(string)
	if !ok {
		err = tkErr.New(cCnt.GRPCPortIDNotFoundErr)
		zap.L().With(
			zap.String(cnt.GRPC, "extraMap[\"instance\"].(string); !ok"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("extra", extraMap),
		).Error(err.Error())
		return
	}

	// call openstack module to associate floating ip
	associateFloatingIpInput := &openstackCom.AssociateFloatingIpInput{
		FloatingIpID: opskGetFipOutput.UUID,
		PortID:       portID,
	}
	associateFloatingIpErr := openstack.Namespace(application.Namespace).Neutron(application.ProjectID, "").AssociateFloatingIp(ctx, associateFloatingIpInput)
	if associateFloatingIpErr != nil {
		zap.L().With(
			zap.String(cnt.Controller, "openstack.Namespace().Neutron().AssociateFloatingIp(...)"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", associateFloatingIpInput),
		).Error(associateFloatingIpErr.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(associateFloatingIpErr)
		return
	}

	// register associated floating IP to VPS
	updateFIPInput := &opskCom.UpdateFloatingIPStatusInput{
		Action:       vps.ActionAssociate,
		FloatingIPID: input.FloatingIPID,
		IAMAuth: opskCom.IAMAuthInfo{
			UserID:    application.CreatorID,
			ProjectID: application.ProjectID,
		},
		Device: opskCom.FIPDeviceInput{
			ID: application.ID,
		},
	}
	_, err = opskresource.Use().UpdateFloatingIPStatus(ctx, updateFIPInput)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Controller, "opskresource.Use().UpdateFloatingIPStatus(...)"),
			zap.Any("input", updateFIPInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)

		// disassociated the FloatingIP from openstack if error
		disassociateFloatingIpInput := &openstackCom.DisassociateFloatingIpInput{
			FloatingIpID: opskGetFipOutput.UUID,
		}
		disassociateFloatingIpErr := openstack.Namespace(application.Namespace).Neutron(application.ProjectID, "").DisassociateFloatingIp(ctx, disassociateFloatingIpInput)
		if disassociateFloatingIpErr != nil {
			zap.L().With(
				zap.String(cnt.Controller, "openstack.Namespace().Neutron().DisassociateFloatingIp(...)"),
				zap.String(cnt.RequestID, requestID),
				zap.Any("input", disassociateFloatingIpInput),
			).Error(disassociateFloatingIpErr.Error())
			return
		}
		return
	}

	// update instance floating ip
	updateInput := &storCom.UpdateInstanceInput{
		ID: input.ID,
		UpdateData: &storCom.InstanceUpdateInfo{
			FloatingIPID:      &input.FloatingIPID,
			FloatingIPAddress: &opskGetFipOutput.Address,
		},
	}
	_, err = storages.Use().UpdateInstance(ctx, updateInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageInstanceNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCInstanceNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().UpdateInstance()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", updateInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	// get latest instance detail
	getOutput, err = storages.Use().GetInstance(ctx, getInput)
	if err != nil {
		if e, ok := tkErr.IsError(err); ok {
			switch e.Code() {
			case cnt.StorageInstanceNotFoundErr.Code():
				err = tkErr.New(cCnt.GRPCInstanceNotFoundErr)
				return
			}
		}
		zap.L().With(
			zap.String(cnt.GRPC, "storages.Use().GetInstance()"),
			zap.String(cnt.RequestID, requestID),
			zap.Any("input", getInput),
		).Error(err.Error())
		err = tkErr.New(cCnt.GRPCInternalServerErr).WithInner(err)
		return
	}

	output = m.storage2pb(&getOutput.Instance)
	return
}
