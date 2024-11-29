package main

import (
	"fmt"
	"os"

	"cloudprojects/current-converter/api"
	"cloudprojects/current-converter/conversion"
	"cloudprojects/current-converter/tui"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (optional)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	// Get the API key
	apiKey := os.Getenv("OXR_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: API key not found in environment variable OXR_API_KEY")
		return
	}

	// Fetch currency rates
	rates, err := api.FetchRates(apiKey)
	if err != nil {
		fmt.Println("Error fetching rates:", err)
		return
	}

	// Run TUI
	conversionParams, err := tui.RunTUI()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Perform conversion
	rateFrom, ok := rates.Rates[conversionParams.CurrencyFrom]
	if !ok {
		fmt.Printf("Error: Unsupported currency %s\n", conversionParams.CurrencyFrom)
		return
	}

	rateTo, ok := rates.Rates[conversionParams.CurrencyTo]
	if !ok {
		fmt.Printf("Error: Unsupported currency %s\n", conversionParams.CurrencyTo)
		return
	}

	convertedValue := conversion.Convert(
		conversionParams.Amount,
		rateFrom,
		rateTo,
	)

	// Display the result
	fmt.Printf("%.2f %s = %.2f %s\n",
		conversionParams.Amount,
		conversionParams.CurrencyFrom,
		convertedValue,
		conversionParams.CurrencyTo,
	)
}
