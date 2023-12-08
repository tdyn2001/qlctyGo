package kafkas

import (
	"context"
	"fmt"
	"log"
	"os"
	"v2/initializers"

	"github.com/segmentio/kafka-go"
)

func Consume(ctx context.Context, topic, group string) {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	l := log.New(os.Stdout, "kafka reader: ", 0)

	config := initializers.GetConfig()
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.KafkaBroker1, config.KafkaBroker2, config.KafkaBroker3},
		Topic:   topic,
		GroupID: group,
		Logger:  l,
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))
	}
}
