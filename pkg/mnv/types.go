package mnv

import "time"

// PhoneCodeInfo содержит информацию о телефонных кодах страны
type PhoneCodeInfo struct {
	// Prefix - телефонный префикс страны (например, "+996")
	Prefix string `json:"prefix"`

	// Pattern - регулярное выражение для строгой валидации
	Pattern string `json:"pattern"`

	// MinLength - минимальная длина номера без префикса
	MinLength int `json:"min_length"`

	// MaxLength - максимальная длина номера без префикса
	MaxLength int `json:"max_length"`

	// CountryName - полное название страны
	CountryName string `json:"country_name"`

	// Description - описание формата номера
	Description string `json:"description"`
}

// ValidatorConfig конфигурация валидатора
type ValidatorConfig struct {
	// AllowSpaces разрешает пробелы в номере телефона
	AllowSpaces bool `json:"allow_spaces"`

	// AllowDashes разрешает тире в номере телефона
	AllowDashes bool `json:"allow_dashes"`

	// AllowParentheses разрешает скобки в номере телефона
	AllowParentheses bool `json:"allow_parentheses"`

	// AllowDots разрешает точки в номере телефона
	AllowDots bool `json:"allow_dots"`

	// StrictMode использует регулярные выражения для точной проверки
	StrictMode bool `json:"strict_mode"`

	// RequirePlusSign требует обязательное наличие знака +
	RequirePlusSign bool `json:"require_plus_sign"`

	// CaseSensitiveCountryCode делает коды стран чувствительными к регистру
	CaseSensitiveCountryCode bool `json:"case_sensitive_country_code"`
}

// ValidationResult результат валидации номера телефона
type ValidationResult struct {
	// IsValid показывает, прошел ли номер валидацию
	IsValid bool `json:"is_valid"`

	// CountryCode код страны (если определен)
	CountryCode string `json:"country_code,omitempty"`

	// CountryName название страны (если определено)
	CountryName string `json:"country_name,omitempty"`

	// FormattedNumber отформатированный номер
	FormattedNumber string `json:"formatted_number,omitempty"`

	// OriginalNumber оригинальный номер
	OriginalNumber string `json:"original_number"`

	// ErrorMessage сообщение об ошибке (если есть)
	ErrorMessage string `json:"error_message,omitempty"`

	// Suggestions предложения по исправлению (если есть)
	Suggestions []string `json:"suggestions,omitempty"`
}

// PhoneInfo детальная информация о номере телефона
type PhoneInfo struct {
	// Number номер телефона
	Number string `json:"number"`

	// CountryCode код страны
	CountryCode string `json:"country_code"`

	// CountryName название страны
	CountryName string `json:"country_name"`

	// Prefix префикс страны
	Prefix string `json:"prefix"`

	// LocalNumber локальная часть номера (без префикса)
	LocalNumber string `json:"local_number"`

	// IsValid валиден ли номер
	IsValid bool `json:"is_valid"`

	// Type тип номера (mobile, landline, etc.)
	Type PhoneType `json:"type"`

	// Carrier информация о операторе (если доступно)
	Carrier *CarrierInfo `json:"carrier,omitempty"`
}

// PhoneType тип телефонного номера
type PhoneType string

const (
	// PhoneTypeMobile мобильный номер
	PhoneTypeMobile PhoneType = "mobile"

	// PhoneTypeLandline стационарный номер
	PhoneTypeLandline PhoneType = "landline"

	// PhoneTypeTollFree бесплатный номер
	PhoneTypeTollFree PhoneType = "toll_free"

	// PhoneTypePremium премиум номер
	PhoneTypePremium PhoneType = "premium"

	// PhoneTypeVoip VoIP номер
	PhoneTypeVoip PhoneType = "voip"

	// PhoneTypeUnknown неизвестный тип
	PhoneTypeUnknown PhoneType = "unknown"
)

// CarrierInfo информация о мобильном операторе
type CarrierInfo struct {
	// Name название оператора
	Name string `json:"name"`

	// Code код оператора
	Code string `json:"code"`

	// Country страна оператора
	Country string `json:"country"`

	// Type тип оператора (GSM, CDMA, etc.)
	Type string `json:"type"`
}

// CountryStats статистика по стране
type CountryStats struct {
	// CountryCode код страны
	CountryCode string `json:"country_code"`

	// CountryName название страны
	CountryName string `json:"country_name"`

	// TotalValidated общее количество проверенных номеров
	TotalValidated int `json:"total_validated"`

	// ValidNumbers количество валидных номеров
	ValidNumbers int `json:"valid_numbers"`

	// InvalidNumbers количество невалидных номеров
	InvalidNumbers int `json:"invalid_numbers"`

	// SuccessRate процент успешных проверок
	SuccessRate float64 `json:"success_rate"`
}

