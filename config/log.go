package config

type Log struct {
	Level           string `mapstructure:"level" json:"level" yaml:"level"`
	RootDir         string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	Filename        string `mapstructure:"filename" json:"filename" yaml:"filename"`
	RequestFilename string `mapstructure:"request_filename" json:"request_filename" yaml:"request_filename"`
	Format          string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine        bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	MaxBackups      int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize         int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
	MaxAge          int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`    // day
	Compress        bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}
