package cmd

// Config keys as constants to avoid repetition.
const (
	configKeyDBPath        = "db.path"
	configKeyDBSyncURL     = "db.sync_url"
	configKeyDBAuthToken   = "db.auth_token" //nolint:gosec // This is a config key name, not a credential
	configKeyDBSyncOnWrite = "db.sync_on_write"
	configKeyUITheme       = "ui.theme"
)
