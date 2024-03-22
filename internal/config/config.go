package config

import (
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	AppMongoURI                      string = "APP_MONGO_URI"
	AppMongoDatabaseName             string = "APP_MONGO_DATABASE_NAME"
	AppMongoPoolMin                  string = "APP_MONGO_POOL_MIN"
	AppMongoPoolMax                  string = "APP_MONGO_POOL_MAX"
	AppMongoMaxIdleTimeSecond        string = "APP_MONGO_MAX_IDLE_TIME_SECOND"
	AppMongoInitConnectionTimeSecond string = "APP_MONGO_INIT_CONNECTION_TIME_SECOND"
	AppMongoQueryTimeoutMs           string = "APP_MONGO_QUERY_TIMEOUT_MS"

	APITimeout   string = "API_TIMEOUT"
	DefaultLimit string = "DEFAULT_LIMIT"
)

type Config struct {
	AppMongoURI                      string `validate:"required"`
	AppMongoDatabaseName             string `validate:"required"`
	AppMongoPoolMin                  int    `validate:"required"`
	AppMongoPoolMax                  int    `validate:"required"`
	AppMongoMaxIdleTimeSecond        int    `validate:"required"`
	AppMongoInitConnectionTimeSecond int    `validate:"required"`
	AppMongoQueryTimeoutMs           int    `validate:"required"`

	APITimeout   int `validate:"required"`
	DefaultLimit int `validate:"required"`
}

func New(validate *validator.Validate) Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("failed to load .env files, loading host env instead")
	}

	cfg := Config{
		AppMongoURI:                      os.Getenv(AppMongoURI),
		AppMongoDatabaseName:             os.Getenv(AppMongoDatabaseName),
		AppMongoPoolMin:                  getEnvInt(AppMongoPoolMin, os.Getenv(AppMongoPoolMin)),
		AppMongoPoolMax:                  getEnvInt(AppMongoPoolMax, os.Getenv(AppMongoPoolMax)),
		AppMongoMaxIdleTimeSecond:        getEnvInt(AppMongoMaxIdleTimeSecond, os.Getenv(AppMongoMaxIdleTimeSecond)),
		AppMongoInitConnectionTimeSecond: getEnvInt(AppMongoInitConnectionTimeSecond, os.Getenv(AppMongoInitConnectionTimeSecond)),
		AppMongoQueryTimeoutMs:           getEnvInt(AppMongoQueryTimeoutMs, os.Getenv(AppMongoQueryTimeoutMs)),

		APITimeout:   getEnvInt(APITimeout, os.Getenv(APITimeout)),
		DefaultLimit: getEnvInt(DefaultLimit, os.Getenv(DefaultLimit)),
	}

	if err := validate.Struct(cfg); err != nil {
		log.Panic().Err(err).Msg("failed to validate config")
	}

	return cfg
}

// convert env to int
func getEnvInt(env string, value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.
			Err(err).
			Stack().
			Str("env", env).
			Str("value", value).
			Msg("failed to convert string to int")
	}
	return i
}
