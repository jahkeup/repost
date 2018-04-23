package config

// Config options
type Config struct {
	General      General
	Notification Notification
	Delivery     Delivery
}

// General configuration options for program run.
type General struct {
	LogLevel string

	// AWS Credential Profile
	Profile string
	// RoleArn to assume using default credentials
	RoleArn string
}
