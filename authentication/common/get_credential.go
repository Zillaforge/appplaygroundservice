package common

type GetCredentialInput struct {
	UserID    string
	ProjectID string
	_         struct{}
}

type GetCredentialOutput struct {
	AccessKey string
	SecretKey string
	_         struct{}
}
