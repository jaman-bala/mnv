package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jaman-bala/mnv"
)

var (
	phone         = flag.String("phone", "", "Phone number to validate")
	country       = flag.String("country", "", "Country code (optional)")
	config        = flag.String("config", "default", "Configuration preset (default, strict, relaxed)")
	interactive   = flag.Bool("interactive", false, "Run in interactive mode")
	batch         = flag.String("batch", "", "File with phone numbers for batch validation")
	listCountries = flag.Bool("list-countries", false, "List all supported countries")
	info          = flag.Bool("info", false, "Show detailed phone information")
	suggestions   = flag.Bool("suggestions", false, "Show correction suggestions for invalid numbers")
	format        = flag.String("format", "text", "Output format: text, json")
	verbose       = flag.Bool("verbose", false, "Verbose output")
)

func main() {
	flag.Parse()

	// Применяем конфигурацию
	if err := mnv.SetPresetConfig(*config); err != nil {
		log.Printf("Warning: Invalid config preset '%s', using default\n", *config)
		mnv.SetPresetConfig("default")
	}

	if *verbose {
		fmt.Printf("Using configuration: %s\n", *config)
		printConfig()
		fmt.Println()
	}

	// Обработка команд
	switch {
	case *listCountries:
		listSupportedCountries()
	case *batch != "":
		processBatchFile(*batch)
	case *interactive:
		runInteractiveMode()
	case *phone != "":
		validateSinglePhone(*phone, *country)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Mobile Number Validator (MNV) CLI Tool")
	fmt.Println("Usage examples:")
	fmt.Println("  mnv -phone=\"+996700123456\" -country=\"kg\"")
	fmt.Println("  mnv -phone=\"+79991234567\" -info")
	fmt.Println("  mnv -list-countries")
	fmt.Println("  mnv -interactive")
	fmt.Println("  mnv -batch=\"phones.txt\"")
	fmt.Println()
	flag.PrintDefaults()
}

func printConfig() {
	config := mnv.GetConfig()
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Allow Spaces: %t\n", config.AllowSpaces)
	fmt.Printf("  Allow Dashes: %t\n", config.AllowDashes)
	fmt.Printf("  Allow Parentheses: %t\n", config.AllowParentheses)
	fmt.Printf("  Strict Mode: %t\n", config.StrictMode)
	fmt.Printf("  Require Plus Sign: %t\n", config.RequirePlusSign)
}

func listSupportedCountries() {
	countries := mnv.GetSupportedCountries()

	if *format == "json" {
		countriesInfo := make(map[string]mnv.PhoneCodeInfo)
		for _, code := range countries {
			if info, exists := mnv.GetCountryInfo(code); exists {
				countriesInfo[code] = info
			}
		}

		data, _ := json.MarshalIndent(map[string]interface{}{
			"countries": countries,
			"count":     len(countries),
			"details":   countriesInfo,
		}, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("Supported countries (%d total):\n\n", len(countries))
	fmt.Printf("%-4s %-20s %-10s %s\n", "Code", "Country", "Prefix", "Description")
	fmt.Println(strings.Repeat("-", 60))

	for _, code := range countries {
		if info, exists := mnv.GetCountryInfo(code); exists {
			fmt.Printf("%-4s %-20s %-10s %s\n",
				strings.ToUpper(code),
				info.CountryName,
				info.Prefix,
				info.Description)
		}
	}
}

func validateSinglePhone(phone, country string) {
	options := &mnv.ValidationOptions{
		ReturnSuggestions: *suggestions,
		MaxSuggestions:    5,
	}

	var result *mnv.ValidationResult

	if country != "" {
		result = mnv.ValidatePhone(phone, country, options)
	} else {
		// Пытаемся определить страну автоматически
		detectedCountry, found := mnv.GetCountryByPhone(phone)
		if found {
			result = mnv.ValidatePhone(phone, detectedCountry, options)
		} else {
			result = &mnv.ValidationResult{
				IsValid:        false,
				OriginalNumber: phone,
				ErrorMessage:   "Cannot determine country for phone number",
			}
		}
	}

	if *format == "json" {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Текстовый вывод
	fmt.Printf("Phone Number: %s\n", phone)
	if result.CountryCode != "" {
		fmt.Printf("Country: %s (%s)\n", result.CountryName, strings.ToUpper(result.CountryCode))
	}

	if result.IsValid {
		fmt.Printf("Status: ✅ VALID\n")
		if result.FormattedNumber != phone {
			fmt.Printf("Formatted: %s\n", result.FormattedNumber)
		}

		if *info {
			phoneInfo := mnv.GetPhoneInfo(phone)
			fmt.Printf("Type: %s\n", phoneInfo.Type)
			fmt.Printf("Prefix: %s\n", phoneInfo.Prefix)
			fmt.Printf("Local Number: %s\n", phoneInfo.LocalNumber)
		}
	} else {
		fmt.Printf("Status: ❌ INVALID\n")
		if result.ErrorMessage != "" {
			fmt.Printf("Error: %s\n", result.ErrorMessage)
		}

		if len(result.Suggestions) > 0 {
			fmt.Printf("Suggestions:\n")
			for _, suggestion := range result.Suggestions {
				fmt.Printf("  - %s\n", suggestion)
			}
		}
	}
}

func processBatchFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var phones []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			phones = append(phones, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	if len(phones) == 0 {
		fmt.Println("No phone numbers found in file")
		return
	}

	fmt.Printf("Processing %d phone numbers from %s...\n\n", len(phones), filename)

	request := &mnv.BatchValidationRequest{
		Phones:   phones,
		Parallel: 10,
		Options: &mnv.ValidationOptions{
			ReturnSuggestions: *suggestions,
			MaxSuggestions:    3,
		},
	}

	response := mnv.BatchValidatePhones(request)

	if *format == "json" {
		data, _ := json.MarshalIndent(response, "", "  ")
		fmt.Println(string(data))
		return
	}

	// Текстовый отчет
	fmt.Printf("Batch Validation Results:\n")
	fmt.Printf("========================\n")
	fmt.Printf("Total: %d | Valid: %d | Invalid: %d | Processing Time: %s\n\n",
		response.Stats.Total,
		response.Stats.Valid,
		response.Stats.Invalid,
		response.ProcessingTime)

	// Статистика по странам
	if len(response.Stats.ByCountry) > 0 {
		fmt.Printf("By Country:\n")
		for country, count := range response.Stats.ByCountry {
			if info, exists := mnv.GetCountryInfo(country); exists {
				fmt.Printf("  %s (%s): %d\n", info.CountryName, strings.ToUpper(country), count)
			} else {
				fmt.Printf("  %s: %d\n", strings.ToUpper(country), count)
			}
		}
		fmt.Println()
	}

	// Детальные результаты
	if *verbose {
		fmt.Printf("Detailed Results:\n")
		fmt.Printf("-----------------\n")
		for i, result := range response.Results {
			status := "❌"
			if result.IsValid {
				status = "✅"
			}
			fmt.Printf("%d. %s %s", i+1, status, result.OriginalNumber)

			if result.CountryCode != "" {
				fmt.Printf(" (%s)", strings.ToUpper(result.CountryCode))
			}

			if !result.IsValid && result.ErrorMessage != "" {
				fmt.Printf(" - %s", result.ErrorMessage)
			}

			fmt.Println()
		}
	}
}

func runInteractiveMode() {
	fmt.Println("Mobile Number Validator - Interactive Mode")
	fmt.Println("Commands:")
	fmt.Println("  validate <phone> [country] - Validate a phone number")
	fmt.Println("  info <phone>              - Get phone information")
	fmt.Println("  countries                 - List supported countries")
	fmt.Println("  config <preset>           - Change configuration")
	fmt.Println("  help                      - Show this help")
	fmt.Println("  quit                      - Exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("mnv> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		command := parts[0]

		switch command {
		case "validate":
			if len(parts) < 2 {
				fmt.Println("Usage: validate <phone> [country]")
				continue
			}

			phone := parts[1]
			country := ""
			if len(parts) > 2 {
				country = parts[2]
			}

			validateSinglePhone(phone, country)

		case "info":
			if len(parts) < 2 {
				fmt.Println("Usage: info <phone>")
				continue
			}

			phoneInfo := mnv.GetPhoneInfo(parts[1])
			data, _ := json.MarshalIndent(phoneInfo, "", "  ")
			fmt.Println(string(data))

		case "countries":
			listSupportedCountries()

		case "config":
			if len(parts) < 2 {
				fmt.Println("Available presets:", strings.Join(mnv.ListPresets(), ", "))
				fmt.Println("Current config:")
				printConfig()
				continue
			}

			preset := parts[1]
			if err := mnv.SetPresetConfig(preset); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Configuration changed to: %s\n", preset)
				printConfig()
			}

		case "help":
			fmt.Println("Commands:")
			fmt.Println("  validate <phone> [country] - Validate a phone number")
			fmt.Println("  info <phone>              - Get phone information")
			fmt.Println("  countries                 - List supported countries")
			fmt.Println("  config <preset>           - Change configuration")
			fmt.Println("  help                      - Show this help")
			fmt.Println("  quit                      - Exit")

		case "quit", "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}

		fmt.Println()
	}
}
