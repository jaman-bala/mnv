package mnv

import (
	"regexp"
	"strings"
	"unicode"
)

// cleanPhoneNumber очищает номер телефона от лишних символов согласно конфигурации
func cleanPhoneNumber(phone string) string {
	config := GetConfig()

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
	if !config.AllowDots {
		phone = strings.ReplaceAll(phone, ".", "")
	}

	return strings.TrimSpace(phone)
}

// normalizeCountryCode нормализует код страны согласно конфигурации
func normalizeCountryCode(countryCode string) string {
	config := GetConfig()

	if !config.CaseSensitiveCountryCode {
		return strings.ToLower(countryCode)
	}

	return countryCode
}

// validatePhoneFormat проверяет базовый формат номера телефона
func validatePhoneFormat(phone string) bool {
	if len(phone) == 0 {
		return false
	}

	config := GetConfig()

	// Проверяем наличие знака + если требуется
	if config.RequirePlusSign && !strings.HasPrefix(phone, "+") {
		return false
	}

	// Убираем + для дальнейшей проверки
	phoneWithoutPlus := strings.TrimPrefix(phone, "+")

	// Проверяем, что остались только цифры и разрешенные символы
	for _, char := range phoneWithoutPlus {
		if unicode.IsDigit(char) {
			continue
		}

		switch char {
		case ' ':
			if !config.AllowSpaces {
				return false
			}
		case '-':
			if !config.AllowDashes {
				return false
			}
		case '(', ')':
			if !config.AllowParentheses {
				return false
			}
		case '.':
			if !config.AllowDots {
				return false
			}
		default:
			return false // Недопустимый символ
		}
	}

	return true
}

// extractDigitsOnly извлекает только цифры из номера телефона
func extractDigitsOnly(phone string) string {
	var digits strings.Builder
	for _, char := range phone {
		if unicode.IsDigit(char) {
			digits.WriteRune(char)
		}
	}
	return digits.String()
}

// suggestCorrections предлагает исправления для неверного номера
func suggestCorrections(phone, countryCode string) []string {
	var suggestions []string

	// Если нет знака +, добавляем его
	if !strings.HasPrefix(phone, "+") {
		suggestion := "+" + phone
		suggestions = append(suggestions, suggestion)
	}

	// Если есть код страны, пробуем добавить правильный префикс
	if countryCode != "" {
		if info, exists := GetCountryInfo(countryCode); exists {
			digitsOnly := extractDigitsOnly(phone)

			// Убираем код страны если он есть в начале
			prefixDigits := extractDigitsOnly(info.Prefix)
			if strings.HasPrefix(digitsOnly, prefixDigits) {
				localPart := digitsOnly[len(prefixDigits):]
				suggestion := info.Prefix + localPart
				suggestions = append(suggestions, suggestion)
			} else {
				// Добавляем префикс к номеру
				suggestion := info.Prefix + digitsOnly
				suggestions = append(suggestions, suggestion)
			}
		}
	}

	// Пробуем определить страну по длине и добавить соответствующий префикс
	digitsOnly := extractDigitsOnly(phone)
	for _, info := range CountryPhoneCodes {
		if len(digitsOnly) >= info.MinLength && len(digitsOnly) <= info.MaxLength {
			// Убираем уже возможный дублирующий префикс перед добавлением нового
			trimmed := strings.TrimPrefix(phone, info.Prefix)
			suggestion := info.Prefix + extractDigitsOnly(trimmed)
			if !contains(suggestions, suggestion) {
				suggestions = append(suggestions, suggestion)
			}
		}
	}

	// Ограничиваем количество предложений
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return suggestions
}

// contains проверяет, содержит ли слайс строку
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// formatPhoneNumber форматирует номер телефона согласно стандартам страны
func formatPhoneNumber(phone, countryCode string) (string, error) {
	_, exists := GetCountryInfo(countryCode)
	if !exists {
		return "", &ValidationError{
			Type:    ErrorTypeUnsupportedCountry,
			Message: "unsupported country code: " + countryCode,
			Phone:   phone,
		}
	}

	cleaned := cleanPhoneNumber(phone)

	if !validatePhoneForCountry(cleaned, countryCode) {
		return "", &ValidationError{
			Type:        ErrorTypeInvalidFormat,
			Message:     "invalid phone number format for country " + countryCode,
			Phone:       phone,
			CountryCode: countryCode,
		}
	}

	return cleaned, nil
}

