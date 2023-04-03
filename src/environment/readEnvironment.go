package environment

import (
	_ "embed"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

const (
	envMongoDbConnectionString = "MONGODB_CONNECTION_STRING"
	envMongoDbDatabase         = "MONGODB_DATABASE_NAME"
	envAppExposePort           = "CAR_EXPOSE_PORT"
	envAppCollectionPrefix     = "CAR_COLLECTION_PREFIX"
	envLocalSetupMode          = "CAR_LOCAL_SETUP"

	defaultAppExposePort       = 80
	defaultAppCollectionPrefix = ""
)

func ptr[T any](v T) *T {
	return &v
}

//go:embed localSetup.env
var localSetup string

// readEnvironment reads the correct environment configuration (also considering local setup mode)
func readEnvironment() *Environment {
	if getBooleanEnvVariable(envLocalSetupMode) {
		fmt.Println("Using local setup mode.")

		localSetupMap, err := godotenv.Unmarshal(localSetup)
		if err != nil {
			panic("Invalid local setup environment variables. This is a bug.")
		}

		// Unfortunately, godotenv does not support reading environment variables from a string
		// directly to the environment. Therefore, we have to use this workaround.
		for key, value := range localSetupMap {
			if os.Getenv(key) != "" {
				// do not overwrite existing environment variables
				continue
			}
			_ = os.Setenv(key, value)
		}
	}
	return readEnvironmentFromEnv()
}

// readEnvironmentFromEnv reads the environment configuration from actual environment variables
// If any of the required environment variables is not set, the program will panic.
func readEnvironmentFromEnv() *Environment {
	return &Environment{
		mongoDbConnectionString: getStringEnvVariable(envMongoDbConnectionString, nil),
		mongoDbDatabase:         getStringEnvVariable(envMongoDbDatabase, nil),
		appExposePort:           getIntegerEnvVariable(envAppExposePort, ptr(defaultAppExposePort)),
		appCollectionPrefix:     getStringEnvVariable(envAppCollectionPrefix, ptr(defaultAppCollectionPrefix)),
		isLocalSetupMode:        getBooleanEnvVariable(envLocalSetupMode),
	}
}

// getStringEnvVariable returns the string value of the environment variable with the given name.
// You can specify a default value that is returned if the environment variable is not set,
// set defaultValue to nil to disable this feature.
// If defaultValue is nil and the environment variable is not set, the program will panic.
func getStringEnvVariable(variableName string, defaultValue *string) string {
	stringValue := os.Getenv(variableName)
	if stringValue != "" {
		return stringValue
	}

	if defaultValue != nil {
		return *defaultValue
	}
	panic("Environment variable not set: " + variableName)
}

// getIntegerEnvVariable returns the integer value of the environment variable with the given name.
// You can specify a default value that is returned if the environment variable is not set,
// set defaultValue to nil to disable this feature.
// If defaultValue is nil and the environment variable is not set, the program will panic.
// If the environment variable is not a valid integer value, the program will panic.
func getIntegerEnvVariable(variableName string, defaultValue *int) int {
	var stringValue string
	if defaultValue != nil {
		defaultValueString := strconv.Itoa(*defaultValue)
		stringValue = getStringEnvVariable(variableName, &defaultValueString)
	} else {
		stringValue = getStringEnvVariable(variableName, nil)
	}

	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		panic(fmt.Sprintf("Invalid value for integer environment variable \"%s\": %s",
			variableName, stringValue))
	}
	return intValue
}

// getBooleanEnvVariable returns the boolean value of the environment variable with the given name.
// If the environment variable is not set, false is returned.
// If the environment variable is not a valid boolean value, the program will panic.
func getBooleanEnvVariable(variableName string) bool {
	stringValue := os.Getenv(variableName)
	if stringValue == "" || stringValue == "false" {
		return false
	}

	if stringValue == "true" {
		return true
	}

	panic(fmt.Sprintf("Invalid value for boolean environment variable \"%s\": %s",
		variableName, stringValue))
}
