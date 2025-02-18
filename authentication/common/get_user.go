package common

import "encoding/json"

type GetUserInput struct {
	ID        string
	Cacheable bool
}

type GetUserOutput struct {
	ID          string                 `json:"id"`
	DisplayName string                 `json:"displayName"`
	Account     string                 `json:"account"`
	Email       string                 `json:"email"`
	Frozen      bool                   `json:"frozen"`
	Extra       map[string]interface{} `json:"extra"`
}

func (p *GetUserOutput) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	b, _ := json.Marshal(p)
	json.Unmarshal(b, &m)
	return m
}
