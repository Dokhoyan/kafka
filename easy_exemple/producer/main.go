package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v7"
)

const (
	brokerAddress = "localhost:9092"
	topicName     = "test-topic"
)

func main() {
	producer, err := newSyncProducer([]string{brokerAddress})
	if err != nil {
		log.Fatalf("failed to start producer: %v\n", err.Error())
	}

	defer func() {
		if err = producer.Close(); err != nil {
			log.Fatalf("failed to close producer: %v\n", err.Error())
		}
	}()

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(gofakeit.StreetName()),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to send message in Kafka: %v\n", err.Error())
		return
	}

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)
}


func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig() 
	config.Producer.RequiredAcks = sarama.WaitForAll //ждём подтверждения от всех in-sync реплик (ISR).
	config.Producer.Retry.Max = 5					 // максимальное число повторных попыток (retries) при ошибке отправки
	config.Producer.Return.Successes = true			 // родьюсер должен возвращать события успешной отправки в канал Successes

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}