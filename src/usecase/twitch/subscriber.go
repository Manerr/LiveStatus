package twitch

import (
	"LiveStatus/src/domain"
	"context"
	"errors"
	esb "github.com/dnsge/twitch-eventsub-bindings"
	esf "github.com/dnsge/twitch-eventsub-framework"
	"log"
)

type Subscriber interface {
	GetSubscriptions() (*esb.RequestStatus, error)
	UnsubscribeAll() error
	SubscribeAll(broadcasterUserIds []string) (int, int, error)
}

func NewSubscriber(clientId string, appToken string, webhookUrl string, webhookSecret string) Subscriber {
	return &subscriber{
		client:        esf.NewSubClient(esf.NewStaticCredentials(clientId, appToken)),
		webhookUrl:    webhookUrl,
		webhookSecret: webhookSecret,
	}
}

type subscriber struct {
	client        *esf.SubClient
	webhookUrl    string
	webhookSecret string
}

func (s *subscriber) GetSubscriptions() (*esb.RequestStatus, error) {
	return s.client.GetSubscriptions(context.Background(), esf.StatusAny)
}

func (s *subscriber) UnsubscribeAll() error {
	subscriptions, err := s.GetSubscriptions()

	for _, sub := range subscriptions.Data {
		err = s.client.Unsubscribe(context.Background(), sub.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// SubscribeAll subscribes to the "stream.online" and "stream.offline" event types for the provided broadcasterUserIds.
// Returns:
//   - int: the total cost used
//   - int: the maximum total cost allowed
//   - error: any error that occurred
func (s *subscriber) SubscribeAll(broadcasterUserIds []string) (int, int, error) {
	if len(broadcasterUserIds) == 0 {
		return 0, 0, errors.New("no broadcasterUserIds provided")
	}

	var latestResponse *esb.RequestStatus
	var err error

	for _, broadcasterUserId := range broadcasterUserIds {
		for _, subType := range domain.SubscriptionList {
			latestResponse, err = s.client.Subscribe(context.Background(), &esf.SubRequest{
				Type: subType,
				Condition: esb.ConditionChannelUpdate{
					BroadcasterUserID: broadcasterUserId,
				},
				Callback: s.webhookUrl,
				Secret:   s.webhookSecret,
			})

			if err != nil {
				log.Fatalf("Error subscribing: %v\n", err)
			}
		}
	}

	return latestResponse.TotalCost, latestResponse.MaxTotalCost, nil
}
