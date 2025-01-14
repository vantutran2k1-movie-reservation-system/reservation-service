package config

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"

	configcat "github.com/configcat/go-sdk/v9"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisClient *redis.Client
var MinioClient *minio.Client
var ConfigcatClient *configcat.Client
var KafkaProducerClient sarama.SyncProducer

func InitDB() {
	l := logger.Default.LogMode(logger.Silent)
	if os.Getenv("GIN_MODE") == constants.GinDebugMode {
		l = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	DB = db
}

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDatabase, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		log.Fatalf("Redis database must be an integer")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDatabase,
	})

	RedisClient = rdb
}

func InitMinio() {
	minioHost := os.Getenv("MINIO_HOST")
	minioPort := os.Getenv("MINIO_PORT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")

	minioClient, err := minio.New(fmt.Sprintf("%s:%s", minioHost, minioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Minio: %v", err)
	}

	MinioClient = minioClient
}

func InitConfigcat() {
	sdkKey := os.Getenv("CONFIGCAT_SDK_KEY")
	client := configcat.NewClient(sdkKey)

	ConfigcatClient = client
}

func InitKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokers := []string{os.Getenv("KAFKA_BROKER")}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	KafkaProducerClient = producer
}
