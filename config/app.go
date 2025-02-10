package config

type App struct {
	Env          string `mapstructure:"env" json:"env" yaml:"env"`
	Port         string `mapstructure:"port" json:"port" yaml:"port"`
	AppName      string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl       string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
	Lang         string `mapstructure:"lang" json:"lang" yaml:"lang"`
	RequestLimit int    `mapstructure:"request_limit" json:"request_limit" yaml:"request_limit"`
}
