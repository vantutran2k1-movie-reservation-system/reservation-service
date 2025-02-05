package config

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/robfig/cron/v3"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"gorm.io/gorm/logger"
	"log"
	"os"
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
var CronJobManager *cron.Cron

func InitDB() {
	l := logger.Default.LogMode(logger.Silent)
	if AppEnv.GinMode == constants.GinDebugMode {
		l = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppEnv.DbHost,
		AppEnv.DbUser,
		AppEnv.DbPassword,
		AppEnv.DbName,
		AppEnv.DbPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	DB = db
}

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppEnv.RedisHost, AppEnv.RedisPort),
		Password: AppEnv.RedisPassword,
		DB:       AppEnv.RedisDatabase,
	})

	RedisClient = rdb
}

func InitMinio() {
	minioClient, err := minio.New(fmt.Sprintf("%s:%s", AppEnv.MinioHost, AppEnv.MinioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(AppEnv.MinioAccessKey, AppEnv.MinioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Minio: %v", err)
	}

	MinioClient = minioClient
}

func InitConfigcat() {
	client := configcat.NewClient(AppEnv.ConfigcatSdkKey)

	ConfigcatClient = client
}

func InitKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	brokers := []string{AppEnv.KafkaBroker}
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	KafkaProducerClient = producer
}

func InitCronjobManager() {
	c := cron.New(cron.WithSeconds())
	CronJobManager = c
}
