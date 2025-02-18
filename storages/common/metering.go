package common

import (
	"AppPlaygroundService/storages/tables"
	"time"
)

type (
	CreateMeteringInput struct {
		Metering tables.Metering
		_        struct{}
	}

	CreateMeteringOutput struct {
		Metering tables.Metering
		_        struct{}
	}

	ListMeteringsInput struct {
		Pagination *Pagination
		_          struct{}
	}

	ListMeteringsOutput struct {
		Meterings []tables.Metering
		_         struct{}
	}

	CountMeteringInput struct {
		_ struct{}
	}

	CountMeteringOutput struct {
		Count int64
		_     struct{}
	}

	MeteringUpdateInfo struct {
		EndedAt         time.Time
		LastPublishedAt time.Time
		_               struct{}
	}

	UpdateMeteringInput struct {
		ApplicationID string
		UpdateData    *MeteringUpdateInfo
		_             struct{}
	}

	UpdateMeteringOutput struct {
		Metering tables.Metering
		_        struct{}
	}

	GetMeteringInput struct {
		ID string
		_  struct{}
	}

	GetMeteringOutput struct {
		Metering tables.Metering
		_        struct{}
	}

	DeleteMeteringInput struct {
		ID string
		_  struct{}
	}

	DeleteMeteringOutput struct {
		_ struct{}
	}
)
