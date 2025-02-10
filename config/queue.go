package config

type Queue struct {
	Rabbitmq Rabbitmq `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	Kafka    Kafka    `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
}

type Rabbitmq struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Vhost    string `mapstructure:"vhost" json:"vhost" yaml:"vhost"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	Topic   string   `mapstructure:"topic" json:"topic" yaml:"topic"`
	GroupId string   `mapstructure:"group_id" json:"group_id" yaml:"group_id"`
}
