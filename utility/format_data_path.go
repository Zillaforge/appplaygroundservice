package utility

import (
	"path/filepath"

	mviper "pegasus-cloud.com/aes/toolkits/mviper"
)

// /<data_dir>/<module_pid>/aps/module/<module-id>/
func FormatModulePath(moduleID string) string {
	return filepath.Join(
		mviper.GetString("app_playground_service.data_path.data_dir"),
		mviper.GetString("app_playground_service.data_path.module_pid"),
		"aps",
		"module",
		moduleID,
	) + string(filepath.Separator)
}

// /<data_dir>/<project-id>/aps/apllication/<apllication-id>/
func FormatApplicationPath(projectID string, applicationID string) string {
	return filepath.Join(
		mviper.GetString("app_playground_service.data_path.data_dir"),
		projectID,
		"aps",
		"application",
		applicationID,
	) + string(filepath.Separator)
}

// /aps/module/<module-id>/
func FormatShortModulePath(moduleID string) string {
	return filepath.Join(
		string(filepath.Separator),
		"aps",
		"module",
		moduleID,
	) + string(filepath.Separator)
}
