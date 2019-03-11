package websockets

import "taskmaster/redis"

func ListenForReservations(redis redis.Connection, hub *Hub) {
	pubsub := redis.Client.Subscribe("worker_reservations")

	_, err := pubsub.Receive()

	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		hub.broadcast <- []byte(msg.Payload)
	}
}