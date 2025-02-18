package tasks

import (
	"fmt"

	"pegasus-cloud.com/aes/toolkits/mviper"
	"pegasus-cloud.com/aes/toolkits/scheduler"
)

func InitSchedulerTasks() {
	{ // Sync Projects
		period := mviper.GetString("app_playground_service_scheduler.tasks.sync_projects.cron_expression")
		task := scheduler.CreateSchedulerV2(period)
		if err := task.Time().Do(SyncProjects); err != nil {
			panic(fmt.Sprintf("start SyncProjects(cronExp: %s) failed: %s", period, err.Error()))
		}
		task.Start()
		SyncProjects()
	}
	{ // Applications 計量排程
		period := mviper.GetString("app_playground_service_scheduler.tasks.applications_metering.cron_expression")
		task := scheduler.CreateSchedulerV2(period)
		if err := task.Time().Do(ApplicationsMetering); err != nil {
			panic(fmt.Sprintf("start ApplicationPricing(cronExp: %s) failed: %s", period, err.Error()))
		}
		task.Start()
		ApplicationsMetering()
	}
}