// ValidationOptions опции для валидации
type ValidationOptions struct {
	// Config конфигурация валидатора
	Config *ValidatorConfig `json:"config,omitempty"`

	// ExpectedCountry ожидаемая страна (для оптимизации)
	ExpectedCountry string `json:"expected_country,omitempty"`

	// AllowedCountries список разрешенных стран
	AllowedCountries []string `json:"allowed_countries,omitempty"`

	// ForbiddenCountries список запрещенных стран
	ForbiddenCountries []string `json:"forbidden_countries,omitempty"`

	// ReturnSuggestions возвращать предложения по исправлению
	ReturnSuggestions bool `json:"return_suggestions"`

	// MaxSuggestions максимальное количество предложений
	MaxSuggestions int `json:"max_suggestions"`
}

// BatchValidationRequest запрос на пакетную валидацию
type BatchValidationRequest struct {
	// Phones список номеров для проверки
	Phones []string `json:"phones"`

	// Options опции валидации
	Options *ValidationOptions `json:"options,omitempty"`

	// Parallel количество параллельных проверок
	Parallel int `json:"parallel"`
}

// BatchValidationResponse ответ на пакетную валидацию
type BatchValidationResponse struct {
	// Results результаты валидации
	Results []ValidationResult `json:"results"`

	// Stats общая статистика
	Stats BatchStats `json:"stats"`

	// ProcessingTime время обработки
	ProcessingTime string `json:"processing_time"`
}

// BatchStats статистика пакетной обработки
type BatchStats struct {
	// Total общее количество номеров
	Total int `json:"total"`

	// Valid количество валидных номеров
	Valid int `json:"valid"`

	// Invalid количество невалидных номеров
	Invalid int `json:"invalid"`

	// ByCountry статистика по странам
	ByCountry map[string]int `json:"by_country"`

	// Errors количество ошибок
	Errors int `json:"errors"`
}

// CustomCountry структура для добавления кастомной страны
type CustomCountry struct {
	// Code код страны (ISO 3166-1 alpha-2)
	Code string `json:"code" binding:"required"`

	// Name название страны
	Name string `json:"name" binding:"required"`

	// Prefix телефонный префикс
	Prefix string `json:"prefix" binding:"required"`

	// Pattern регулярное выражение для валидации
	Pattern string `json:"pattern" binding:"required"`

	// MinLength минимальная длина номера
	MinLength int `json:"min_length" binding:"required,min=1"`

	// MaxLength максимальная длина номера
	MaxLength int `json:"max_length" binding:"required,min=1"`

	// Description описание
	Description string `json:"description"`
}

// ErrorType тип ошибки валидации
type ErrorType string

const (
	// ErrorTypeInvalidFormat неверный формат
	ErrorTypeInvalidFormat ErrorType = "invalid_format"

	// ErrorTypeInvalidLength неверная длина
	ErrorTypeInvalidLength ErrorType = "invalid_length"

	// ErrorTypeInvalidPrefix неверный префикс
	ErrorTypeInvalidPrefix ErrorType = "invalid_prefix"

	// ErrorTypeUnsupportedCountry неподдерживаемая страна
	ErrorTypeUnsupportedCountry ErrorType = "unsupported_country"

	// ErrorTypeInvalidCharacters недопустимые символы
	ErrorTypeInvalidCharacters ErrorType = "invalid_characters"

	// ErrorTypeMissingPlus отсутствует знак +
	ErrorTypeMissingPlus ErrorType = "missing_plus"

	// ErrorTypeUnknown неизвестная ошибка
	ErrorTypeUnknown ErrorType = "unknown"
)

// ValidationError кастомная ошибка валидации
type ValidationError struct {
	// Type тип ошибки
	Type ErrorType `json:"type"`

	// Message сообщение об ошибке
	Message string `json:"message"`

	// Phone номер телефона, вызвавший ошибку
	Phone string `json:"phone"`

	// CountryCode код страны (если определен)
	CountryCode string `json:"country_code,omitempty"`

	// Suggestions предложения по исправлению
	Suggestions []string `json:"suggestions,omitempty"`
}

// Error реализует интерфейс error
func (ve *ValidationError) Error() string {
	return ve.Message
}

// CacheEntry запись в кеше валидации
type CacheEntry struct {
	// Phone номер телефона
	Phone string `json:"phone"`

	// Result результат валидации
	Result ValidationResult `json:"result"`

	// Timestamp время создания записи
	Timestamp int64 `json:"timestamp"`

	// TTL время жизни записи в секундах
	TTL int64 `json:"ttl"`
}

// IsExpired проверяет, истекла ли запись в кеше
func (ce *CacheEntry) IsExpired() bool {
	if ce.TTL <= 0 {
		return false // Бессрочная запись
	}
	return (ce.Timestamp + ce.TTL) < time.Now().Unix()
}
