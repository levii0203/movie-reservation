package rabbitmq

import (
	"booking-service/internal/config"
	"booking-service/internal/model"
	"booking-service/pkg/helper"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


var (
	AmqpClient = amqpConnection()

	Booking = createChannel()

	SeatQueue = NewSeatQueue(Booking)
)


func amqpConnection() *amqp.Connection {
	config.LoadEnv()

	URI := os.Getenv("RABBITMQ_URI")
	conn , err := amqp.Dial(URI)
	if err!=nil {
		helper.FailOnError(err,"Failed to connect to RabbitMQ")
		return nil
	}

	return conn
}

func createChannel() *amqp.Channel {
	ch,err:=AmqpClient.Channel()
	if(err!=nil){
		helper.FailOnError(err,"Failed to open a channel")
		return nil
	}
	return ch
}

func NewSeatQueue(ch *amqp.Channel) *amqp.Queue {
	q,err := ch.QueueDeclare(
		"seat",
		false,
		false,
		false,
		false,
		nil,
	)
	if err!=nil {
		helper.FailOnError(err, "Failed to declare a queue")
		return nil
	}

	return &q
}

func AlertSeatLocked(ch *amqp.Channel, s model.SeatLock ) error {
	if ch == nil {
        return fmt.Errorf("AMQP channel is nil")
    }

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	data , err := json.Marshal(s)
	if err!=nil {
		return fmt.Errorf("seat lock request failed")
	}
	err = ch.PublishWithContext(
		ctx,
		"",
		"seat",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: []byte(data),
		},
	)	
	if err!=nil{
		helper.FailOnError(err, "Failed to publish a message")
		return err
	}
	return nil
}