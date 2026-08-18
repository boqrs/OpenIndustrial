package provider


type ProviderConfig struct {
	Provider string      `mapstructure:"provider"` // "aws", "aliyun", or "local"
	AWS      AWSConfig   `mapstructure:"aws"`
	Aliyun   AliyunConfig `mapstructure:"aliyun"`
	Local    LocalConfig `mapstructure:"local"`
}

type AWSConfig struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	CAArn     string `mapstructure:"ca_arn"` // Certificate Authority ARN
}

type AliyunConfig struct {
	Endpoint         string `mapstructure:"endpoint"`
	AccessKeyID      string `mapstructure:"access_key_id"`
	AccessKeySecret  string `mapstructure:"access_key_secret"`
	ParentIdentifier string `mapstructure:"parent_identifier"` // CA 实例ID
}

type LocalConfig struct {
	RootCertPath string `mapstructure:"root_cert_path"`
	RootKeyPath  string `mapstructure:"root_key_path"`
}
