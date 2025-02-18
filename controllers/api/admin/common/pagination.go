package common

type Pagination struct {
	Limit  int `json:"limit" form:"limit,default=100" binding:"max=100"`
	Offset int `json:"offset" form:"offset,default=0" binding:"min=0"`
	_      struct{}
}
