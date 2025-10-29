package configs

import (
	"os"
	"strings"

	"github.com/ShopOnGO/ShopOnGO/configs"
	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	Db           DbConfig
	Elastic      ElasticSearch
	Kafka        KafkaConfig
	LogLevel     logger.LogLevel
	FileLogLevel logger.LogLevel
}

type DbConfig struct {
	Dsn string
}

type ElasticSearch struct {
	URL   string
	Index string
}

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	GroupID  string
	ClientID string
}

func LoadConfig() *Config {
	if _, err := os.Stat(".env"); err == nil {
		// Локально есть .env → загружаем
		if loadErr := godotenv.Load(); loadErr != nil {
			logger.Error("Failed to load .env file", loadErr.Error())
		}
	} else {
		// В контейнере файла нет → просто идём дальше
		logger.Info(".env not found, using environment variables only")
	}

	brokersRaw := os.Getenv("KAFKA_SEARCH_BROKERS")
	brokers := strings.Split(brokersRaw, ",")

	// logger
	logLevelStr := os.Getenv("SEARCH_SERVICE_LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = "INFO"
	}
	LogLevel := configs.ParseLogLevel(logLevelStr)
	fileLogLevelStr := os.Getenv("SEARCH_SERVICE_FILE_LOG_LEVEL")
	if fileLogLevelStr == "" {
		fileLogLevelStr = "INFO"
	}
	FileLogLevel := configs.ParseLogLevel(fileLogLevelStr)

	return &Config{
		Db: DbConfig{
			Dsn: os.Getenv("DSN"),
		},
		Elastic: ElasticSearch{
			URL:   os.Getenv("ELASTIC_URL"),
			Index: os.Getenv("ELASTIC_INDEX"),
		},
		Kafka: KafkaConfig{
			Brokers:  brokers,
			Topic:    os.Getenv("KAFKA_SEARCH_TOPIC"),
			GroupID:  os.Getenv("KAFKA_SEARCH_GROUP_ID"),
			ClientID: os.Getenv("KAFKA_SEARCH_CLIENT_ID"),
		},
		LogLevel:     LogLevel,
		FileLogLevel: FileLogLevel,
	}
}
