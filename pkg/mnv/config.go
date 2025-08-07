package mnv

import (
	"sync"
	"time"
)

// глобальная конфигурация с мьютексом для thread-safety
var (
	globalConfig ValidatorConfig
	configMutex  sync.RWMutex
)

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() ValidatorConfig {
	return ValidatorConfig{
		AllowSpaces:              false,
		AllowDashes:              false,
		AllowParentheses:         false,
		AllowDots:                false,
		StrictMode:               true,
		RequirePlusSign:          true,
		CaseSensitiveCountryCode: false,
	}
}

// RelaxedConfig возвращает более мягкую конфигурацию
func RelaxedConfig() ValidatorConfig {
	return ValidatorConfig{
		AllowSpaces:              true,
		AllowDashes:              true,
		AllowParentheses:         true,
		AllowDots:                true,
		StrictMode:               false,
		RequirePlusSign:          false,
		CaseSensitiveCountryCode: false,
	}
}

// StrictConfig возвращает строгую конфигурацию
func StrictConfig() ValidatorConfig {
	return ValidatorConfig{
		AllowSpaces:              false,
		AllowDashes:              false,
		AllowParentheses:         false,
		AllowDots:                false,
		StrictMode:               true,
		RequirePlusSign:          true,
		CaseSensitiveCountryCode: true,
	}
}

// SetConfig устанавливает глобальную конфигурацию (thread-safe)
func SetConfig(cfg ValidatorConfig) {
	configMutex.Lock()
	defer configMutex.Unlock()
	globalConfig = cfg
}

// GetConfig возвращает текущую глобальную конфигурацию (thread-safe)
func GetConfig() ValidatorConfig {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return globalConfig
}

// ResetConfig сбрасывает конфигурацию к значениям по умолчанию
func ResetConfig() {
	SetConfig(DefaultConfig())
}

// ConfigBuilder строитель конфигурации для более удобного создания
type ConfigBuilder struct {
	config ValidatorConfig
}

// NewConfigBuilder создает новый строитель конфигурации
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// AllowSpaces разрешает пробелы
func (cb *ConfigBuilder) AllowSpaces(allow bool) *ConfigBuilder {
	cb.config.AllowSpaces = allow
	return cb
}

// AllowDashes разрешает тире
func (cb *ConfigBuilder) AllowDashes(allow bool) *ConfigBuilder {
	cb.config.AllowDashes = allow
	return cb
}

// AllowParentheses разрешает скобки
func (cb *ConfigBuilder) AllowParentheses(allow bool) *ConfigBuilder {
	cb.config.AllowParentheses = allow
	return cb
}

// AllowDots разрешает точки
func (cb *ConfigBuilder) AllowDots(allow bool) *ConfigBuilder {
	cb.config.AllowDots = allow
	return cb
}

// StrictMode включает/выключает строгий режим
func (cb *ConfigBuilder) StrictMode(strict bool) *ConfigBuilder {
	cb.config.StrictMode = strict
	return cb
}

// RequirePlusSign требует знак плюс
func (cb *ConfigBuilder) RequirePlusSign(require bool) *ConfigBuilder {
	cb.config.RequirePlusSign = require
	return cb
}

// CaseSensitiveCountryCode делает коды стран чувствительными к регистру
func (cb *ConfigBuilder) CaseSensitiveCountryCode(sensitive bool) *ConfigBuilder {
	cb.config.CaseSensitiveCountryCode = sensitive
	return cb
}

// Build возвращает построенную конфигурацию
func (cb *ConfigBuilder) Build() ValidatorConfig {
	return cb.config
}

// ValidateConfig проверяет корректность конфигурации
func ValidateConfig(cfg ValidatorConfig) error {
	// Проверяем логические противоречия
	if !cfg.StrictMode && cfg.RequirePlusSign {
		return &ValidationError{
			Type:    ErrorTypeInvalidFormat,
			Message: "RequirePlusSign can only be true when StrictMode is enabled",
		}
	}

	return nil
}

