package config

type SmsTencent struct {
	SecretID         string `mapstructure:"secret_id" json:"secret_id" yaml:"secret_id"`
	SecretKey        string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	Endpoint         string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	SignMethod       string `mapstructure:"sign_method" json:"sign_method" yaml:"sign_method"`
	Region           string `mapstructure:"region" json:"region" yaml:"region"`
	SmsSdkAppId      string `mapstructure:"sms_sdk_app_id" json:"sms_sdk_app_id" yaml:"sms_sdk_app_id"`
	SignName         string `mapstructure:"sign_name" json:"sign_name" yaml:"sign_name"`
	TemplateIdCommon string `mapstructure:"template_id_common" json:"template_id_common" yaml:"template_id_common"`
	SessionKey       string `mapstructure:"session_key" json:"session_key" yaml:"session_key"`
	CaptchaExpiring  int    `mapstructure:"captcha_expiring" json:"captcha_expiring" yaml:"captcha_expiring"`
	CaptchaInterval  int    `mapstructure:"captcha_interval" json:"captcha_interval" yaml:"captcha_interval"`
}
