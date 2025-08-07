package mnv_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/jaman-bala/mnv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePhoneForCountry(t *testing.T) {
	// Устанавливаем конфигурацию по умолчанию
	mnv.SetConfig(mnv.DefaultConfig())

	tests := []struct {
		name        string
		phone       string
		countryCode string
		expected    bool
	}{
		// Кыргызстан
		{"Valid KG phone", "+996700123456", "kg", true},
		{"Valid KG phone 2", "+996555987654", "kg", true},
		{"Valid KG phone 3", "+996777123456", "kg", true},
		{"Invalid KG phone - short", "+99670012345", "kg", false},
		{"Invalid KG phone - long", "+99670012345678", "kg", false},
		{"Invalid KG phone - wrong prefix", "+79991234567", "kg", false},
		{"Invalid KG phone - no plus", "996700123456", "kg", false},

		// Россия
		{"Valid RU phone", "+79991234567", "ru", true},
		{"Valid RU phone 2", "+78001234567", "ru", true},
		{"Valid RU phone 3", "+79161234567", "ru", true},
		{"Invalid RU phone - short", "+7999123456", "ru", false},
		{"Invalid RU phone - long", "+799912345678", "ru", false},
		{"Invalid RU phone - Kazakhstan format", "+77001234567", "ru", false},

		// Казахстан
		{"Valid KZ phone", "+77001234567", "kz", true},
		{"Valid KZ phone 2", "+76001234567", "kz", true},
		{"Invalid KZ phone - Russia format", "+79991234567", "kz", false},

		// США
		{"Valid US phone", "+14155552671", "us", true},
		{"Valid US phone 2", "+12125551234", "us", true},
		{"Valid US phone 3", "+13105551234", "us", true},
		{"Invalid US phone - short", "+1415555267", "us", false},
		{"Invalid US phone - long", "+141555526712", "us", false},
		{"Invalid US phone - starts with 1", "+11555552671", "us", false},

		// Великобритания
		{"Valid UK phone", "+447911123456", "uk", true},
		{"Valid UK phone 2", "+447700123456", "uk", true},
		{"Invalid UK phone - short", "+44791112345", "uk", false},
		{"Invalid UK phone - wrong format", "+441234567890", "uk", false},

		// Узбекистан
		{"Valid UZ phone", "+998901234567", "uz", true},
		{"Valid UZ phone 2", "+998971234567", "uz", true},
		{"Invalid UZ phone - short", "+99890123456", "uz", false},
		{"Invalid UZ phone - long", "+9989012345678", "uz", false},

		// Турция
		{"Valid TR phone", "+905123456789", "tr", true},
		{"Valid TR phone 2", "+905321234567", "tr", true},
		{"Invalid TR phone - wrong prefix", "+901123456789", "tr", false},

		// Неподдерживаемая страна
		{"Unsupported country", "+123456789", "xx", false},
		{"Empty phone", "", "kg", false},
		{"Empty country", "+996700123456", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mnv.IsPhoneValid(tt.phone, tt.countryCode)
			assert.Equal(t, tt.expected, result,
				"Phone: %s, Country: %s", tt.phone, tt.countryCode)
		})
	}
}

