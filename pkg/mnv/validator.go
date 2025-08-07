package mnv

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// RegisterValidators регистрирует все кастомные валидаторы
func RegisterValidators(v *validator.Validate) error {
	validators := map[string]validator.Func{
		"phonebycountry": validatePhoneByCountry,
		"kg":             validateKGPhone,
		"ru":             validateRUPhone,
		"kz":             validateKZPhone,
		"uz":             validateUZPhone,
		"tj":             validateTJPhone,
		"tm":             validateTMPhone,
		"us":             validateUSPhone,
		"ca":             validateCAPhone,
		"uk":             validateUKPhone,
		"de":             validateDEPhone,
		"fr":             validateFRPhone,
		"it":             validateITPhone,
		"es":             validateESPhone,
		"nl":             validateNLPhone,
		"tr":             validateTRPhone,
		"cn":             validateCNPhone,
		"in":             validateINPhone,
		"jp":             validateJPPhone,
		"kr":             validateKRPhone,
		"phone":          validateGeneralPhone,
		"mobile":         validateMobilePhone,
	}

	for tag, fn := range validators {
		if err := v.RegisterValidation(tag, fn); err != nil {
			return fmt.Errorf("failed to register validator %s: %w", tag, err)
		}
	}

	return nil
}

// validatePhoneByCountry проверяет телефон по коду страны из поля Country
func validatePhoneByCountry(fl validator.FieldLevel) bool {
	parent := fl.Parent()
	countryField := parent.FieldByName("Country")
	if !countryField.IsValid() {
		return false
	}

	countryCode := normalizeCountryCode(countryField.String())
	phone := cleanPhoneNumber(fl.Field().String())

	return validatePhoneForCountry(phone, countryCode)
}

// validatePhoneForCountry проверяет номер телефона для конкретной страны
func validatePhoneForCountry(phone, countryCode string) bool {
	normalizedCode := normalizeCountryCode(countryCode)
	phoneInfo, ok := CountryPhoneCodes[normalizedCode]
	if !ok {
		return false
	}

	// Базовая проверка формата
	if !validatePhoneFormat(phone) {
		return false
	}

	config := GetConfig()

	if config.StrictMode {
		// Строгая проверка с использованием regex
		matched, err := regexp.MatchString(phoneInfo.Pattern, phone)
		return err == nil && matched
	}

	// Базовая проверка по префиксу и длине
	if !strings.HasPrefix(phone, phoneInfo.Prefix) {
		return false
	}

	phoneWithoutPrefix := strings.TrimPrefix(phone, phoneInfo.Prefix)
	phoneDigits := removeNonDigits(phoneWithoutPrefix)
	phoneLen := len(phoneDigits)

	return phoneLen >= phoneInfo.MinLength && phoneLen <= phoneInfo.MaxLength
}

// ValidatePhone выполняет полную валидацию номера телефона с детальными результатами
func ValidatePhone(phone, countryCode string, options ...*ValidationOptions) *ValidationResult {

	result := &ValidationResult{
		OriginalNumber: phone,
		CountryCode:    countryCode,
	}

	// Применяем опции если переданы
	var opts *ValidationOptions
	if len(options) > 0 && options[0] != nil {
		opts = options[0]
		if opts.Config != nil {
			// Временно устанавливаем конфигурацию
			oldConfig := GetConfig()
			SetConfig(*opts.Config)
			defer SetConfig(oldConfig)
		}
	}

	// Нормализуем код страны
	normalizedCountry := normalizeCountryCode(countryCode)

	// Проверяем, поддерживается ли страна
	info, exists := GetCountryInfo(normalizedCountry)
	if !exists {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("Unsupported country code: %s", countryCode)
		if opts != nil && opts.ReturnSuggestions {
			result.Suggestions = findSimilarCountries(countryCode)
		}
		return result
	}

	result.CountryName = info.CountryName

	// Очищаем номер
	cleanedPhone := cleanPhoneNumber(phone)

	// Базовая валидация
	if !validatePhoneFormat(cleanedPhone) {
		result.IsValid = false
		result.ErrorMessage = "Invalid phone number format"
		if opts != nil && opts.ReturnSuggestions {
			result.Suggestions = suggestCorrections(phone, countryCode)
		}
		return result
	}

	// Валидация для конкретной страны
	if validatePhoneForCountry(cleanedPhone, normalizedCountry) {
		result.IsValid = true
		result.FormattedNumber = cleanedPhone
	} else {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("Invalid phone number for country %s", countryCode)
		if opts != nil && opts.ReturnSuggestions {
			result.Suggestions = suggestCorrections(phone, countryCode)
		}
	}

	return result
}

