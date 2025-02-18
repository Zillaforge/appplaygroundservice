package tables

import (
	"time"

	"gorm.io/datatypes"
)

// Metering
type Metering struct {
	ApplicationID string         `gorm:"not null;primary_key;type:varchar(36)"`
	Name          string         `gorm:"not null;type:varchar(36)"`
	ProjectID     string         `gorm:"not null;type:varchar(36)"`
	Creator       string         `gorm:"not null;type:varchar(36)"`
	Instances     datatypes.JSON `gorm:"check:json_valid(instances);default:'{}'"`
	// CreatedAt 表示 Application 建立時間
	CreatedAt time.Time
	// EndedAt 表示 Application 結束時間
	EndedAt *time.Time
	// LastPublishedAt 表示上次拋量時間
	LastPublishedAt *time.Time
	_               struct{}
}

/*
MeteringInstance 與 MeteringInstanceInfo 為定義 Metering 資料表中 Instances 的資料結構
*/
type MeteringInstanceInfo struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Name     string `json:"name"`
	PortID   string `json:"port_id"`
	Provider string `json:"provider"`
	Type     string `json:"type"`
	FlavorID string `json:"flavor_id"`
	_        struct{}
}

type MeteringInstance struct {
	Instance MeteringInstanceInfo `json:"instance"`
	_        struct{}
}
