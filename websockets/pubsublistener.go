package websockets

import (
	"github.com/go-redis/redis"
)

func ListenForReservations(redisClient *redis.Client, hub *Hub) {
	pubsub := redisClient.Subscribe("worker_reservations")

	_, err := pubsub.Receive()

	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		hub.broadcast <- []byte(msg.Payload)
	}
}