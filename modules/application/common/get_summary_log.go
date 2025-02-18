package common

type GetSummaryLogInput struct {
	ApplicationID string
	ProjectID     string
}

type GetSummaryLogOutput struct {
	Log string
}
