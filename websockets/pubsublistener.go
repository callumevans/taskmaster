package websockets

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func ListenForReservations(redis *redis.Client, hub *Hub) {
	pubsub := redis.Subscribe("worker_reservations")

	_, err := pubsub.Receive()

	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	for msg := range ch {
		var parsedInbound OutboundMessage
		err := json.Unmarshal([]byte(msg.Payload), &parsedInbound)

		if err != nil {
			logrus.Errorf("Error parsing inbound message %s. %s", msg.Payload, err.Error())
		} else {
			hub.Broadcast(parsedInbound)
		}
	}
}