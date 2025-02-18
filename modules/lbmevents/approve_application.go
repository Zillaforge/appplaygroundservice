package lbmevents

import "pegasus-cloud.com/aes/toolkits/littlebell"

type ApproveApplicationEvent struct{ littlebell.MessageStruct }

type ApproveApplication struct {
	// @message id
	ID string `json:"id"`
	// @message name
	Name string `json:"name"`
	// @message description
	Description string `json:"description"`
	// @message moduleID
	ModuleID string `json:"moduleID"`
	// @message state
	State string `json:"state"`
	// @message answer
	Answer interface{} `json:"answer"`
	// @message namespace
	Namespace string `json:"namespace"`
	// @message shiftable
	Shiftable bool `json:"shiftable"`
	// @message projectID
	ProjectID string `json:"projectID"`
	// @message projectName
	ProjectName string `json:"projectName"`
	// @message reviewerID
	ReviewerID string `json:"reviewerID"`
	// @message reviewerName
	ReviewerName string `json:"reviewerName"`
	// @message userID
	UserID string `json:"userID"`
	// @message userName
	UserName string `json:"userName"`
	// @message updaterID
	UpdaterID string `json:"updaterID"`
	// @message createdAt
	CreatedAt string `json:"createdAt"`
	// @message updatedAt
	UpdatedAt string `json:"updatedAt"`
	// @message ad
	AvailabilityDistrict string `json:"ad"`
	_                    struct{}
}

func (e *ApproveApplicationEvent) Name() string {
	return "APS_APPROVE_APPLICATION"
}
