package server

import (
	ctlAppCredential "AppPlaygroundService/controllers/grpc/app_credential"
	ctlApplication "AppPlaygroundService/controllers/grpc/application"
	ctlInstance "AppPlaygroundService/controllers/grpc/instance"
	ctlMetering "AppPlaygroundService/controllers/grpc/metering"
	ctlModule "AppPlaygroundService/controllers/grpc/module"
	ctlModuleAcl "AppPlaygroundService/controllers/grpc/module_acl"
	ctlModuleCategory "AppPlaygroundService/controllers/grpc/module_category"
	ctlModuleJoinModuleAcl "AppPlaygroundService/controllers/grpc/module_join_module_acl"
	ctlProject "AppPlaygroundService/controllers/grpc/project"

	"google.golang.org/grpc"
	"github.com/Zillaforge/appplaygroundserviceclient/pb"
)

func grpcRouters(srv *grpc.Server) {
	pb.RegisterProjectCRUDControllerServer(srv, new(ctlProject.Method))
	pb.RegisterModuleCategoryCRUDControllerServer(srv, new(ctlModuleCategory.Method))
	pb.RegisterModuleCRUDControllerServer(srv, new(ctlModule.Method))
	pb.RegisterModuleAclCRUDControllerServer(srv, new(ctlModuleAcl.Method))
	pb.RegisterModuleJoinModuleAclCRUDControllerServer(srv, new(ctlModuleJoinModuleAcl.Method))
	pb.RegisterApplicationCRUDControllerServer(srv, new(ctlApplication.Method))
	pb.RegisterInstanceCRUDControllerServer(srv, new(ctlInstance.Method))
	pb.RegisterMeteringCRUDControllerServer(srv, new(ctlMetering.Method))
	pb.RegisterAppCredentialCRUDControllerServer(srv, new(ctlAppCredential.Method))
}
