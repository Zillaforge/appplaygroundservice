package scheduler

import (
	gArgs "AppPlaygroundService/cmd/args"
	gComm "AppPlaygroundService/cmd/common"
	"AppPlaygroundService/configs"
	cnt "AppPlaygroundService/constants"
	"AppPlaygroundService/server"
	"fmt"

	"github.com/spf13/cobra"
)

func StartCmd() (cmd *cobra.Command) {
	description := "Start %s Scheduler Service"
	cmd = &cobra.Command{
		Use:   "start",
		Short: fmt.Sprintf(description, cnt.UpperAbbrName),
		Long:  fmt.Sprintf(description, cnt.PascalCaseName),
		Run: func(cmd *cobra.Command, args []string) {
			server.RunScheduler()
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			configs.InitScheduler()
			gComm.MergeConfig(gArgs.CfgFileScheduler)
		},
	}
	return
}
