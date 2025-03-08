package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Message      string    `json:"message" binding:"required"`
	Time         time.Time `json:"time" binding:"required"`
	SenderName   string    `json:"senderName" binding:"required"`
	RecieverName string    `json:"recieverName" binding:"required"`
}

func main() {
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		var message Message
		err = json.Unmarshal(requestBody, &message)
		if err != nil {
			http.Error(w, "Failed to unmarshal", http.StatusInternalServerError)
			return
		}
		message.Time = time.Now()
		jsonData, err := json.Marshal(message)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		resp, err := http.Post("http://localhost:8081/send", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}

		// Forward the response to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(responseBody)
		defer resp.Body.Close()

		fmt.Fprintf(w, "Message sent via API service!")
	})

	log.Println("API Service running on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
