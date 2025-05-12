package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Key struct {
	YaMapsAPIKeyGeo     string
	YaMapsAPIKeyMap     string
	YaMapsAPIKeyLocator string
}

var apiKeys = initKeys()

func initKeys() Key {
	k := Key{}
	k.init()
	return k
}

func (k *Key) init() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Получаем значения из переменных окружения
	k.YaMapsAPIKeyGeo = getEnv("YA_MAPS_API_KEY_GEO", "")
	k.YaMapsAPIKeyMap = getEnv("YA_MAPS_API_KEY_MAP", "")
	k.YaMapsAPIKeyLocator = getEnv("YA_MAPS_API_KEY_LOCATOR", "")

	// Проверяем обязательные переменные
	if k.YaMapsAPIKeyGeo == "" || k.YaMapsAPIKeyMap == "" || k.YaMapsAPIKeyLocator == "" {
		log.Fatal("API keys not set in .env file")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			log.Printf("Warning: Environment variable %s not set", key)
		}
		return defaultValue
	}
	return value
}

func GetAPIGeo() string {
	return apiKeys.YaMapsAPIKeyGeo
}

func GetAPIMap() string {
	return apiKeys.YaMapsAPIKeyMap
}

