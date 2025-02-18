package configs

import (
	cnt "AppPlaygroundService/constants"
	"math"

	tkCfg "pegasus-cloud.com/aes/toolkits/configs"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

func init() {
	mviper.SetDefault("version", cnt.Version, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("kind", cnt.PascalCaseName, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("location_id", getLocation(), "-", tkCfg.TypeSystem, tkCfg.RegionLocal)
	mviper.SetDefault("host_id", getHostname(), "-", tkCfg.TypeSystem, tkCfg.RegionLocal)

	mviper.SetDefault("app_playground_service.developer", true, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service.instance", cnt.PascalCaseName, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// http
	mviper.SetDefault("app_playground_service.http.host", "0.0.0.0:8111", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.http.access_control.allow_origins", []string{"*"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.http.access_control.allow_credentials", true, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.http.access_control.allow_headers", []string{"X-Language", "Origin", "Content-Length", "Content-Type", "Authorization"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.http.access_control.allow_methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.http.access_control.expose_headers", []string{"Content-Language", "host-id", "version-id", "location-id"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// grpc
	mviper.SetDefault("app_playground_service.grpc.host", "0.0.0.0:5111", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service.grpc.unix_socket.path", "/run/AppPlaygroundService.sock", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.unix_socket.conn_count", 20, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.write_buffer_size", 32*1024, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.read_buffer_size", 32*1024, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.max_receive_message_size", 1024*1024*4, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.max_send_message_size", math.MaxInt32, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.hosts", []string{"127.0.0.1:5106"}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.grpc.conn_per_host", 3, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// tls
	mviper.SetDefault("app_playground_service.tls.enable", false, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.tls.cert_path", "/var/lib/ASUS/hcs/tls/cert.pem", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.tls.key_path", "/var/lib/ASUS/hcs/tls/cert-key.pem", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// tracer
	mviper.SetDefault("app_playground_service.tracer.enable", false, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.tracer.host", "http://127.0.0.1:14268/api/traces", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.tracer.timeout", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// scopes
	mviper.SetDefault("app_playground_service.scopes.memcache_ttl", 60, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service.scopes.allow_namespaces", []string{}, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service.scopes.default_language", "zh-TW", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service.scopes.enable_resource_review", false, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service.scopes.availability_district", "default", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// log
	// system-log
	mviper.SetDefault("app_playground_service.logger.system_log.path", "/var/log/ASUS/", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.system_log.max_size", 100, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.system_log.max_backups", 5, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.system_log.max_age", 10, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.system_log.compress", false, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.system_log.mode", "error", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.system_log.show_in_console", true, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	// access-log
	mviper.SetDefault("app_playground_service.logger.access_log.path", "/var/log/ASUS/", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.access_log.max_size", 100, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.access_log.max_backups", 5, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.access_log.max_age", 10, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.access_log.compress", false, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.access_log.default_source_ip", "127.0.0.1", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	// event-log
	mviper.SetDefault("app_playground_service.logger.event_consume_log.path", "/var/log/ASUS/", "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.event_consume_log.max_size", 100, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.event_consume_log.max_backups", 5, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.event_consume_log.max_age", 10, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.logger.event_consume_log.compress", false, "-", tkCfg.TypeStale, tkCfg.RegionGlobal)

	// data_path
	mviper.SetDefault("app_playground_service.data_path.data_dir", "/var/lib/ASUS/appplaygroundservice/data", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.data_path.module_pid", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// application
	mviper.SetDefault("app_playground_service.application.provider", "terraform", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.application.summary_log_length", -1, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.application.terraform.bin_path", "/usr/local/bin/terraform", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// opstkidentity
	mviper.SetDefault("app_playground_service.opstkidentity.pid_source", "id", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service.opstkidentity.uid_source", "id", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// storage
	mviper.SetDefault("storage.provider", "mariadb", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.auto_migrate", true, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.host", "mariadb-galera.pegasus-system:3306", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.account", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.password", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.name", "pt", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.timeout", 5, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.max_open_conns", 150, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.conn_max_lifetime", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("storage.mariadb.max_idle_conns", 150, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// authentication
	mviper.SetDefault("authentication.service", "", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// openstack_resource
	mviper.SetDefault("openstack_resource.service", "", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// services
	mviper.SetDefault("services", []interface{}{}, "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)

	// plugins
	mviper.SetDefault("event_publish.plugins", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// openstack
	mviper.SetDefault("openstack", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// metering service
	mviper.SetDefault("metering_service.account", "account", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.password", "password", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.host", "127.0.0.1:5672", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.manage_host", "127.0.0.1:15672", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.timeout", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.rpc_timeout", 10, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.vhost", "/", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.operation_connection_num", 2, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.channel_num", 1, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.replica_num", 0, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("metering_service.consumer_connection_num", 1, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)

	// littlebell
	mviper.SetDefault("littlebell.arn", "arn:aws:sns:default:14735dfa-5553-46cc-b4bd-405e711b223f:lbm-svc-event-publish-topic", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.host", "http://sns-service.pegasus-system:8092", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.region", "default", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.access_key", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.secret_key", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.credential.project_id", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
	mviper.SetDefault("littlebell.credential.user_id", "", "-", tkCfg.TypeLocal, tkCfg.RegionGlobal)
}
