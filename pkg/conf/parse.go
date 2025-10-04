package conf

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/yanking/app-skeleton/pkg/log"
)

var (
	mu sync.RWMutex
)

func Parse(configFile string, obj any, reloads ...func()) error {
	// 创建独立的viper实例，避免全局实例带来的冲突
	v := viper.New()
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read configs file %s: %w", configFile, err)
	}

	mu.Lock()
	err := v.Unmarshal(obj)
	mu.Unlock()

	if err != nil {
		return fmt.Errorf("failed to unmarshal configs: %w", err)
	}

	if len(reloads) > 0 {
		watchConfig(v, obj, reloads...)
	}

	return nil
}

func watchConfig(v *viper.Viper, obj any, reloads ...func()) {
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		mu.Lock()
		err := v.Unmarshal(obj)
		mu.Unlock()

		if err != nil {
			log.Errorf("conf.watchConfig: viper.Unmarshal error: %v", err)
		} else {
			// 将 defer/recover 移到循环外面，对所有 reload 函数提供统一的 panic 保护
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("conf.watchConfig: reload function panic: %v", r)
				}
			}()

			for _, reload := range reloads {
				reload()
			}
		}
	})
}
