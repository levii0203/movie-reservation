package rabbitmq

var CLIENT *RabbitMQClient = NewRabbitMQClient()

func Init() {
	//creating channel
	CLIENT.CreateChannel()
	//creating exchange
	CLIENT.CreateExchange()
	//creating a seat lock queue
	CLIENT.NewSeatQueue()
}