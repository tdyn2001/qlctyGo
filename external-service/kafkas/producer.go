package kafkas

import (
	"context"
	"fmt"
	"time"
	"v2/initializers"

	"github.com/segmentio/kafka-go"
)

func Produce(ctx context.Context, topic string) {
	// initialize a counter
	config := initializers.GetConfig()
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{config.KafkaBroker1, config.KafkaBroker2, config.KafkaBroker3},
		Topic:   topic,
	})

	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		dt := time.Now()
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(dt.String()),
			// create an arbitrary message payload for the value
			Value: []byte("Time: " + dt.String()),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", dt.String())
		// sleep for a second
		time.Sleep(time.Second)
	}
}
