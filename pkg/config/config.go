package config

// Config is the configuration for the bsf cli
type Config struct {
	BuildSafeAPI    string `json:"buildsafe_api"`
	BuildSafeAPITLS bool   `json:"buildsafe_api_tls"`
}
