package common

import "context"

const (
	DefaultErrorMsg = "An Error Occurred."
)

type Provider interface {
	Deploy(ctx context.Context, input DeployInput) (output DeployOutput, err error)
	Destroy(ctx context.Context, input DestroyInput) (err error)
	GetLogs(ctx context.Context, input GetLogsInput) (output GetLogsOutput, err error)
	GetSummaryLog(ctx context.Context, input GetSummaryLogInput) (output GetSummaryLogOutput, err error)
}
