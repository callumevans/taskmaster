package websockets

import "github.com/go-redis/redis"

func ListenForReservations(redis *redis.Client, hub *Hub) {
	pubsub := redis.Subscribe("worker_reservations")

	_, err := pubsub.Receive()

	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		hub.broadcast <- []byte(msg.Payload)
	}
}