package translations

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// TranslationHelperFunc is a function that translates a key to its localized value
type TranslationHelperFunc func(key string, defaultValue string) string

// NullTranslationHelper returns the default value without translation
func NullTranslationHelper(_ string, defaultValue string) string {
	return defaultValue
}

// TranslationHelper returns a function that handles translations and a function to dump translations
func TranslationHelper() (TranslationHelperFunc, func()) {
	var translationKeyMap = map[string]string{}
	v := viper.New()

	v.SetEnvPrefix("KB_CLOUD_MCP")
	v.AutomaticEnv()

	// Load from JSON file
	v.SetConfigName("kb-cloud-mcp-server-config")
	v.SetConfigType("json")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		// ignore error if file not found as it is not required
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Could not read JSON config: %v", err)
		}
	}

	// Create a function that takes both a key, and a default value and returns either the default value or an override value
	return func(key string, defaultValue string) string {
			key = strings.ToUpper(key)
			if value, exists := translationKeyMap[key]; exists {
				return value
			}
			// Check if the env var exists
			if value, exists := os.LookupEnv("KB_CLOUD_MCP_" + key); exists {
				translationKeyMap[key] = value
				return value
			}

			v.SetDefault(key, defaultValue)
			translationKeyMap[key] = v.GetString(key)
			return translationKeyMap[key]
		}, func() {
			// Dump the translationKeyMap to a JSON file
			if err := DumpTranslationKeyMap(translationKeyMap); err != nil {
				log.Fatalf("Could not dump translation key map: %v", err)
			}
		}
}

// DumpTranslationKeyMap dumps the translation key map to a JSON file
func DumpTranslationKeyMap(translationKeyMap map[string]string) error {
	file, err := os.Create("kb-cloud-mcp-server-config.json")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer func() { _ = file.Close() }()

	// Marshal the map to JSON
	jsonData, err := json.MarshalIndent(translationKeyMap, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling map to JSON: %v", err)
	}

	// Write the JSON data to the file
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
