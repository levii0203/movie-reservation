package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/levii0203/booking-service/internal/config"
	"github.com/levii0203/booking-service/internal/model"
	"github.com/levii0203/booking-service/pkg/helper"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct{
	conn *amqp.Connection
	ch *amqp.Channel
	seatQueue *amqp.Queue
}

func NewRabbitMQClient() *RabbitMQClient {
	return &RabbitMQClient{
		conn: NewAmqpConnection(),
		ch:nil,
		seatQueue:nil,
	}
}


func NewAmqpConnection() *amqp.Connection {
	config.LoadEnv()

	URI := os.Getenv("RABBITMQ_URI")
	conn , err := amqp.Dial(URI)
	if err!=nil {
		helper.FailOnError(err,"Failed to connect to RabbitMQ")
		return nil
	}

	return conn
}

func (r *RabbitMQClient) CreateChannel() {
	if r.conn==nil {
		helper.FailOnError(nil,"amqp connection not established")
	}
	ch,err:=r.conn.Channel()
	if(err!=nil){
		helper.FailOnError(err,"Failed to open a channel")
		return
	}
	
	r.ch = ch
}

func (r *RabbitMQClient) CreateExchange(){
	err:=r.ch.ExchangeDeclare(
		"reservation",
		"direct",
		true,
		true,
		false,
		false,
		nil,
	)
	if err!=nil {
		helper.FailOnError(err, "Failed to declare the exchange")
	}
}

func (r *RabbitMQClient) NewSeatQueue() {
	q,err := r.ch.QueueDeclare(
		"seat",
		false,
		false,
		false,
		false,
		nil,
	)
	if err!=nil {
		helper.FailOnError(err, "Failed to declare a queue")
		return
	}

	err = r.ch.QueueBind(
		q.Name,
		"seat",
		"reservation",
		true,
		nil,
	)
	if err!=nil {
		helper.FailOnError(err, "Failed to declare a queue")
		return
	}
	
	r.seatQueue = &q

}

func (r *RabbitMQClient) AlertSeatLocked(s model.SeatLock) error {
	if r.ch == nil {
        return fmt.Errorf("AMQP channel is nil")
    }

	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()

	data , err := json.Marshal(s)
	if err!=nil {
		return fmt.Errorf("seat lock request failed")
	}
	
	err = r.ch.PublishWithContext(
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

func (r *RabbitMQClient) Close() {
	r.ch.Close()
	r.conn.Close()
}