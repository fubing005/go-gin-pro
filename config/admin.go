package config

type Admin struct {
	ExcludeAuthPath string `mapstructure:"exclude_auth_path" json:"exclude_auth_path" yaml:"exclude_auth_path"`
}
