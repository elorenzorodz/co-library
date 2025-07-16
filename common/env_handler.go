package common

import (
	"log"
	"os"
)

func GetEnvVariable(name string) string {
	envValue := os.Getenv(name)

	if envValue == "" {
		log.Fatal(name, " is not found in the environment")
	}

	return envValue
}