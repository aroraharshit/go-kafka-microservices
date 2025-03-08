This project demonstrates a simple microservices architecture in Go using Kafka for real-time messaging. It consists of three services:
1)API Service (api/main.go): Exposes an endpoint to send messages to the Producer service.
2)Producer Service (producer/main.go): Receives messages from the API service and publishes them to a Kafka topic.
3)Consumer Service (consumer/main.go): Listens to the Kafka topic and processes incoming messages.

Architecture
(API Service) → (Producer Service) → Kafka → (Consumer Service)

-)The API Service provides an HTTP endpoint (/publish) where users can send messages.
-)The Producer Service pushes the received message to a Kafka topic (messages).
-)The Consumer Service listens for new messages in the Kafka topic and logs them.

Technologies Used--
Go (Golang)
Kafka
Sarama (Go client for Apache Kafka)

How It Works

-)The API service provides an endpoint where users can manually send a message.
-)The Producer receives the message from the API and sends it to the Kafka topic "messages".
-)The Consumer listens for new messages in the Kafka topic and logs them.
