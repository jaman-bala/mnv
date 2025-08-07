package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jaman-bala/mnv"
)

func main() {
	// Инициализация валидатора
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := mnv.RegisterValidators(v); err != nil {
			log.Fatal("Failed to register validators:", err)
		}
	}

	// Настройка конфигурации валидатора
	mnv.SetConfig(mnv.ValidatorConfig{
		AllowSpaces:      true,
		AllowDashes:      false,
		AllowParentheses: false,
		AllowDots:        false,
		StrictMode:       true,
		RequirePlusSign:  true,
	})

	// Инициализация Gin
	r := gin.Default()

	// Middleware для CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Группа API v1
	v1 := r.Group("/api/v1")
	{
		// Регистрация пользователя
		v1.POST("/register", handleRegister)
		v1.POST("/register-with-country", handleRegisterWithCountry)
		v1.POST("/profile", handleProfile)

		// Валидация номеров
		v1.POST("/validate", handleValidatePhone)
		v1.POST("/validate/batch", handleBatchValidation)
		v1.GET("/phone/:phone/info", handlePhoneInfo)

		// Информация о странах
		v1.GET("/countries", handleGetCountries)
		v1.GET("/countries/:code", handleGetCountryInfo)
		v1.POST("/countries", handleAddCountry)
		v1.DELETE("/countries/:code", handleRemoveCountry)

		// Конфигурация
		v1.GET("/config", handleGetConfig)
		v1.PUT("/config", handleSetConfig)
		v1.GET("/config/presets", handleGetPresets)
		v1.PUT("/config/preset/:preset", handleSetPreset)

		// Утилиты
		v1.GET("/format/:phone/:country", handleFormatPhone)
		v1.GET("/detect/:phone", handleDetectCountry)
	}

	// Статические файлы для документации
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Главная страница с документацией
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Mobile Number Validator API",
		})
	})

	// Запуск сервера
	log.Println("Server starting on :8080")
	log.Println("API documentation: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// RegisterRequestDTO базовая структура для регистрации
type RegisterRequestDTO struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone" binding:"required,kg"` // Валидация для Кыргызстана
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
}

// RegisterWithCountryDTO структура с динамической валидацией
type RegisterWithCountryDTO struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	PhoneNumber string `json:"phone" binding:"required,phonebycountry"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	Country     string `json:"country" binding:"required,min=2,max=2"`
}

// UserProfileDTO структура профиля пользователя
type UserProfileDTO struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone" binding:"required,phone"`
	Bio         string `json:"bio"`
}

// ValidationRequestDTO запрос на валидацию
type ValidationRequestDTO struct {
	Phone             string `json:"phone" binding:"required"`
	Country           string `json:"country"`
	ReturnInfo        bool   `json:"return_info"`
	ReturnSuggestions bool   `json:"return_suggestions"`
}

// ConfigUpdateDTO обновление конфигурации
type ConfigUpdateDTO struct {
	AllowSpaces      *bool `json:"allow_spaces"`
	AllowDashes      *bool `json:"allow_dashes"`
	AllowParentheses *bool `json:"allow_parentheses"`
	AllowDots        *bool `json:"allow_dots"`
	StrictMode       *bool `json:"strict_mode"`
	RequirePlusSign  *bool `json:"require_plus_sign"`
}

// Обработчики запросов

func handleRegister(c *gin.Context) {
	var req RegisterRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Дополнительная проверка и форматирование
	formatted, err := mnv.FormatPhone(req.PhoneNumber, "kg")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Phone formatting failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Registration successful",
		"formatted_phone": formatted,
		"data":            req,
	})
}

func handleRegisterWithCountry(c *gin.Context) {
	var req RegisterWithCountryDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Форматирование номера
	formatted, err := mnv.FormatPhone(req.PhoneNumber, req.Country)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Phone formatting failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Registration successful",
		"formatted_phone": formatted,
		"data":            req,
	})
}

func handleProfile(c *gin.Context) {
	var req UserProfileDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Определяем страну по номеру телефона
	country, found := mnv.GetCountryByPhone(req.PhoneNumber)
	if !found {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot determine country by phone number",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Profile updated successfully",
		"detected_country": country,
		"data":             req,
	})
}

func handleValidatePhone(c *gin.Context) {
	var req ValidationRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Создаем опции валидации
	options := &mnv.ValidationOptions{
		ReturnSuggestions: req.ReturnSuggestions,
		MaxSuggestions:    5,
	}

	var result *mnv.ValidationResult

	if req.Country != "" {
		// Валидация для конкретной страны
		result = mnv.ValidatePhone(req.Phone, req.Country, options)
	} else {
		// Общая валидация - пытаемся определить страну
		country, found := mnv.GetCountryByPhone(req.Phone)
		if found {
			result = mnv.ValidatePhone(req.Phone, country, options)
		} else {
			result = &mnv.ValidationResult{
				IsValid:        false,
				OriginalNumber: req.Phone,
				ErrorMessage:   "Cannot determine country for phone number",
			}
		}
	}

	response := gin.H{
		"phone":   req.Phone,
		"valid":   result.IsValid,
		"country": result.CountryCode,
	}

	if result.FormattedNumber != "" {
		response["formatted"] = result.FormattedNumber
	}

	if result.ErrorMessage != "" {
		response["error"] = result.ErrorMessage
	}

	if len(result.Suggestions) > 0 {
		response["suggestions"] = result.Suggestions
	}

	if req.ReturnInfo && result.IsValid {
		phoneInfo := mnv.GetPhoneInfo(req.Phone)
		response["info"] = phoneInfo
	}

	c.JSON(http.StatusOK, response)
}

