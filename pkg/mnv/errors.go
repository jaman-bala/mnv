package mnv

import (
	"fmt"
)

// Предопределенные ошибки валидации
var (
	// ErrInvalidFormat ошибка неверного формата номера
	ErrInvalidFormat = &ValidationError{
		Type:    ErrorTypeInvalidFormat,
		Message: "invalid phone number format",
	}

	// ErrInvalidLength ошибка неверной длины номера
	ErrInvalidLength = &ValidationError{
		Type:    ErrorTypeInvalidLength,
		Message: "invalid phone number length",
	}

	// ErrInvalidPrefix ошибка неверного префикса
	ErrInvalidPrefix = &ValidationError{
		Type:    ErrorTypeInvalidPrefix,
		Message: "invalid country prefix",
	}

	// ErrUnsupportedCountry ошибка неподдерживаемой страны
	ErrUnsupportedCountry = &ValidationError{
		Type:    ErrorTypeUnsupportedCountry,
		Message: "unsupported country code",
	}

	// ErrInvalidCharacters ошибка недопустимых символов
	ErrInvalidCharacters = &ValidationError{
		Type:    ErrorTypeInvalidCharacters,
		Message: "phone number contains invalid characters",
	}

	// ErrMissingPlus ошибка отсутствующего знака +
	ErrMissingPlus = &ValidationError{
		Type:    ErrorTypeMissingPlus,
		Message: "phone number must start with + sign",
	}
)

// NewValidationError создает новую ошибку валидации
func NewValidationError(errorType ErrorType, message, phone, countryCode string, suggestions []string) *ValidationError {
	return &ValidationError{
		Type:        errorType,
		Message:     message,
		Phone:       phone,
		CountryCode: countryCode,
		Suggestions: suggestions,
	}
}

// NewInvalidFormatError создает ошибку неверного формата
func NewInvalidFormatError(phone, countryCode string) *ValidationError {
	return &ValidationError{
		Type:        ErrorTypeInvalidFormat,
		Message:     "invalid phone number format",
		Phone:       phone,
		CountryCode: countryCode,
		Suggestions: suggestCorrections(phone, countryCode),
	}
}

// NewInvalidLengthError создает ошибку неверной длины
func NewInvalidLengthError(phone, countryCode string, expected, actual int) *ValidationError {
	message := fmt.Sprintf("invalid phone number length: expected %d digits, got %d", expected, actual)
	return &ValidationError{
		Type:        ErrorTypeInvalidLength,
		Message:     message,
		Phone:       phone,
		CountryCode: countryCode,
		Suggestions: suggestCorrections(phone, countryCode),
	}
}

// NewInvalidLengthRangeError создает ошибку неверной длины с диапазоном
func NewInvalidLengthRangeError(phone, countryCode string, minLen, maxLen, actual int) *ValidationError {
	message := fmt.Sprintf("invalid phone number length: expected %d-%d digits, got %d", minLen, maxLen, actual)
	return &ValidationError{
		Type:        ErrorTypeInvalidLength,
		Message:     message,
		Phone:       phone,
		CountryCode: countryCode,
		Suggestions: suggestCorrections(phone, countryCode),
	}
}

// NewInvalidPrefixError создает ошибку неверного префикса
func NewInvalidPrefixError(phone, countryCode, expectedPrefix, actualPrefix string) *ValidationError {
	message := fmt.Sprintf("invalid country prefix: expected %s for %s, got %s", expectedPrefix, countryCode, actualPrefix)
	return &ValidationError{
		Type:        ErrorTypeInvalidPrefix,
		Message:     message,
		Phone:       phone,
		CountryCode: countryCode,
		Suggestions: suggestCorrections(phone, countryCode),
	}
}

// NewUnsupportedCountryError создает ошибку неподдерживаемой страны
func NewUnsupportedCountryError(countryCode string) *ValidationError {
	similar := findSimilarCountries(countryCode)
	var suggestions []string

	if len(similar) > 0 {
		for _, code := range similar {
			suggestions = append(suggestions, fmt.Sprintf("Did you mean '%s'?", code))
		}
	}

	message := fmt.Sprintf("unsupported country code: %s", countryCode)
	return &ValidationError{
		Type:        ErrorTypeUnsupportedCountry,
		Message:     message,
		CountryCode: countryCode,
		Suggestions: suggestions,
	}
}

// NewInvalidCharactersError создает ошибку недопустимых символов
func NewInvalidCharactersError(phone string, invalidChars []rune) *ValidationError {
	message := fmt.Sprintf("phone number contains invalid characters: %v", invalidChars)
	return &ValidationError{
		Type:        ErrorTypeInvalidCharacters,
		Message:     message,
		Phone:       phone,
		Suggestions: []string{removeInvalidChars(phone, invalidChars)},
	}
}

// NewMissingPlusError создает ошибку отсутствующего знака +
func NewMissingPlusError(phone string) *ValidationError {
	suggestions := []string{"+" + phone}
	return &ValidationError{
		Type:        ErrorTypeMissingPlus,
		Message:     "phone number must start with + sign",
		Phone:       phone,
		Suggestions: suggestions,
	}
}

