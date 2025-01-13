package repositories

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/payloads"
	"os"
)

type NotificationRepository interface {
	SendUserRegistrationEvent(event payloads.UserRegistrationEvent) error
}

func NewNotificationRepository(kafkaProducer sarama.SyncProducer) NotificationRepository {
	return &notificationRepository{
		kafkaProducer: kafkaProducer,
	}
}

type notificationRepository struct {
	kafkaProducer sarama.SyncProducer
}

func (r *notificationRepository) SendUserRegistrationEvent(event payloads.UserRegistrationEvent) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, _, err = r.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: os.Getenv("KAFKA_USER_REGISTRATION_TOPIC"),
		Value: sarama.ByteEncoder(messageBytes),
	})
	if err != nil {
		return err
	}

	return nil
}