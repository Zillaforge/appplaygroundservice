package versions

import "github.com/go-gormigrate/gormigrate/v2"

// Get the all migration list
func Get() (versions []*gormigrate.Migration) {
	return []*gormigrate.Migration{
		getID001Migrate(), getID002Migrate(), getID003Migrate(), getID004Migrate(), getID005Migrate(),
		getID006Migrate(), getID007Migrate(), getID008Migrate(),
	}
}
