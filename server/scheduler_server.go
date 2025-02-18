package server

import (
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/logger"
	"AppPlaygroundService/modules/eventconsume"
	"AppPlaygroundService/modules/opskresource"
	"AppPlaygroundService/services"
	"AppPlaygroundService/tasks"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	aps "pegasus-cloud.com/aes/appplaygroundserviceclient/aps"
	"pegasus-cloud.com/aes/meteringtoolkits/metering"
	tkErr "pegasus-cloud.com/aes/toolkits/errors"
	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/tracer"
)

func RunScheduler() {
	startScheduler()
	signActionForScheduler()
	stopScheduler()
}

func startScheduler() {
	prepareUpstreamServicesForScheduler()
	startUpstreamServicesForScheduler()
	prepareSchedulerServer()
	startSchedulerServer()
}

func stopScheduler() {
	stopSchedulerServer()
	stopUpstreamServicesForScheduler()
}

func signActionForScheduler() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	//lint:ignore S1000 ...
	for {
		select {
		case sign := <-quit:
			switch sign {
			case syscall.SIGINT, syscall.SIGTERM:
				zap.L().Info("Shutdown Service.")
				return
			case syscall.SIGUSR2:
				stopScheduler()
				startScheduler()
			default:
				fmt.Println("Other Signal ", sign)
				return
			}
		}
	}
}

func prepareUpstreamServicesForScheduler() {
	{ // logger
		// 初始化 Logger
		logger.Init("app_playground_service_scheduler.log")
		// 初始化 Event Consume Logger
		logger.InitEventConsumeLogger("app_playground_service_eventconsume.log")
	}
	{ // tracer (jaeger)
		if mviper.GetBool("app_playground_service.tracer.enable") {
			tracer.Init(&tracer.Config{
				ServiceName: mviper.GetString("app_playground_service_scheduler.instance"),
				Endpoint:    mviper.GetString("app_playground_service.tracer.host"),
				Timeout:     mviper.GetInt("app_playground_service.tracer.timeout"),
			})
		}
	}
}

func startUpstreamServicesForScheduler() {}

func prepareSchedulerServer() {
	{ // 初始化 gRPC Connection
		aps.Init(aps.PoolProvider{
			Mode: aps.TCPMode,
			TCPProvider: aps.TCPProvider{
				Hosts: mviper.GetStringSlice("app_playground_service_scheduler.core_grpc.hosts"),
				TLS: aps.TLSConfig{
					Enable:   mviper.GetBool("app_playground_service_scheduler.core_grpc.tls.enable"),
					CertPath: mviper.GetString("app_playground_service_scheduler.core_grpc.tls.cert_path"),
				},
				ConnPerHost: mviper.GetInt("app_playground_service_scheduler.core_grpc.connection_per_host"),
			},
		})
	}

	{ // services
		// 初始化 Services
		if err := services.InitServices(); err != nil {
			fmt.Println(tkErr.New(cnt.ServerInternalServerErr).WithInner(err))
			os.Exit(1)
		}
	}

	{ // event consume
		eventconsume.New(mviper.GetString("event_consume.service"))
	}

	{ // metering init
		metering.Init(&metering.AMQP{
			Account:                mviper.GetString("metering_service.account"),
			Password:               mviper.GetString("metering_service.password"),
			Host:                   mviper.GetString("metering_service.host"),
			ManageHost:             mviper.GetString("metering_service.manage_host"),
			Timeout:                mviper.GetInt("metering_service.timeout"),
			RPCTimeout:             mviper.GetInt("metering_service.rpc_timeout"),
			Vhost:                  mviper.GetString("metering_service.vhost"),
			OperationConnectionNum: mviper.GetInt("metering_service.operation_connection_num"),
			ChannelNum:             mviper.GetInt("metering_service.channel_num"),
			ReplicaNum:             mviper.GetInt("metering_service.replica_num"),
			ConsumerConnectionNum:  mviper.GetInt("metering_service.consumer_connection_num"),
		})
	}

	{ // opskresource
		opskresource.New(mviper.GetString("openstack_resource.service"))
	}
}

func startSchedulerServer() {
	// 啟動定期同步 Resources
	tasks.InitSchedulerTasks()
}

func stopSchedulerServer() {
}

func stopUpstreamServicesForScheduler() {
	tracer.Shutdown()
}
