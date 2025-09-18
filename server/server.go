package server

import (
	auth "AppPlaygroundService/authentication"
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/eventpublish"
	epGRPC "AppPlaygroundService/eventpublish/grpc"
	"AppPlaygroundService/logger"
	modApplication "AppPlaygroundService/modules/application"
	"AppPlaygroundService/modules/fsmhandler/handler/application"
	"AppPlaygroundService/modules/lbmevents"
	"AppPlaygroundService/modules/openstack"
	"AppPlaygroundService/modules/opskresource"
	"AppPlaygroundService/modules/opstkidentity"
	"AppPlaygroundService/services"
	"AppPlaygroundService/storages"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	midUnary "AppPlaygroundService/middlewares/unary"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	client "github.com/Zillaforge/appplaygroundserviceclient/aps"
	tkErr "github.com/Zillaforge/toolkits/errors"
	"github.com/Zillaforge/toolkits/mviper"
	"github.com/Zillaforge/toolkits/tracer"
)

var srvHTTP *http.Server
var srvGRPC *grpc.Server

func Run() {
	start()
	signAction()
	stop()
}

func start() {
	prepareUpstreamServices()
	startGRPCServer()
	startUpstreamServices()
	startHTTPServer()
}

func stop() {
	stopHTTPServer()
	stopGRPCServer()
	stopUpstreamServices()
}

func prepareUpstreamServices() {
	{ // logger
		// 初始化 Logger
		logger.Init("app_playground_service.log")
		// 初始化 Access Logger
		logger.InitAccessLogger("app_playground_service_access.log")
	}

	{ // tracer (jaeger)
		if mviper.GetBool("app_playground_service.tracer.enable") {
			tracer.Init(&tracer.Config{
				ServiceName: cnt.Kind,
				Endpoint:    mviper.GetString("app_playground_service.tracer.host"),
				Timeout:     mviper.GetInt("app_playground_service.tracer.timeout"),
			})
		}
	}

	{ // storages (database)
		storages.New(mviper.GetString("storage.provider"))

		// migrate database to migrate map latest version
		if mviper.GetBool("storage.auto_migrate") {
			if err := storages.Exec().AutoMigration(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func startUpstreamServices() {
	{ // services
		// 初始化 Services
		if err := services.InitServices(); err != nil {
			fmt.Println(tkErr.New(cnt.ServerInternalServerErr).WithInner(err))
			os.Exit(1)
		}
	}

	{ // authentication
		auth.Init(mviper.GetString("authentication.service"))
	}

	{ // opskresource
		opskresource.New(mviper.GetString("openstack_resource.service"))
	}

	{ // fsm
		application.Init()
	}

	{ // openstack
		if err := openstack.Init(mviper.Get("openstack")); err != nil {
			fmt.Println(tkErr.New(cnt.ServerInternalServerErr).WithInner(err))
			os.Exit(1)
		}
	}

	{ // openskidentity
		opstkidentity.Init()
	}

	{ // application
		modApplication.New(mviper.GetString("app_playground_service.application.provider"))
	}

	{
		// Initialize all of event publish plugins
		eventpublish.InitAllPlugins()
	}

	{
		// 啟動 LBM
		lbmevents.Init()
	}
}

func startGRPCServer() {
	zap.L().Info(serverStartInfo("GRPC", mviper.GetString("app_playground_service.grpc.host"), mviper.GetBool("app_playground_service.tls.enable")))

	lis, err := net.Listen("tcp", mviper.GetString("app_playground_service.grpc.host"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to listen of GRPC Port: %v, %s", err, mviper.GetString("app_playground_service.grpc.host")))
	}
	grpcOptions := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(tracer.RequestIDParser(), midUnary.ActionParser()),
		grpc.UnaryInterceptor(tracer.NewGRPCUnaryServerInterceptor()),
		grpc.WriteBufferSize(mviper.GetInt("app_playground_service.grpc.write_buffer_size")),
		grpc.ReadBufferSize(mviper.GetInt("app_playground_service.grpc.read_buffer_size")),
		grpc.MaxRecvMsgSize(mviper.GetInt("app_playground_service.grpc.max_receive_message_size")),
		grpc.MaxSendMsgSize(mviper.GetInt("app_playground_service.grpc.max_send_message_size")),
	}

	if mviper.GetBool("app_playground_service.tls.enable") {
		c, err := credentials.NewServerTLSFromFile(mviper.GetString("app_playground_service.tls.cert_path"), mviper.GetString("app_playground_service.tls.key_path"))
		if err != nil {
			log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
		}
		grpcOptions = append(grpcOptions, grpc.Creds(c))
	}

	srvGRPC = grpc.NewServer(grpcOptions...)
	grpcRouters(srvGRPC)
	go func() {
		if err := srvGRPC.Serve(lis); err != nil {
			zap.L().Error(fmt.Sprintf("failed to GRPC Server: %v", err))
		}
	}()
	os.Remove(mviper.GetString("app_playground_service.grpc.unix_socket.path"))
	lis2, err := net.Listen("unix", mviper.GetString("app_playground_service.grpc.unix_socket.path"))
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to listen of GRPC Port: %v, %s", err, mviper.GetString("app_playground_service.grpc.host")))
	}
	go func() {
		if err := srvGRPC.Serve(lis2); err != nil {
			zap.L().Error(fmt.Sprintf("failed to GRPC Server: %v", err))
		}
	}()
	client.Init(client.PoolProvider{
		Mode: client.UnixMode,
		UnixProvider: client.UnixProvider{
			SocketPath: mviper.GetString("app_playground_service.grpc.unix_socket.path"),
			ConnCount:  mviper.GetInt("app_playground_service.grpc.unix_socket.conn_count"),
		},
		WriteBufferSize:       mviper.GetInt("app_playground_service.grpc.write_buffer_size"),
		ReadBufferSize:        mviper.GetInt("app_playground_service.grpc.read_buffer_size"),
		MaxReceiveMessageSize: mviper.GetInt("app_playground_service.grpc.max_receive_message_size"),
		MaxSendMessageSize:    mviper.GetInt("app_playground_service.grpc.max_send_message_size"),
	})
}

func startHTTPServer() {
	zap.L().Info(serverStartInfo("HTTP",
		mviper.GetString("app_playground_service.http.host"),
		mviper.GetBool("app_playground_service.tls.enable")))

	if mviper.GetBool("app_playground_service.tls.enable") {
		srvHTTP = &http.Server{
			Addr:    mviper.GetString("app_playground_service.http.host"),
			Handler: router(),
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			},
		}
		go func() {
			if err := srvHTTP.ListenAndServeTLS(
				mviper.GetString("app_playground_service.tls.cert_path"),
				mviper.GetString("app_playground_service.tls.key_path")); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	} else {
		srvHTTP = &http.Server{
			Addr:    mviper.GetString("app_playground_service.http.host"),
			Handler: router(),
		}
		go func() {
			if err := srvHTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	}
}

func signAction() {
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
				stop()
				start()
			default:
				fmt.Println("Other Signal ", sign)
				return
			}
		}
	}
}

func stopHTTPServer() {
	srvHTTP.Shutdown(context.Background())
}
func stopGRPCServer() {
	srvGRPC.Stop()
}
func stopUpstreamServices() {
	epGRPC.ClosePlugins()
	tracer.Shutdown()
}

func serverStartInfo(serverName, host string, tlsEnable bool) string {
	tls := "disabled"
	if tlsEnable {
		tls = "enabled"
	}
	return fmt.Sprintf("%s server started at %s and TLS is %s", serverName, host, tls)
}