// detectPhoneType определяет тип номера телефона (примитивная реализация)
func detectPhoneType(phone, countryCode string) PhoneType {
	// Базовая логика определения типа номера
	// В реальном приложении здесь была бы более сложная логика

	switch countryCode {
	case "us", "ca":
		// В США и Канаде номера 800, 888, 877, 866, 855, 844, 833, 822 - бесплатные
		if strings.HasPrefix(phone, "+1800") ||
			strings.HasPrefix(phone, "+1888") ||
			strings.HasPrefix(phone, "+1877") ||
			strings.HasPrefix(phone, "+1866") {
			return PhoneTypeTollFree
		}
		// 900 номера - премиум
		if strings.HasPrefix(phone, "+1900") {
			return PhoneTypePremium
		}
	case "uk":
		// В UK номера начинающиеся с +447 обычно мобильные
		if strings.HasPrefix(phone, "+447") {
			return PhoneTypeMobile
		}
		// +4480x - бесплатные
		if strings.HasPrefix(phone, "+44800") {
			return PhoneTypeTollFree
		}
	}

	// По умолчанию считаем мобильным, так как библиотека в основном для мобильных номеров
	return PhoneTypeMobile
}

// calculateDistance вычисляет расстояние Левенштейна между строками
func calculateDistance(s1, s2 string) int {
	len1, len2 := len(s1), len(s2)
	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Создаем матрицу
	matrix := make([][]int, len1+1)
	for i := range matrix {
		matrix[i] = make([]int, len2+1)
	}

	// Инициализируем первую строку и столбец
	for i := 0; i <= len1; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len2; j++ {
		matrix[0][j] = j
	}

	// Заполняем матрицу
	for i := 1; i <= len1; i++ {
		for j := 1; j <= len2; j++ {
			cost := 0
			if s1[i-1] != s2[j-1] {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // удаление
				matrix[i][j-1]+1,      // вставка
				matrix[i-1][j-1]+cost, // замена
			)
		}
	}

	return matrix[len1][len2]
}

// min возвращает минимальное из трех чисел
func min(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}
	if b <= c {
		return b
	}
	return c
}

// findSimilarCountries находит похожие коды стран
func findSimilarCountries(countryCode string) []string {
	var similar []string
	countryCode = strings.ToLower(countryCode)

	for code := range CountryPhoneCodes {
		distance := calculateDistance(countryCode, code)
		if distance <= 1 && distance > 0 { // Максимум 1 символ различия
			similar = append(similar, code)
		}
	}

	return similar
}

// validateRegexPattern проверяет, является ли строка валидным регулярным выражением
func validateRegexPattern(pattern string) bool {
	_, err := regexp.Compile(pattern)
	return err == nil
}

// parsePhoneComponents разбирает номер телефона на компоненты
func parsePhoneComponents(phone string) (prefix, localNumber string) {
	cleaned := cleanPhoneNumber(phone)

	if !strings.HasPrefix(cleaned, "+") {
		return "", cleaned
	}

	// Ищем подходящий префикс
	for _, info := range CountryPhoneCodes {
		if strings.HasPrefix(cleaned, info.Prefix) {
			prefix = info.Prefix
			localNumber = strings.TrimPrefix(cleaned, info.Prefix)
			return prefix, localNumber
		}
	}

	return "", cleaned
}

// isNumericOnly проверяет, содержит ли строка только цифры
func isNumericOnly(s string) bool {
	for _, char := range s {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return len(s) > 0
}

// removeNonDigits удаляет все символы кроме цифр
func removeNonDigits(s string) string {
	var result strings.Builder
	for _, char := range s {
		if unicode.IsDigit(char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// countryCodeExists проверяет существование кода страны
func countryCodeExists(countryCode string) bool {
	normalizedCode := normalizeCountryCode(countryCode)
	_, exists := CountryPhoneCodes[normalizedCode]
	return exists
}

// getPhoneLength возвращает длину номера без префикса
func getPhoneLength(phone, prefix string) int {
	cleaned := cleanPhoneNumber(phone)
	if strings.HasPrefix(cleaned, prefix) {
		localNumber := strings.TrimPrefix(cleaned, prefix)
		return len(removeNonDigits(localNumber))
	}
	return len(removeNonDigits(cleaned))
}

// generatePhoneExample генерирует пример номера для страны
func generatePhoneExample(countryCode string) string {
	info, exists := GetCountryInfo(countryCode)
	if !exists {
		return ""
	}

	// Генерируем простой пример
	example := info.Prefix
	targetLength := info.MinLength

	// Добавляем цифры для достижения минимальной длины
	for i := 0; i < targetLength; i++ {
		example += "0"
	}

	return example
}