// removeInvalidChars удаляет недопустимые символы из номера
func removeInvalidChars(phone string, invalidChars []rune) string {
	invalidSet := make(map[rune]bool)
	for _, char := range invalidChars {
		invalidSet[char] = true
	}

	var result []rune
	for _, char := range phone {
		if !invalidSet[char] {
			result = append(result, char)
		}
	}

	return string(result)
}

// IsValidationError проверяет, является ли ошибка ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// GetValidationError приводит ошибку к типу ValidationError
func GetValidationError(err error) (*ValidationError, bool) {
	ve, ok := err.(*ValidationError)
	return ve, ok
}

// WrapError оборачивает обычную ошибку в ValidationError
func WrapError(err error, errorType ErrorType, phone, countryCode string) *ValidationError {
	return &ValidationError{
		Type:        errorType,
		Message:     err.Error(),
		Phone:       phone,
		CountryCode: countryCode,
	}
}

// ErrorMessages содержит локализованные сообщения об ошибках
var ErrorMessages = map[ErrorType]map[string]string{
	ErrorTypeInvalidFormat: {
		"en": "Invalid phone number format",
		"ru": "Неверный формат номера телефона",
		"kg": "Телефон номерунун форматы туура эмес",
	},
	ErrorTypeInvalidLength: {
		"en": "Invalid phone number length",
		"ru": "Неверная длина номера телефона",
		"kg": "Телефон номерунун узундугу туура эмес",
	},
	ErrorTypeInvalidPrefix: {
		"en": "Invalid country prefix",
		"ru": "Неверный код страны",
		"kg": "Өлкө коду туура эмес",
	},
	ErrorTypeUnsupportedCountry: {
		"en": "Unsupported country code",
		"ru": "Неподдерживаемый код страны",
		"kg": "Колдоого алынбаган өлкө коду",
	},
	ErrorTypeInvalidCharacters: {
		"en": "Phone number contains invalid characters",
		"ru": "Номер содержит недопустимые символы",
		"kg": "Номерде жараксыз символдор бар",
	},
	ErrorTypeMissingPlus: {
		"en": "Phone number must start with + sign",
		"ru": "Номер должен начинаться со знака +",
		"kg": "Номер + белгиси менен башталышы керек",
	},
}

// GetLocalizedMessage возвращает локализованное сообщение об ошибке
func (ve *ValidationError) GetLocalizedMessage(lang string) string {
	if messages, exists := ErrorMessages[ve.Type]; exists {
		if message, langExists := messages[lang]; langExists {
			return message
		}
		// Fallback to English
		if message, enExists := messages["en"]; enExists {
			return message
		}
	}
	return ve.Message
}

// WithSuggestions добавляет предложения к ошибке
func (ve *ValidationError) WithSuggestions(suggestions []string) *ValidationError {
	ve.Suggestions = suggestions
	return ve
}

// WithPhone добавляет номер телефона к ошибке
func (ve *ValidationError) WithPhone(phone string) *ValidationError {
	ve.Phone = phone
	return ve
}

// WithCountryCode добавляет код страны к ошибке
func (ve *ValidationError) WithCountryCode(countryCode string) *ValidationError {
	ve.CountryCode = countryCode
	return ve
}

// String возвращает строковое представление ошибки
func (ve *ValidationError) String() string {
	result := fmt.Sprintf("[%s] %s", ve.Type, ve.Message)

	if ve.Phone != "" {
		result += fmt.Sprintf(" (Phone: %s)", ve.Phone)
	}

	if ve.CountryCode != "" {
		result += fmt.Sprintf(" (Country: %s)", ve.CountryCode)
	}

	if len(ve.Suggestions) > 0 {
		result += fmt.Sprintf(" (Suggestions: %v)", ve.Suggestions)
	}

	return result
}

// ErrorCode возвращает числовой код ошибки
func (ve *ValidationError) ErrorCode() int {
	switch ve.Type {
	case ErrorTypeInvalidFormat:
		return 1001
	case ErrorTypeInvalidLength:
		return 1002
	case ErrorTypeInvalidPrefix:
		return 1003
	case ErrorTypeUnsupportedCountry:
		return 1004
	case ErrorTypeInvalidCharacters:
		return 1005
	case ErrorTypeMissingPlus:
		return 1006
	default:
		return 1000
	}
}

// IsRetryable определяет, можно ли повторить операцию после исправления ошибки
func (ve *ValidationError) IsRetryable() bool {
	switch ve.Type {
	case ErrorTypeInvalidFormat, ErrorTypeInvalidCharacters, ErrorTypeMissingPlus:
		return true
	case ErrorTypeUnsupportedCountry, ErrorTypeInvalidLength, ErrorTypeInvalidPrefix:
		return false
	default:
		return false
	}
}

// ToJSON преобразует ошибку в JSON-совместимую структуру
func (ve *ValidationError) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"type":         string(ve.Type),
		"message":      ve.Message,
		"phone":        ve.Phone,
		"country_code": ve.CountryCode,
		"suggestions":  ve.Suggestions,
		"error_code":   ve.ErrorCode(),
		"retryable":    ve.IsRetryable(),
	}
}
