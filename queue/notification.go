package queue

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Notification struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func SendNotification(n Notification) {
	msg, err := json.Marshal(n)
	if err != nil {
		log.Printf("Error marshalling notification: %v", err)
	}
	err = ch.Publish(
		"amq.topic",
		"device6",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/json",
			Body:        msg,
		})

	if err != nil {
		log.Printf("Error sending notification: %v", err)
	}
}
