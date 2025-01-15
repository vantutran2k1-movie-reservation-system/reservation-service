package config

import (
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/utils"
	"log"
	"os"
	"strconv"
)

type Env struct {
	AppPort                         string
	GinMode                         string
	DbHost                          string
	DbPort                          string
	DbUser                          string
	DbPassword                      string
	DbName                          string
	RedisHost                       string
	RedisPort                       string
	RedisPassword                   string
	RedisDatabase                   int
	MinioHost                       string
	MinioPort                       string
	MinioConsolePort                int
	MinioRootUser                   string
	MinioRootPassword               string
	MinioAccessKey                  string
	MinioSecretKey                  string
	MinioProfilePictureBucket       string
	MaxProfilePictureFileSize       int
	ConfigcatSdkKey                 string
	LoginTokenExpireTime            int
	PassResetTokenExpireTime        int
	UserRegistrationTokenExpireTime int
	UserLocationApiTimeout          int
	UserLocationApiUrl              string
	KafkaBroker                     string
	KafkaUserRegistrationTopic      string
	PlatformUiEndpoint              string
}

var AppEnv Env

func InitAppEnv() {
	AppEnv.AppPort = getOrDefault("APP_PORT", "8080")
	AppEnv.GinMode = mustBeOneOf("GIN_MODE", utils.GetPointerOf(constants.GinDebugMode), constants.GinDebugMode, constants.GinReleaseMode)

	AppEnv.DbHost = getOrDefault("DB_HOST", "localhost")
	AppEnv.DbPort = getOrDefault("DB_PORT", "5432")
	AppEnv.DbUser = mustGetEnv("DB_USER")
	AppEnv.DbPassword = mustGetEnv("DB_PASSWORD")
	AppEnv.DbName = getOrDefault("DB_NAME", "booking")

	AppEnv.RedisHost = getOrDefault("REDIS_HOST", "localhost")
	AppEnv.RedisPort = getOrDefault("REDIS_PORT", "6379")
	AppEnv.RedisPassword = getOrDefault("REDIS_PASSWORD", "redis")
	AppEnv.RedisDatabase = getOrDefaultInt("REDIS_DATABASE", 0)

	AppEnv.MinioHost = getOrDefault("MINIO_HOST", "localhost")
	AppEnv.MinioPort = getOrDefault("MINIO_PORT", "9000")
	AppEnv.MinioRootUser = getOrDefault("MINIO_ROOT_USER", "minioadmin")
	AppEnv.MinioRootPassword = getOrDefault("MINIO_ROOT_PASSWORD", "minioadmin")
	AppEnv.MinioAccessKey = mustGetEnv("MINIO_ACCESS_KEY")
	AppEnv.MinioSecretKey = mustGetEnv("MINIO_SECRET_KEY")

	AppEnv.MinioProfilePictureBucket = getOrDefault("MINIO_PROFILE_PICTURE_BUCKET_NAME", "users.profile-pictures")
	AppEnv.MaxProfilePictureFileSize = getOrDefaultInt("MAX_USER_PROFILE_PICTURE_FILE_SIZE_MB", 10)

	AppEnv.ConfigcatSdkKey = mustGetEnv("CONFIGCAT_SDK_KEY")

	AppEnv.LoginTokenExpireTime = getOrDefaultInt("LOGIN_TOKEN_EXPIRES_AFTER_MINUTES", 60)
	AppEnv.PassResetTokenExpireTime = getOrDefaultInt("PASSWORD_RESET_TOKEN_EXPIRES_AFTER_MINUTES", 5)
	AppEnv.UserRegistrationTokenExpireTime = getOrDefaultInt("USER_REGISTRATION_TOKEN_EXPIRES_AFTER_MINUTES", 5)

	AppEnv.UserLocationApiTimeout = getOrDefaultInt("USER_LOCATION_API_TIMEOUT_SECONDS", 10)
	AppEnv.UserLocationApiUrl = getOrDefault("USER_LOCATION_API_URL", "http://ip-api.com/json/")

	AppEnv.KafkaBroker = getOrDefault("KAFKA_BROKER", "localhost:9092")
	AppEnv.KafkaUserRegistrationTopic = getOrDefault("KAFKA_USER_REGISTRATION_TOPIC", "users.user_registrations")

	AppEnv.PlatformUiEndpoint = getOrDefault("PLATFORM_UI_ENDPOINT", "http://localhost:5173")
}

func getOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getOrDefaultInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return valueInt
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("environment variable %s is not set", key)
	}

	return value
}

func mustGetEnvAsInt(key string) int {
	valueStr := mustGetEnv(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("environment variable %s must be an integer, got %s", key, valueStr)
	}

	return value
}

func mustBeOneOf(key string, defaultValues *string, acceptedValues ...string) string {
	var v string
	if defaultValues != nil {
		v = getOrDefault(key, *defaultValues)
	} else {
		v = mustGetEnv(key)
	}

	accepted := false
	for _, value := range acceptedValues {
		if v == value {
			accepted = true
			break
		}
	}

	if !accepted {
		log.Fatalf("environment variable %s must be one of %v", key, acceptedValues)
	}

	return v
}
