package common

type CreateAppCredentialInput struct {
	Name         string
	Description  string
	Unrestricted bool
	OpstkUserID  string
}

type CreateAppCredentialOutput struct {
	ID             string
	Name           string
	Secret         string
	OpstkProjectID string
}

type DeleteAppCredentialInput struct {
	ID          string
	OpstkUserID string
}
