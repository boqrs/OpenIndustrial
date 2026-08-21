package config

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"

	config "github.com/boqrs/nexus/config/v2"
	"github.com/boqrs/nexus/email"
	"github.com/boqrs/nexus/log"
	"github.com/boqrs/nexus/media"
	"github.com/boqrs/nexus/redis"
	"github.com/boqrs/nexus/tracing"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security/provider"

)

type MyAppConfig struct {
	DBCfg    config.DBConfig `json:"db_cfg" yaml:"db_cfg" mapstructure:"db_cfg"`
	RedisCfg redis.Config    `json:"redis_cfg" yaml:"redis_cfg" mapstructure:"redis_cfg"`
	Trace    tracing.Config  `json:"tracing_cfg" yaml:"tracing_cfg" mapstructure:"tracing_cfg"`
	Media    media.Config    `json:"media_cfg" yaml:"media_cfg" mapstructure:"media_cfg"`
	LogCfg   log.LogConfig   `json:"log_cfg" yaml:"log_cfg" mapstructure:"log_cfg"`
	EmailCfg email.Config    `json:"email_cfg" yaml:"email_cfg" mapstructure:"email_cfg"`
	// Demo 业务配置
	UserJwtSecret string `json:"user_jwt_secret" yaml:"user_jwt_secret" mapstructure:"user_jwt_secret"`
	Ca			provider.ProviderConfig `json:"ca" yaml:"ca" mapstructure:"ca"`
}

// Reload implements comm/config.ConfigReloader.
// It unmarshals the provided map into the configuration struct.
func (c *MyAppConfig) Reload(config map[string]interface{}) error {
	// DecodeHook 可以处理一些特殊的类型转换，比如 string -> time.Duration
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           c,
		WeaklyTypedInput: true,   // 允许弱类型转换，例如 float64 转 int
		TagName:          "yaml", // 【关键修改】指定使用 json tag 进行字段匹配
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(config); err != nil {
		return fmt.Errorf("failed to decode config: %w", err)
	}

	fmt.Printf("Config reloaded successfully, Result is: %#v", config)
	return nil
}

func InitConfig() (*MyAppConfig, error) {
	cfg, _, err := InitConfigWithManager()
	return cfg, err
}

// InitConfigWithManager returns both the config struct and the underlying ConfigManager.
//
// This is useful when you want to register additional reloaders (e.g., provider
// wrappers for database/redis/log/tracing/media) so they can be hot-reloaded
// together with the config.
func InitConfigWithManager() (*MyAppConfig, *config.ConfigManager, error) {
	watcherConfig := &config.FileWatcherConfig{
		FilePath:       "config.yaml",
		FileType:       "yaml",
		EnableFsNotify: true,
	}

	watcher, err := config.NewFileConfigWatcher(watcherConfig)
	if err != nil {
		fmt.Printf("Failed to create watcher: %v", err)
		return nil, nil, err
	}

	configManager := config.NewConfigManager(watcher)

	appConfig := &MyAppConfig{}
	configManager.AddReloader(appConfig)

	ctx := context.Background()
	if err := configManager.Start(ctx); err != nil {
		fmt.Printf("Failed to start: %v", err)
		return nil, nil, err
	}
	// defer configManager.Stop()
	fmt.Printf("Current Service: %#v\n", appConfig)

	return appConfig, configManager, nil
}