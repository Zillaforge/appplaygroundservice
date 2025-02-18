package common

type GetKeypairInput struct {
	ID string
}

type GetKeypairOutput struct {
	ID     string
	Name   string
	UserID string
}
