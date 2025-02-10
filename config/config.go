package config

type Configuration struct {
	App           App           `mapstructure:"app" json:"app" yaml:"app"`
	Admin         Admin         `mapstructure:"admin" json:"admin" yaml:"admin"`
	Log           Log           `mapstructure:"log" json:"log" yaml:"log"`
	Database      Database      `mapstructure:"database" json:"database" yaml:"database"`
	Jwt           Jwt           `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis         Redis         `mapstructure:"redis" json:"redis" yaml:"redis"`
	Storage       Storage       `mapstructure:"storage" json:"storage" yaml:"storage"`
	SmsTencent    SmsTencent    `mapstructure:"sms_tencent" json:"sms_tencent" yaml:"sms_tencent"`
	Queue         Queue         `mapstructure:"queue" json:"queue" yaml:"queue"`
	MongoDB       MongoDB       `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	Elasticsearch Elasticsearch `mapstructure:"elasticsearch" json:"elasticsearch" yaml:"elasticsearch"`
	Clickhouse    Clickhouse    `mapstructure:"clickhouse" json:"clickhouse" yaml:"clickhouse"`
}
