package topic

import (
	"context"
	"cloud.google.com/go/pubsub"
	"github.com/akihokurino/crypto-assistant-server/domain/models"
	"google.golang.org/appengine/log"
	"strings"
)

type PubsubClient interface {
	SendAddress(ctx context.Context, addresses []*models.Address) error
}

type pubsubClient struct {

}

func NewPubsubClient() PubsubClient {
	return &pubsubClient{}
}

func (m *pubsubClient) SendAddress(ctx context.Context, addresses []*models.Address) error {
	client, err := pubsub.NewClient(ctx, "crypto-assistant-dev")
	if err != nil {
		return err
	}

	topic := m.createTopicIfNotExists(ctx, client)

	addressIds := make([]string, len(addresses))
	for i, v := range addresses {
		addressIds[i] = string(v.Id)
	}

	joinedAddressIds := strings.Join(addressIds, ",")

	message := &pubsub.Message{Data: []byte(joinedAddressIds)}

	test := topic.Publish(ctx, message)

	if _, err := test.Get(ctx); err != nil {
		return err
	}

	log.Infof(ctx, "success publish address of %v", joinedAddressIds)

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