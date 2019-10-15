package queue

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/systemz/tasktab/config"
	"log"
)

var (
	ch *amqp.Channel
)

func Listen() {
	conn, err := amqp.Dial("amqp://" + config.RABBITMQ_USERNAME + ":" + config.RABBITMQ_PASSWORD + "@" + config.RABBITMQ_HOST + ":" + config.RABBITMQ_PORT + config.RABBITMQ_VHOST)
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Println("Failed to open a ch")
	}
	//defer ch.Close()

	err = nil
	err = ch.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		log.Println("Failed to set QoS")
	}
	logrus.Println("Connection to RabbitMQ seems OK!")
}
