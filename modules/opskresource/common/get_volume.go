package common

type GetVolumeInput struct {
	ID string
}

type GetVolumeOutput struct {
	ID        string
	Name      string
	ProjectID string
	UserID    string
	Namespace string
	Status    string
	Type      string
	Size      int32
}
