# Mobile Number Validator (MNV)

Библиотека для валидации номеров мобильных телефонов для различных стран с поддержкой кастомных валидаторов для `go-playground/validator`.

## 🚀 Особенности

- ✅ Поддержка 19+ стран
- ✅ Интеграция с `gin-gonic/gin` и `go-playground/validator`
- ✅ Гибкая конфигурация валидации
- ✅ Строгий и базовый режимы валидации
- ✅ Определение страны по номеру телефона
- ✅ Форматирование номеров
- ✅ Возможность добавления новых стран
- ✅ Thread-safe операции

## 📦 Установка

```bash
go get github.com/jaman-bala/mnv
```

## 🌍 Поддерживаемые страны

| Страна | Код | Префикс | Пример |
|--------|-----|---------|--------|
| Кыргызстан | kg | +996 | +996700123456 |
| Россия | ru | +7 | +79991234567 |
| Казахстан | kz | +7 | +77771234567 |
| Узбекистан | uz | +998 | +998901234567 |
| США | us | +1 | +14155552671 |
| Великобритания | uk | +44 | +441234567890 |
| Германия | de | +49 | +4915123456789 |
| Франция | fr | +33 | +33612345678 |
| Италия | it | +39 | +393123456789 |
| Турция | tr | +90 | +905123456789 |
| И другие... | | | |

## 🔧 Быстрый старт

### 1. Базовое использование с Gin

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
    "github.com/jaman-bala/mnv"
)

type UserRequestDTO struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required,kg"` // Валидация для Кыргызстана
}
```

### 2. Динамическая валидация по стране

```go
type User struct {
    PhoneNumber string `json:"phone" binding:"required,phonebycountry"`
    Country     string `json:"country" binding:"required"` // Обязательно для phonebycountry
}
```

### 3. Программная валидация

```go
// Валидация для конкретной страны
isValid := mnv.ValidatePhoneForCountry("+996700123456", "kg")

// Определение страны по номеру
country, found := mnv.GetCountryByPhone("+996700123456")

// Форматирование номера
formatted, err := mnv.FormatPhone("+996 700 123 456", "kg")
```

## 🏷️ Доступные валидаторы

| Тег | Описание | Пример |
|-----|----------|--------|
| `kg` | Кыргызстан | `binding:"kg"` |
| `ru` | Россия | `binding:"ru"` |
| `us` | США | `binding:"us"` |
| `phone` | Любая страна | `binding:"phone"` |
| `phonebycountry` | По полю Country | `binding:"phonebycountry"` |

## ⚙️ Конфигурация

```go
// Настройка валидатора
mnv.SetConfig(mnv.ValidatorConfig{
    AllowSpaces:      true,  // Разрешить пробелы: "+996 700 123 456"
    AllowDashes:      false, // Запретить тире: "+996-700-123-456"
    AllowParentheses: false, // Запретить скобки: "+996(700)123456"
    StrictMode:       true,  // Строгая валидация с regex
})
```

## 🔍 Примеры валидации

### Правильные номера
```go
✅ "+996700123456"  // Кыргызстан
✅ "+79991234567"   // Россия
✅ "+14155552671"   // США
✅ "+441234567890"  // Великобритания
```

### Неправильные номера
```go
❌ "996700123456"   // Без +
❌ "+99670012345"   // Короткий номер
❌ "+996abc123456"  // Содержит буквы
❌ "+99670012345678" // Слишком длинный
```

## 🛠️ Расширенное использование

### Добавление новой страны

```go
mnv.AddCountry("md", "+373", `^\+373[0-9]{8}$`, 8, 8) // Молдова
```

### Получение списка стран

```go
countries := mnv.GetSupportedCountries()
fmt.Println(countries) // [kg, ru, us, uk, ...]
```

### Удаление страны

```go
mnv.RemoveCountry("md")
```

## 🌐 HTTP API пример

```go
// POST /validate
{
    "phone": "+996700123456",
    "country": "kg"  // опционально
}

// Response
{
    "phone": "+996700123456",
    "valid": true,
    "detected_country": "kg"
}
```

## 🧪 Тестирование

```bash
go test ./...
```

Примеры тестов:

```go
func TestKGPhone(t *testing.T) {
    tests := []struct {
        phone string
        valid bool
    }{
        {"+996700123456", true},
        {"+996555123456", true},
        {"+99670012345", false},   // короткий
        {"+7991234567", false},    // не кыргызский код
    }
    
    for _, tt := range tests {
        result := mnv.ValidatePhoneForCountry(tt.phone, "kg")
        assert.Equal(t, tt.valid, result)
    }
}
```

## 📝 Структуры данных

### PhoneCodeInfo
```go
type PhoneCodeInfo struct {
    Prefix    string `json:"prefix"`    // "+996"
    Pattern   string `json:"pattern"`   // Регулярное выражение
    MinLength int    `json:"min_length"`// Минимальная длина без префикса
    MaxLength int    `json:"max_length"`// Максимальная длина без префикса
}
```

### ValidatorConfig
```go
type ValidatorConfig struct {
    AllowSpaces      bool // Разрешить пробелы
    AllowDashes      bool // Разрешить тире
    AllowParentheses bool // Разрешить скобки
    StrictMode       bool // Строгий режим с regex
}
```

## 🤝 Вклад в проект

1. Fork репозиторий
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Сделайте коммит (`git commit -m 'Add amazing feature'`)
4. Push в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📄 Лицензия

Этот проект распространяется под лицензией MIT. См. файл `LICENSE` для подробностей.

## 🐛 Сообщить о проблеме

Если вы нашли баг или у вас есть предложения, пожалуйста, создайте [issue](https://github.com/jaman-bala/mnv/issues).

## 📞 Поддержка

- 📧 Email: ermekoffdastan@gmail.com
- 💬 Telegram: @dermekoff
- 🐛 Issues: [GitHub Issues](https://github.com/jaman-bala/mnv/issues)

---

Сделано с ❤️ для мирового сообщества разработчиков