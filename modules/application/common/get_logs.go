package common

type GetLogsInput struct {
	ApplicationID string
	ProjectID     string
}

type GetLogsOutput struct {
	Logs string
}
