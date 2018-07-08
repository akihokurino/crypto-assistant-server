package topic

import (
	"context"
	"cloud.google.com/go/pubsub"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"google.golang.org/appengine/log"
)

type PubsubClient interface {
	SendAddress(ctx context.Context, address *models.Address) error
}

type pubsubClient struct {

}

func NewPubsubClient() PubsubClient {
	return &pubsubClient{}
}

func (m *pubsubClient) SendAddress(ctx context.Context, address *models.Address) error {
	client, err := pubsub.NewClient(ctx, "crypto-assistant-dev")
	if err != nil {
		return err
	}

	topic := m.createTopicIfNotExists(ctx, client)

	message := &pubsub.Message{Data: []byte(address.Id)}

	test := topic.Publish(ctx, message)

	if _, err := test.Get(ctx); err != nil {
		return err
	}

	log.Infof(ctx, "success publish address of %v", address.Value)

	return nil
}

func (m *pubsubClient) createTopicIfNotExists(ctx context.Context, c *pubsub.Client) *pubsub.Topic {
	topicName := "calc-assets"
	t := c.Topic(topicName)

	ok, err := t.Exists(ctx)
	if err != nil {
		panic(err)
	}

	if ok {
		return t
	}

	t, err = c.CreateTopic(ctx, topicName)
	if err != nil {
		panic(err)
	}

	return t
}