// ConfigPresets предустановленные конфигурации
var ConfigPresets = map[string]ValidatorConfig{
	"default": DefaultConfig(),
	"relaxed": RelaxedConfig(),
	"strict":  StrictConfig(),
	"international": {
		AllowSpaces:              true,
		AllowDashes:              false,
		AllowParentheses:         false,
		AllowDots:                false,
		StrictMode:               true,
		RequirePlusSign:          true,
		CaseSensitiveCountryCode: false,
	},
	"local": {
		AllowSpaces:              true,
		AllowDashes:              true,
		AllowParentheses:         true,
		AllowDots:                true,
		StrictMode:               false,
		RequirePlusSign:          false,
		CaseSensitiveCountryCode: false,
	},
	"api": {
		AllowSpaces:              false,
		AllowDashes:              false,
		AllowParentheses:         false,
		AllowDots:                false,
		StrictMode:               true,
		RequirePlusSign:          true,
		CaseSensitiveCountryCode: false,
	},
}

// GetPresetConfig возвращает предустановленную конфигурацию
func GetPresetConfig(preset string) (ValidatorConfig, bool) {
	cfg, exists := ConfigPresets[preset]
	return cfg, exists
}

// SetPresetConfig устанавливает предустановленную конфигурацию как глобальную
func SetPresetConfig(preset string) error {
	cfg, exists := GetPresetConfig(preset)
	if !exists {
		return &ValidationError{
			Type:    ErrorTypeUnknown,
			Message: "unknown configuration preset: " + preset,
		}
	}

	if err := ValidateConfig(cfg); err != nil {
		return err
	}

	SetConfig(cfg)
	return nil
}

// ListPresets возвращает список доступных предустановок
func ListPresets() []string {
	presets := make([]string, 0, len(ConfigPresets))
	for name := range ConfigPresets {
		presets = append(presets, name)
	}
	return presets
}

// CacheConfig конфигурация кеширования
type CacheConfig struct {
	// Enabled включает кеширование
	Enabled bool `json:"enabled"`

	// MaxSize максимальный размер кеша
	MaxSize int `json:"max_size"`

	// TTL время жизни записей в секундах
	TTL int64 `json:"ttl"`

	// CleanupInterval интервал очистки кеша в секундах
	CleanupInterval int64 `json:"cleanup_interval"`
}

// DefaultCacheConfig возвращает конфигурацию кеша по умолчанию
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Enabled:         false, // По умолчанию отключен
		MaxSize:         10000, // 10K записей
		TTL:             3600,  // 1 час
		CleanupInterval: 300,   // 5 минут
	}
}

// кеш конфигурация
var (
	cacheConfig      CacheConfig
	cacheConfigMutex sync.RWMutex
)

// SetCacheConfig устанавливает конфигурацию кеша
func SetCacheConfig(cfg CacheConfig) {
	cacheConfigMutex.Lock()
	defer cacheConfigMutex.Unlock()
	cacheConfig = cfg
}

// GetCacheConfig возвращает текущую конфигурацию кеша
func GetCacheConfig() CacheConfig {
	cacheConfigMutex.RLock()
	defer cacheConfigMutex.RUnlock()
	return cacheConfig
}

// PerformanceConfig конфигурация производительности
type PerformanceConfig struct {
	// EnableProfiling включает профилирование
	EnableProfiling bool `json:"enable_profiling"`

	// MaxConcurrentValidations максимальное количество одновременных валидаций
	MaxConcurrentValidations int `json:"max_concurrent_validations"`

	// ValidationTimeout таймаут валидации в миллисекундах
	ValidationTimeout time.Duration `json:"validation_timeout"`

	// EnableMetrics включает сбор метрик
	EnableMetrics bool `json:"enable_metrics"`
}

// DefaultPerformanceConfig возвращает конфигурацию производительности по умолчанию
func DefaultPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		EnableProfiling:          false,
		MaxConcurrentValidations: 100,
		ValidationTimeout:        time.Millisecond * 100,
		EnableMetrics:            false,
	}
}

// конфигурация производительности
var (
	performanceConfig      PerformanceConfig
	performanceConfigMutex sync.RWMutex
)

// SetPerformanceConfig устанавливает конфигурацию производительности
func SetPerformanceConfig(cfg PerformanceConfig) {
	performanceConfigMutex.Lock()
	defer performanceConfigMutex.Unlock()
	performanceConfig = cfg
}

// GetPerformanceConfig возвращает текущую конфигурацию производительности
func GetPerformanceConfig() PerformanceConfig {
	performanceConfigMutex.RLock()
	defer performanceConfigMutex.RUnlock()
	return performanceConfig
}

// init инициализирует конфигурации по умолчанию
func init() {
	globalConfig = DefaultConfig()
	cacheConfig = DefaultCacheConfig()
	performanceConfig = DefaultPerformanceConfig()
}
