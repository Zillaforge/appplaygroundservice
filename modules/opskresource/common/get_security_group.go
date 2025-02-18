package common

type GetSecurityGroupInput struct {
	ID string
}

type GetSecurityGroupOutput struct {
	ID        string
	Name      string
	UserID    string
	ProjectID string
	Namespace string
}
