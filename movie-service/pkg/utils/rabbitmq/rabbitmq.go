package rabbitmq

import (
	"os"
	"log"
	"movie-service/internal/config"
	"movie-service/pkg/helper"

	amqp "github.com/rabbitmq/amqp091-go"
)


var (
	AmqpClient = amqpConnection()

	Booking = createChannel()
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


func KeepConsumingSeatRequests(ch *amqp.Channel) {
	msgs, err := ch.Consume(
        "seat", 
        "",         
        false,      
        false,      
        false,      
        false,     
        nil,       
    )
	if(err!=nil){
		helper.FailOnError(err, "Failed to consume a queue")
		return
	}
	
	for {
		select {
			case d,ok := <-msgs:
				if !ok {
					log.Panic("Failed to consume seat requests")
					return
				}
				d.Ack(false)
		}
	}

}