package configs

import (
	cnt "AppPlaygroundService/constants"
	"fmt"

	tkCfg "pegasus-cloud.com/aes/toolkits/configs"
	"pegasus-cloud.com/aes/toolkits/mviper"
)

func InitScheduler() {
	mviper.SetDefault("version", cnt.Version, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("kind", "Proxy", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	mviper.SetDefault("app_playground_service_scheduler.instance", fmt.Sprintf("%s-Scheduler", cnt.PascalCaseName), "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service_scheduler.core_grpc.hosts", []string{"0.0.0.0:5111"}, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service_scheduler.core_grpc.tls.enable", false, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service_scheduler.core_grpc.tls.cert_path", "/root/server-csr/APS.pem", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service_scheduler.core_grpc.connection_per_host", 3, "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// tasks
	mviper.SetDefault("app_playground_service_scheduler.tasks.sync_projects.cron_expression", "0 */10 * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service_scheduler.tasks.applications_metering.cron_expression", "0 * * * * * *", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("app_playground_service_scheduler.tasks.applications_metering.metering_service.exchange", "mts", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)
	mviper.SetDefault("app_playground_service_scheduler.tasks.applications_metering.metering_service.routing_key", "applications_metering", "-", tkCfg.TypeLocal, tkCfg.RegionLocal)

	// Event consume
	mviper.SetDefault("event_consume.service", "", "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
	mviper.SetDefault("event_consume.channels", []interface{}{}, "-", tkCfg.TypeRestart, tkCfg.RegionGlobal)
}