func handleBatchValidation(c *gin.Context) {
	var req mnv.BatchValidationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	// Ограничиваем количество номеров в одном запросе
	if len(req.Phones) > 1000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Too many phone numbers. Maximum 1000 per request",
		})
		return
	}

	response := mnv.BatchValidatePhones(&req)
	c.JSON(http.StatusOK, response)
}

func handlePhoneInfo(c *gin.Context) {
	phone := c.Param("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone number is required",
		})
		return
	}

	phoneInfo := mnv.GetPhoneInfo(phone)
	c.JSON(http.StatusOK, phoneInfo)
}

func handleGetCountries(c *gin.Context) {
	countries := mnv.GetSupportedCountries()

	// Получаем детальную информацию для каждой страны
	countriesInfo := make(map[string]mnv.PhoneCodeInfo)
	for _, code := range countries {
		if info, exists := mnv.GetCountryInfo(code); exists {
			countriesInfo[code] = info
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"countries": countries,
		"count":     len(countries),
		"details":   countriesInfo,
	})
}

func handleGetCountryInfo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Country code is required",
		})
		return
	}

	info, exists := mnv.GetCountryInfo(code)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Country not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"country_code": code,
		"info":         info,
	})
}

func handleAddCountry(c *gin.Context) {
	var country mnv.CustomCountry

	if err := c.ShouldBindJSON(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid country data",
			"details": err.Error(),
		})
		return
	}

	if err := mnv.AddCustomCountry(&country); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to add country",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Country added successfully",
		"country": country,
	})
}

func handleRemoveCountry(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Country code is required",
		})
		return
	}

	// Проверяем, существует ли страна
	if _, exists := mnv.GetCountryInfo(code); !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Country not found",
		})
		return
	}

	mnv.RemoveCountry(code)

	c.JSON(http.StatusOK, gin.H{
		"message": "Country removed successfully",
		"code":    code,
	})
}

func handleGetConfig(c *gin.Context) {
	config := mnv.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"config": config,
	})
}

func handleSetConfig(c *gin.Context) {
	var updateReq ConfigUpdateDTO

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid configuration",
			"details": err.Error(),
		})
		return
	}

	// Получаем текущую конфигурацию
	config := mnv.GetConfig()

	// Обновляем только переданные поля
	if updateReq.AllowSpaces != nil {
		config.AllowSpaces = *updateReq.AllowSpaces
	}
	if updateReq.AllowDashes != nil {
		config.AllowDashes = *updateReq.AllowDashes
	}
	if updateReq.AllowParentheses != nil {
		config.AllowParentheses = *updateReq.AllowParentheses
	}
	if updateReq.AllowDots != nil {
		config.AllowDots = *updateReq.AllowDots
	}
	if updateReq.StrictMode != nil {
		config.StrictMode = *updateReq.StrictMode
	}
	if updateReq.RequirePlusSign != nil {
		config.RequirePlusSign = *updateReq.RequirePlusSign
	}

	mnv.SetConfig(config)

	c.JSON(http.StatusOK, gin.H{
		"message": "Configuration updated successfully",
		"config":  config,
	})
}

func handleGetPresets(c *gin.Context) {
	presets := mnv.ListPresets()
	c.JSON(http.StatusOK, gin.H{
		"presets": presets,
	})
}

func handleSetPreset(c *gin.Context) {
	preset := c.Param("preset")
	if preset == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Preset name is required",
		})
		return
	}

	if err := mnv.SetPresetConfig(preset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to set preset",
			"details": err.Error(),
		})
		return
	}

	config := mnv.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"message": "Preset applied successfully",
		"preset":  preset,
		"config":  config,
	})
}

func handleFormatPhone(c *gin.Context) {
	phone := c.Param("phone")
	country := c.Param("country")

	if phone == "" || country == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone number and country are required",
		})
		return
	}

	formatted, err := mnv.FormatPhone(phone, country)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Formatting failed",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original":  phone,
		"country":   country,
		"formatted": formatted,
	})
}

func handleDetectCountry(c *gin.Context) {
	phone := c.Param("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone number is required",
		})
		return
	}

	country, found := mnv.GetCountryByPhone(phone)

	response := gin.H{
		"phone": phone,
		"found": found,
	}

	if found {
		if info, exists := mnv.GetCountryInfo(country); exists {
			response["country"] = country
			response["country_name"] = info.CountryName
			response["prefix"] = info.Prefix
		}
	}

	c.JSON(http.StatusOK, response)
}
