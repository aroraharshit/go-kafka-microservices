package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
)

type Message struct {
	Message      string    `json:"message" `
	Time         time.Time `json:"time" `
	SenderName   string    `json:"senderName" `
	RecieverName string    `json:"recieverName"`
}

func main() {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("Failed to start Kafka producer:", err)
	}
	defer producer.Close()

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		fmt.Println(string(body))
		var message Message
		err = json.Unmarshal(body, &message)
		if err != nil {
			fmt.Println("Error in unmarshalling:", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			http.Error(w, "Failed to marshal message", http.StatusInternalServerError)
			return
		}
		kafkaMessage := &sarama.ProducerMessage{
			Topic: "messages",
			Value: sarama.StringEncoder(messageJSON),
		}

		partition, offset, err := producer.SendMessage(kafkaMessage)
		if err != nil {
			fmt.Println("Error in sending message:", err)
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Message sent to Kafka (Partition: %d, Offset: %d)\n", partition, offset)

		// Send response back to client
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Message sent to Kafka successfully!")
	})

	log.Println("Producer Service running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
