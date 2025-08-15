package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/levii0203/movie-service/internal/config"
	"github.com/levii0203/movie-service/internal/model"
	"github.com/levii0203/movie-service/internal/repository"
	"github.com/levii0203/movie-service/pkg/helper"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RabbitMQClient struct{
	conn *amqp.Connection
	ch *amqp.Channel
	repo repository.MovieRepository
}

func NewRabbitMQClient() *RabbitMQClient {
	return &RabbitMQClient{
		conn: NewAmqpConnection(),
		ch:nil,
		repo:repository.NewMovieRepository(),
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

func (r *RabbitMQClient) ConsumeSeatLock() {
	if r.repo == nil {
		helper.FailOnError(nil,"repo not declared")
		return
	}
	var wg sync.WaitGroup
	if r.ch == nil {
		helper.FailOnError(nil,"channel not declared")
		return
	}

	msgs, err := r.ch.Consume(
		"seat",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err!=nil {
		helper.FailOnError(err,"Failed to consume")
		return
	}

	wg.Add(1)
	go func(){
		defer wg.Done()
		for d := range msgs {
			fmt.Println(string(d.Body))

			var s model.SeatLock
			err := json.Unmarshal(d.Body, &s)
			if err!=nil {
				continue
			}
			movie_id , _:= primitive.ObjectIDFromHex(s.MovieID)
			fmt.Println(movie_id)
			err = r.repo.UpdateFilledMovie(context.Background(), movie_id, s.Seat)
			if err!=nil {
				log.Println(err.Error())
			}

			d.Ack(false)
		}
	}()

	wg.Wait()
}

func (r *RabbitMQClient) Close() {
	r.ch.Close()
	r.conn.Close()
}