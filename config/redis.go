package config

import "time"

type Redis struct {
	Host         string        `mapstructure:"host" json:"host" yaml:"host"`
	Port         int           `mapstructure:"port" json:"port" yaml:"port"`
	DB           int           `mapstructure:"db" json:"db" yaml:"db"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password"`
	Addrs        []string      `mapstructure:"addrs" json:"addrs" yaml:"addrs"`
	PoolSize     int           `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout"`
}
