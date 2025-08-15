package rabbitmq

var (
	CLIENT *RabbitMQClient = NewRabbitMQClient()
)

func Init() {
	CLIENT.CreateChannel()

	go CLIENT.ConsumeSeatLock()
}