package common

import "github.com/Zillaforge/pegasusiamclient/pb"

type ListProjectsInput struct {
	Limit  int32
	Offset int32
	_      struct{}
}

type ListProjectsOutput struct {
	Projects []*pb.ProjectInfo `json:"projects"`
	Total    int64             `json:"total"`
	_        struct{}
}
