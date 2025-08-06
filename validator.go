package mnv

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate     *validator.Validate
	registerOnce sync.Once
)

// GetValidator возвращает экземпляр валидатора с зарегистрированными правилами.
func GetValidator() *validator.Validate {
	registerOnce.Do(func() {
		validate = validator.New()
		if err := RegisterValidators(validate); err != nil {
			panic(err) // или логировать ошибку
		}
	})
	return validate
}

// PhoneCodeInfo содержит информацию о телефонных кодах страны
type PhoneCodeInfo struct {
	Prefix    string `json:"prefix"`
	Pattern   string `json:"pattern"`
	MinLength int    `json:"min_length"`
	MaxLength int    `json:"max_length"`
}

// ValidatorConfig конфигурация валидатора
type ValidatorConfig struct {
	AllowSpaces      bool
	AllowDashes      bool
	AllowParentheses bool
	StrictMode       bool // Если true, использует регулярные выражения для точной проверки
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() ValidatorConfig {
	return ValidatorConfig{
		AllowSpaces:      false,
		AllowDashes:      false,
		AllowParentheses: false,
		StrictMode:       true,
	}
}

var config = DefaultConfig()

// SetConfig устанавливает глобальную конфигурацию валидатора
func SetConfig(cfg ValidatorConfig) {
	config = cfg
}

// RegisterValidators регистрирует все кастомные валидаторы
func RegisterValidators(v *validator.Validate) error {
	validators := map[string]validator.Func{
		"phonebycountry": validatePhoneByCountry,
		"kg":             validateKGPhone,
		"ru":             validateRUPhone,
		"kz":             validateKZPhone,
		"uz":             validateUZPhone,
		"us":             validateUSPhone,
		"uk":             validateUKPhone,
		"de":             validateDEPhone,
		"fr":             validateFRPhone,
		"it":             validateITPhone,
		"tr":             validateTRPhone,
		"phone":          validateGeneralPhone,
	}

	for tag, fn := range validators {
		if err := v.RegisterValidation(tag, fn); err != nil {
			return fmt.Errorf("failed to register validator %s: %w", tag, err)
		}
	}

	return nil
}

// cleanPhoneNumber очищает номер телефона от лишних символов
func cleanPhoneNumber(phone string) string {
	if !config.AllowSpaces {
		phone = strings.ReplaceAll(phone, " ", "")
	}
	if !config.AllowDashes {
		phone = strings.ReplaceAll(phone, "-", "")
	}
	if !config.AllowParentheses {
		phone = strings.ReplaceAll(phone, "(", "")
		phone = strings.ReplaceAll(phone, ")", "")
	}
	return strings.TrimSpace(phone)
}

// validatePhoneByCountry проверяет телефон по коду страны из поля Country
func validatePhoneByCountry(fl validator.FieldLevel) bool {
	parent := fl.Parent()
	countryField := parent.FieldByName("Country")
	if !countryField.IsValid() {
		return false
	}

	countryCode := strings.ToLower(countryField.String())
	phone := cleanPhoneNumber(fl.Field().String())

	return validatePhoneForCountry(phone, countryCode)
}

// validatePhoneForCountry проверяет номер телефона для конкретной страны
func validatePhoneForCountry(phone, countryCode string) bool {
	phoneInfo, ok := CountryPhoneCodes[countryCode]
	if !ok {
		return false
	}

	if config.StrictMode {
		matched, err := regexp.MatchString(phoneInfo.Pattern, phone)
		return err == nil && matched
	}

	// Базовая проверка по префиксу и длине
	if !strings.HasPrefix(phone, phoneInfo.Prefix) {
		return false
	}

	phoneWithoutPrefix := strings.TrimPrefix(phone, phoneInfo.Prefix)
	phoneLen := len(phoneWithoutPrefix)

	return phoneLen >= phoneInfo.MinLength && phoneLen <= phoneInfo.MaxLength
}

// Валидаторы для конкретных стран
func validateKGPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "kg")
}

func validateRUPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "ru")
}

func validateKZPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "kz")
}

func validateUZPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "uz")
}

func validateUSPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "us")
}

func validateUKPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "uk")
}

func validateDEPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "de")
}

func validateFRPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "fr")
}

func validateITPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "it")
}

func validateTRPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "tr")
}

// validateGeneralPhone общий валидатор телефона (проверяет по всем странам)
func validateGeneralPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())

	for countryCode := range CountryPhoneCodes {
		if validatePhoneForCountry(phone, countryCode) {
			return true
		}
	}
	return false
}

// GetCountryByPhone определяет страну по номеру телефона
func GetCountryByPhone(phone string) (string, bool) {
	phone = cleanPhoneNumber(phone)

	for countryCode, info := range CountryPhoneCodes {
		if strings.HasPrefix(phone, info.Prefix) {
			return countryCode, true
		}
	}
	return "", false
}

// FormatPhone форматирует номер телефона по стандарту страны
func FormatPhone(phone, countryCode string) (string, error) {
	phone = cleanPhoneNumber(phone)
	countryCode = strings.ToLower(countryCode)

	_, ok := CountryPhoneCodes[countryCode]
	if !ok {
		return "", fmt.Errorf("unsupported country code: %s", countryCode)
	}

	if !validatePhoneForCountry(phone, countryCode) {
		return "", fmt.Errorf("invalid phone number for country %s", countryCode)
	}

	return phone, nil
}

// AddCountry добавляет новую страну в валидатор
func AddCountry(countryCode, prefix, pattern string, minLen, maxLen int) {
	CountryPhoneCodes[strings.ToLower(countryCode)] = PhoneCodeInfo{
		Prefix:    prefix,
		Pattern:   pattern,
		MinLength: minLen,
		MaxLength: maxLen,
	}
}

// RemoveCountry удаляет страну из валидатора
func RemoveCountry(countryCode string) {
	delete(CountryPhoneCodes, strings.ToLower(countryCode))
}

// GetSupportedCountries возвращает список поддерживаемых стран
func GetSupportedCountries() []string {
	countries := make([]string, 0, len(CountryPhoneCodes))
	for code := range CountryPhoneCodes {
		countries = append(countries, code)
	}
	return countries
}