func TestValidatePhoneWithDifferentConfigs(t *testing.T) {
	tests := []struct {
		name     string
		config   mnv.ValidatorConfig
		phone    string
		country  string
		expected bool
	}{
		{
			name:     "Strict mode - valid",
			config:   mnv.StrictConfig(),
			phone:    "+996700123456",
			country:  "kg",
			expected: true,
		},
		{
			name:     "Strict mode - no plus",
			config:   mnv.StrictConfig(),
			phone:    "996700123456",
			country:  "kg",
			expected: false,
		},
		{
			name:     "Relaxed mode - no plus",
			config:   mnv.RelaxedConfig(),
			phone:    "996700123456",
			country:  "kg",
			expected: true,
		},
		{
			name:     "Allow spaces",
			config:   mnv.RelaxedConfig(),
			phone:    "+996 700 123 456",
			country:  "kg",
			expected: true,
		},
		{
			name:     "Disallow spaces",
			config:   mnv.StrictConfig(),
			phone:    "+996 700 123 456",
			country:  "kg",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mnv.SetConfig(tt.config)
			result := mnv.IsPhoneValid(tt.phone, tt.country)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Возвращаем конфигурацию по умолчанию
	mnv.SetConfig(mnv.DefaultConfig())
}

func TestGetCountryByPhone(t *testing.T) {
	mnv.SetConfig(mnv.DefaultConfig())

	tests := []struct {
		name            string
		phone           string
		expectedCountry string
		expectedFound   bool
	}{
		{"KG phone", "+996700123456", "kg", true},
		{"RU phone", "+79991234567", "ru", true},
		{"US phone", "+14155552671", "us", true},
		{"UK phone", "+447911123456", "uk", true},
		{"UZ phone", "+998901234567", "uz", true},
		{"TR phone", "+905123456789", "tr", true},
		{"Invalid phone", "+123456789", "", false},
		{"Empty phone", "", "", false},
		{"Incomplete phone", "+996", "", false},
		{"Wrong format", "8996700123456", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			country, found := mnv.GetCountryByPhone(tt.phone)
			assert.Equal(t, tt.expectedCountry, country)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestFormatPhone(t *testing.T) {
	mnv.SetConfig(mnv.DefaultConfig())

	tests := []struct {
		name        string
		phone       string
		countryCode string
		expectError bool
		expected    string
	}{
		{"Valid KG phone", "+996700123456", "kg", false, "+996700123456"},
		{"Valid RU phone", "+79991234567", "ru", false, "+79991234567"},
		{"Valid US phone", "+14155552671", "us", false, "+14155552671"},
		{"Invalid phone", "+123456789", "kg", true, ""},
		{"Unsupported country", "+996700123456", "xx", true, ""},
		{"Empty phone", "", "kg", true, ""},
		{"Wrong country for phone", "+996700123456", "us", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mnv.FormatPhone(tt.phone, tt.countryCode)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAddAndRemoveCountry(t *testing.T) {
	// Сохраняем оригинальное состояние
	originalCodes := make(map[string]mnv.PhoneCodeInfo)
	for _, code := range mnv.GetSupportedCountries() {
		if info, exists := mnv.GetCountryInfo(code); exists {
			originalCodes[code] = info
		}
	}

	// Тест добавления новой страны
	t.Run("Add valid country", func(t *testing.T) {
		err := mnv.AddCountry("test", "+999", `^\+999[0-9]{7}$`, 7, 7)
		require.NoError(t, err)

		// Проверяем, что страна добавлена
		info, exists := mnv.GetCountryInfo("test")
		require.True(t, exists)
		assert.Equal(t, "+999", info.Prefix)
		assert.Equal(t, 7, info.MinLength)
		assert.Equal(t, 7, info.MaxLength)

		// Проверяем валидацию для новой страны
		assert.True(t, mnv.IsPhoneValid("+9991234567", "test"))
		assert.False(t, mnv.IsPhoneValid("+999123456", "test"))
	})

	// Тест добавления невалидной страны
	t.Run("Add invalid country", func(t *testing.T) {
		// Пустой код страны
		err := mnv.AddCountry("", "+888", `^\+888[0-9]{7}$`, 7, 7)
		assert.Error(t, err)

		// Пустой префикс
		err = mnv.AddCountry("test2", "", `^\+888[0-9]{7}$`, 7, 7)
		assert.Error(t, err)

		// Префикс без +
		err = mnv.AddCountry("test3", "888", `^\+888[0-9]{7}$`, 7, 7)
		assert.Error(t, err)

		// Неверные параметры длины
		err = mnv.AddCountry("test4", "+888", `^\+888[0-9]{7}$`, 0, 7)
		assert.Error(t, err)

		err = mnv.AddCountry("test5", "+888", `^\+888[0-9]{7}$`, 10, 5)
		assert.Error(t, err)
	})

	// Тест удаления страны
	t.Run("Remove country", func(t *testing.T) {
		// Удаляем тестовую страну
		mnv.RemoveCountry("test")
		_, exists := mnv.GetCountryInfo("test")
		assert.False(t, exists)
	})

	// Восстанавливаем оригинальное состояние
	// (В реальном приложении нужно было бы сохранить/восстановить состояние)
}

func TestValidatorWithGoPlaygroundValidator(t *testing.T) {
	// Создаем новый валидатор
	validate := validator.New()

	// Регистрируем наши кастомные валидаторы
	err := mnv.RegisterValidators(validate)
	require.NoError(t, err)

	// Тестовые структуры
	type UserKG struct {
		Name  string `validate:"required"`
		Phone string `validate:"required,kg"`
	}

	type UserWithCountry struct {
		Name    string `validate:"required"`
		Phone   string `validate:"required,phonebycountry"`
		Country string `validate:"required"`
	}

	type UserGeneral struct {
		Name  string `validate:"required"`
		Phone string `validate:"required,phone"`
	}

	// Тест валидации для Кыргызстана
	t.Run("KG validator", func(t *testing.T) {
		validUser := UserKG{Name: "Test", Phone: "+996700123456"}
		err := validate.Struct(validUser)
		assert.NoError(t, err)

		invalidUser := UserKG{Name: "Test", Phone: "+79991234567"}
		err = validate.Struct(invalidUser)
		assert.Error(t, err)
	})

	// Тест валидации по стране
	t.Run("Country-based validator", func(t *testing.T) {
		validUser := UserWithCountry{
			Name:    "Test",
			Phone:   "+996700123456",
			Country: "kg",
		}
		err := validate.Struct(validUser)
		assert.NoError(t, err)

		// Неверная страна для номера
		invalidUser := UserWithCountry{
			Name:    "Test",
			Phone:   "+996700123456",
			Country: "us",
		}
		err = validate.Struct(invalidUser)
		assert.Error(t, err)
	})

	// Тест общего валидатора
	t.Run("General validator", func(t *testing.T) {
		validUsers := []UserGeneral{
			{Name: "Test1", Phone: "+996700123456"}, // KG
			{Name: "Test2", Phone: "+79991234567"},  // RU
			{Name: "Test3", Phone: "+14155552671"},  // US
		}

		for _, user := range validUsers {
			err := validate.Struct(user)
			assert.NoError(t, err, "Should validate %s", user.Phone)
		}

		invalidUser := UserGeneral{Name: "Test", Phone: "+123456789"}
		err := validate.Struct(invalidUser)
		assert.Error(t, err)
	})
}

func TestValidatePhoneFunction(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		country string
		valid   bool
	}{
		{"Valid KG", "+996700123456", "kg", true},
		{"Valid RU", "+79991234567", "ru", true},
		{"Valid US", "+14155552671", "us", true},
		{"Invalid format", "invalid", "kg", false},
		{"Wrong country", "+996700123456", "us", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mnv.ValidatePhone(tt.phone, tt.country)
			assert.Equal(t, tt.valid, result.IsValid)
			assert.Equal(t, tt.phone, result.OriginalNumber)
			assert.Equal(t, tt.country, result.CountryCode)

			if result.IsValid {
				assert.NotEmpty(t, result.FormattedNumber)
				assert.NotEmpty(t, result.CountryName)
			} else {
				assert.NotEmpty(t, result.ErrorMessage)
			}
		})
	}
}

func TestBatchValidation(t *testing.T) {
	phones := []string{
		"+996700123456", // valid KG
		"+79991234567",  // valid RU
		"+14155552671",  // valid US
		"+123456789",    // invalid
		"invalid",       // invalid
	}

	request := &mnv.BatchValidationRequest{
		Phones:   phones,
		Parallel: 3,
		Options: &mnv.ValidationOptions{
			ReturnSuggestions: true,
			MaxSuggestions:    3,
		},
	}

	response := mnv.BatchValidatePhones(request)

	assert.Equal(t, len(phones), len(response.Results))
	assert.Equal(t, len(phones), response.Stats.Total)
	assert.Equal(t, 3, response.Stats.Valid)   // 3 валидных
	assert.Equal(t, 2, response.Stats.Invalid) // 2 невалидных
	assert.NotEmpty(t, response.ProcessingTime)
}

func TestGetPhoneInfo(t *testing.T) {
	tests := []struct {
		phone       string
		expectValid bool
		expectType  mnv.PhoneType
	}{
		{"+996700123456", true, mnv.PhoneTypeMobile},
		{"+79991234567", true, mnv.PhoneTypeMobile},
		{"+14155552671", true, mnv.PhoneTypeMobile},
		{"+123456789", false, mnv.PhoneTypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.phone, func(t *testing.T) {
			info := mnv.GetPhoneInfo(tt.phone)
			assert.Equal(t, tt.phone, info.Number)
			assert.Equal(t, tt.expectValid, info.IsValid)
			assert.Equal(t, tt.expectType, info.Type)

			if info.IsValid {
				assert.NotEmpty(t, info.CountryCode)
				assert.NotEmpty(t, info.CountryName)
				assert.NotEmpty(t, info.Prefix)
			}
		})
	}
}

func TestConfigPresets(t *testing.T) {
	// Тест получения предустановок
	presets := mnv.ListPresets()
	assert.Contains(t, presets, "default")
	assert.Contains(t, presets, "strict")
	assert.Contains(t, presets, "relaxed")

	// Тест применения предустановки
	err := mnv.SetPresetConfig("strict")
	assert.NoError(t, err)

	config := mnv.GetConfig()
	assert.True(t, config.StrictMode)
	assert.True(t, config.RequirePlusSign)

	// Тест несуществующей предустановки
	err = mnv.SetPresetConfig("nonexistent")
	assert.Error(t, err)

	// Возвращаем к умолчанию
	mnv.SetPresetConfig("default")
}

// Бенчмарки
func BenchmarkValidatePhone(b *testing.B) {
	mnv.SetConfig(mnv.DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mnv.IsPhoneValid("+996700123456", "kg")
	}
}

func BenchmarkValidatePhoneStrict(b *testing.B) {
	mnv.SetConfig(mnv.StrictConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mnv.IsPhoneValid("+996700123456", "kg")
	}
}

func BenchmarkGetCountryByPhone(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mnv.GetCountryByPhone("+996700123456")
	}
}

func BenchmarkBatchValidation(b *testing.B) {
	phones := []string{
		"+996700123456",
		"+79991234567",
		"+14155552671",
		"+447911123456",
		"+998901234567",
	}

	request := &mnv.BatchValidationRequest{
		Phones:   phones,
		Parallel: 3,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mnv.BatchValidatePhones(request)
	}
}