// BatchValidatePhones выполняет пакетную валидацию номеров телефонов
func BatchValidatePhones(request *BatchValidationRequest) *BatchValidationResponse {
	startTime := time.Now()

	response := &BatchValidationResponse{
		Results: make([]ValidationResult, len(request.Phones)),
		Stats: BatchStats{
			Total:     len(request.Phones),
			ByCountry: make(map[string]int),
		},
	}

	// Определяем количество параллельных воркеров
	parallel := request.Parallel
	if parallel <= 0 {
		parallel = 10 // По умолчанию
	}

	// Канал для результатов
	resultChan := make(chan struct {
		index  int
		result ValidationResult
	}, len(request.Phones))

	// Семафор для ограничения количества горутин
	semaphore := make(chan struct{}, parallel)

	// Запускаем валидацию
	for i, phone := range request.Phones {
		go func(idx int, ph string) {
			semaphore <- struct{}{}        // Захватываем место в семафоре
			defer func() { <-semaphore }() // Освобождаем место

			result := ValidatePhone(ph, "", request.Options)

			resultChan <- struct {
				index  int
				result ValidationResult
			}{idx, *result}
		}(i, phone)
	}

	// Собираем результаты
	for i := 0; i < len(request.Phones); i++ {
		res := <-resultChan
		response.Results[res.index] = res.result

		// Обновляем статистику
		if res.result.IsValid {
			response.Stats.Valid++
			if res.result.CountryCode != "" {
				response.Stats.ByCountry[res.result.CountryCode]++
			}
		} else {
			response.Stats.Invalid++
			response.Stats.Errors++
		}
	}

	response.ProcessingTime = time.Since(startTime).String()
	return response
}

// GetPhoneInfo возвращает детальную информацию о номере телефона
func GetPhoneInfo(phone string) *PhoneInfo {
	// Определяем страну по номеру
	country, found := GetCountryByPhone(phone)
	if !found {
		return &PhoneInfo{
			Number:  phone,
			IsValid: false,
			Type:    PhoneTypeUnknown,
		}
	}

	info, _ := GetCountryInfo(country)
	prefix, localNumber := parsePhoneComponents(phone)
	phoneType := detectPhoneType(phone, country)

	return &PhoneInfo{
		Number:      phone,
		CountryCode: country,
		CountryName: info.CountryName,
		Prefix:      prefix,
		LocalNumber: localNumber,
		IsValid:     validatePhoneForCountry(phone, country),
		Type:        phoneType,
	}
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

func validateTJPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "tj")
}

func validateTMPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "tm")
}

func validateUSPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "us")
}

func validateCAPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "ca")
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

func validateESPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "es")
}

func validateNLPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "nl")
}

func validateTRPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "tr")
}

func validateCNPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "cn")
}

func validateINPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "in")
}

func validateJPPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "jp")
}

func validateKRPhone(fl validator.FieldLevel) bool {
	phone := cleanPhoneNumber(fl.Field().String())
	return validatePhoneForCountry(phone, "kr")
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

// validateMobilePhone валидатор мобильных телефонов (алиас для общего валидатора)
func validateMobilePhone(fl validator.FieldLevel) bool {
	return validateGeneralPhone(fl)
}

// GetCountryByPhone определяет страну по номеру телефона
func GetCountryByPhone(phone string) (string, bool) {
	phone = cleanPhoneNumber(phone)

	for countryCode := range CountryPhoneCodes {
		if validatePhoneForCountry(phone, countryCode) {
			return countryCode, true
		}
	}
	return "", false
}

// FormatPhone форматирует номер телефона по стандарту страны
func FormatPhone(phone, countryCode string) (string, error) {
	return formatPhoneNumber(phone, countryCode)
}

// AddCountry добавляет новую страну в валидатор
func AddCountry(countryCode, prefix, pattern string, minLen, maxLen int) error {
	// Валидация входных параметров
	if countryCode == "" {
		return NewValidationError(ErrorTypeInvalidFormat, "country code cannot be empty", "", "", nil)
	}

	if prefix == "" {
		return NewValidationError(ErrorTypeInvalidFormat, "prefix cannot be empty", "", countryCode, nil)
	}

	if !strings.HasPrefix(prefix, "+") {
		return NewValidationError(ErrorTypeInvalidFormat, "prefix must start with +", "", countryCode, nil)
	}

	if pattern != "" && !validateRegexPattern(pattern) {
		return NewValidationError(ErrorTypeInvalidFormat, "invalid regex pattern", "", countryCode, nil)
	}

	if minLen <= 0 || maxLen <= 0 || minLen > maxLen {
		return NewValidationError(ErrorTypeInvalidFormat, "invalid length parameters", "", countryCode, nil)
	}

	normalizedCode := normalizeCountryCode(countryCode)

	CountryPhoneCodes[normalizedCode] = PhoneCodeInfo{
		Prefix:      prefix,
		Pattern:     pattern,
		MinLength:   minLen,
		MaxLength:   maxLen,
		CountryName: strings.ToUpper(countryCode),
		Description: fmt.Sprintf("%s phone numbers", strings.ToUpper(countryCode)),
	}

	return nil
}

// AddCustomCountry добавляет кастомную страну с полной информацией
func AddCustomCountry(country *CustomCountry) error {
	return AddCountry(
		country.Code,
		country.Prefix,
		country.Pattern,
		country.MinLength,
		country.MaxLength,
	)
}

// RemoveCountry удаляет страну из валидатора
func RemoveCountry(countryCode string) {
	normalizedCode := normalizeCountryCode(countryCode)
	delete(CountryPhoneCodes, normalizedCode)
}

// GetSupportedCountries возвращает список поддерживаемых стран
func GetSupportedCountries() []string {
	countries := make([]string, 0, len(CountryPhoneCodes))
	for code := range CountryPhoneCodes {
		countries = append(countries, code)
	}
	return countries
}

// IsPhoneValid простая проверка валидности номера
func IsPhoneValid(phone, countryCode string) bool {
	return validatePhoneForCountry(cleanPhoneNumber(phone), normalizeCountryCode(countryCode))
}
