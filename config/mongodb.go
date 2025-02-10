package config

type MongoDB struct {
	Uri      string `mapstructure:"uri" json:"uri" yaml:"uri"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
}
