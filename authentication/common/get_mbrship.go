package common

// GetMembershipInput defines input structure for get membership
type GetMembershipInput struct {
	UserId    string
	ProjectId string
}

// GetMembershipOutput defines output structure for get membership
type GetMembershipOutput struct {
	TenantRole string                 `json:"tenantRole"`
	Frozen     bool                   `json:"frozen"`
	Extra      map[string]interface{} `json:"extra"`
}
