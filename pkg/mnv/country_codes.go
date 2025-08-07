package mnv

// CountryPhoneCodes содержит коды стран и их телефонные префиксы
var CountryPhoneCodes = map[string]PhoneCodeInfo{
	// Центральная Азия
	"kg": {
		Prefix:      "+996",
		Pattern:     `^\+996[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Kyrgyzstan",
		Description: "Kyrgyzstan mobile numbers",
	},
	"kz": {
		Prefix:      "+7",
		Pattern:     `^\+7[67][0-9]{8}$`, // Казахстан использует +7 6xx и +7 7xx
		MinLength:   10,
		MaxLength:   10,
		CountryName: "Kazakhstan",
		Description: "Kazakhstan mobile numbers (6xx, 7xx)",
	},
	"uz": {
		Prefix:      "+998",
		Pattern:     `^\+998[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Uzbekistan",
		Description: "Uzbekistan mobile numbers",
	},
	"tj": {
		Prefix:      "+992",
		Pattern:     `^\+992[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Tajikistan",
		Description: "Tajikistan mobile numbers",
	},
	"tm": {
		Prefix:      "+993",
		Pattern:     `^\+993[0-9]{8}$`,
		MinLength:   8,
		MaxLength:   8,
		CountryName: "Turkmenistan",
		Description: "Turkmenistan mobile numbers",
	},

	// Россия и СНГ
	"ru": {
		Prefix:      "+7",
		Pattern:     `^\+7[9][0-9]{9}$`, // Россия использует +7 9xx
		MinLength:   10,
		MaxLength:   10,
		CountryName: "Russia",
		Description: "Russia mobile numbers (9xx)",
	},
	"ua": {
		Prefix:      "+380",
		Pattern:     `^\+380[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Ukraine",
		Description: "Ukraine mobile numbers",
	},
	"by": {
		Prefix:      "+375",
		Pattern:     `^\+375[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Belarus",
		Description: "Belarus mobile numbers",
	},
	"am": {
		Prefix:      "+374",
		Pattern:     `^\+374[0-9]{8}$`,
		MinLength:   8,
		MaxLength:   8,
		CountryName: "Armenia",
		Description: "Armenia mobile numbers",
	},
	"az": {
		Prefix:      "+994",
		Pattern:     `^\+994[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Azerbaijan",
		Description: "Azerbaijan mobile numbers",
	},
	"ge": {
		Prefix:      "+995",
		Pattern:     `^\+995[0-9]{9}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Georgia",
		Description: "Georgia mobile numbers",
	},
	"md": {
		Prefix:      "+373",
		Pattern:     `^\+373[0-9]{8}$`,
		MinLength:   8,
		MaxLength:   8,
		CountryName: "Moldova",
		Description: "Moldova mobile numbers",
	},

	// Западная Европа
	"de": {
		Prefix:      "+49",
		Pattern:     `^\+49[1][5-7][0-9]{8,9}$`,
		MinLength:   10,
		MaxLength:   12,
		CountryName: "Germany",
		Description: "Germany mobile numbers",
	},
	"fr": {
		Prefix:      "+33",
		Pattern:     `^\+33[6-7][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   10,
		CountryName: "France",
		Description: "France mobile numbers",
	},
	"uk": {
		Prefix:      "+44",
		Pattern:     `^\+44[7][0-9]{9}$`,
		MinLength:   10,
		MaxLength:   11,
		CountryName: "United Kingdom",
		Description: "UK mobile numbers",
	},
	"it": {
		Prefix:      "+39",
		Pattern:     `^\+39[3][0-9]{8,9}$`,
		MinLength:   9,
		MaxLength:   11,
		CountryName: "Italy",
		Description: "Italy mobile numbers",
	},
	"es": {
		Prefix:      "+34",
		Pattern:     `^\+34[6-7][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Spain",
		Description: "Spain mobile numbers",
	},
	"nl": {
		Prefix:      "+31",
		Pattern:     `^\+31[6][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Netherlands",
		Description: "Netherlands mobile numbers",
	},

	// Северная Америка
	"us": {
		Prefix:      "+1",
		Pattern:     `^\+1[2-9][0-9]{9}$`,
		MinLength:   10,
		MaxLength:   10,
		CountryName: "United States",
		Description: "US mobile numbers",
	},
	"ca": {
		Prefix:      "+1",
		Pattern:     `^\+1[2-9][0-9]{9}$`,
		MinLength:   10,
		MaxLength:   10,
		CountryName: "Canada",
		Description: "Canada mobile numbers",
	},

	// Азия
	"tr": {
		Prefix:      "+90",
		Pattern:     `^\+90[5][0-9]{9}$`,
		MinLength:   10,
		MaxLength:   10,
		CountryName: "Turkey",
		Description: "Turkey mobile numbers",
	},
	"cn": {
		Prefix:      "+86",
		Pattern:     `^\+86[1][3-9][0-9]{9}$`,
		MinLength:   11,
		MaxLength:   11,
		CountryName: "China",
		Description: "China mobile numbers",
	},
	"in": {
		Prefix:      "+91",
		Pattern:     `^\+91[6-9][0-9]{9}$`,
		MinLength:   10,
		MaxLength:   10,
		CountryName: "India",
		Description: "India mobile numbers",
	},
	"jp": {
		Prefix:      "+81",
		Pattern:     `^\+81[7-9][0-9]{8}$`,
		MinLength:   10,
		MaxLength:   11,
		CountryName: "Japan",
		Description: "Japan mobile numbers",
	},
	"kr": {
		Prefix:      "+82",
		Pattern:     `^\+82[1][0-9]{8,9}$`,
		MinLength:   9,
		MaxLength:   10,
		CountryName: "South Korea",
		Description: "South Korea mobile numbers",
	},

	// Ближний Восток
	"ae": {
		Prefix:      "+971",
		Pattern:     `^\+971[5][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "United Arab Emirates",
		Description: "UAE mobile numbers",
	},
	"sa": {
		Prefix:      "+966",
		Pattern:     `^\+966[5][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Saudi Arabia",
		Description: "Saudi Arabia mobile numbers",
	},
	"il": {
		Prefix:      "+972",
		Pattern:     `^\+972[5][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Israel",
		Description: "Israel mobile numbers",
	},

	// Африка
	"za": {
		Prefix:      "+27",
		Pattern:     `^\+27[6-8][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "South Africa",
		Description: "South Africa mobile numbers",
	},
	"eg": {
		Prefix:      "+20",
		Pattern:     `^\+20[1][0-9]{9}$`,
		MinLength:   10,
		MaxLength:   10,
		CountryName: "Egypt",
		Description: "Egypt mobile numbers",
	},

	// Океания
	"au": {
		Prefix:      "+61",
		Pattern:     `^\+61[4][0-9]{8}$`,
		MinLength:   9,
		MaxLength:   9,
		CountryName: "Australia",
		Description: "Australia mobile numbers",
	},
	"nz": {
		Prefix:      "+64",
		Pattern:     `^\+64[2][0-9]{7,9}$`,
		MinLength:   8,
		MaxLength:   10,
		CountryName: "New Zealand",
		Description: "New Zealand mobile numbers",
	},

	// Латинская Америка
	"br": {
		Prefix:      "+55",
		Pattern:     `^\+55[1-9][1-9][9][0-9]{8}$`,
		MinLength:   11,
		MaxLength:   11,
		CountryName: "Brazil",
		Description: "Brazil mobile numbers",
	},
	"ar": {
		Prefix:      "+54",
		Pattern:     `^\+54[9][1-9][0-9]{8}$`,
		MinLength:   10,
		MaxLength:   10,
		CountryName: "Argentina",
		Description: "Argentina mobile numbers",
	},
	"mx": {
		Prefix:      "+52",
		Pattern:     `^\+52[1][0-9]{10}$`,
		MinLength:   11,
		MaxLength:   11,
		CountryName: "Mexico",
		Description: "Mexico mobile numbers",
	},
}

// GetCountryInfo возвращает информацию о стране по коду
func GetCountryInfo(countryCode string) (PhoneCodeInfo, bool) {
	info, exists := CountryPhoneCodes[countryCode]
	return info, exists
}

// GetCountriesByPrefix возвращает страны по префиксу (например, +7 для России и Казахстана)
func GetCountriesByPrefix(prefix string) []string {
	var countries []string
	for code, info := range CountryPhoneCodes {
		if info.Prefix == prefix {
			countries = append(countries, code)
		}
	}
	return countries
}

// GetAllPrefixes возвращает все уникальные префиксы
func GetAllPrefixes() []string {
	prefixSet := make(map[string]bool)
	for _, info := range CountryPhoneCodes {
		prefixSet[info.Prefix] = true
	}

	var prefixes []string
	for prefix := range prefixSet {
		prefixes = append(prefixes, prefix)
	}
	return prefixes
}

// IsValidCountryCode проверяет, поддерживается ли код страны
func IsValidCountryCode(countryCode string) bool {
	_, exists := CountryPhoneCodes[countryCode]
	return exists
}
