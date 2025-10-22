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

func LoadEnvConfig() EnvConfig {
	return EnvConfig{
		APIVersion:           GetEnvVariable("API_VERSION"),
		Port:                 GetEnvVariable("PORT"),
		DBUrl:                GetEnvVariable("DB_URL"),
		MailgunAPIKey:        GetEnvVariable("MAILGUN_API_KEY"),
		MailgunSendingDomain: GetEnvVariable("MAILGUN_SENDING_DOMAIN"),
	}
}