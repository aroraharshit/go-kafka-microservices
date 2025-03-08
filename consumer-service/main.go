package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	Message      string    `json:"message" binding:"required"`
	Time         time.Time `json:"time" binding:"required"`
	SenderName   string    `json:"senderName" binding:"required"`
	RecieverName string    `json:"recieverName" binding:"required"`
}

func main() {
	brokers := []string{"localhost:9092"}
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatal("Failed to start Kafka consumer", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("messages", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to create partion consumer", err)
	}
	defer partitionConsumer.Close()

	log.Println("Consumer listening for messages...")

	for msg := range partitionConsumer.Messages() {
		var Message Message
		json.Unmarshal(msg.Value, &Message)
		fmt.Printf("New message received: Message=%s, SenderName=%s\n", Message.Message, Message.SenderName)
	}
}
