package registry

import (
	"context"
	"log"

	"sync"

	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/golang/glog"
	"github.com/n3wscott/ledhouse-broker/pkg/registry/api"
)

func (c *ControllerInstance) PubSubControllerRun(ctx context.Context, projectID, subscription string) error {

	var err error
	c.Client, err = pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	c.Subscription = c.Client.Subscription(subscription)
	if ok, err := c.Subscription.Exists(ctx); err != nil || !ok {
		glog.Info("Exists: ", ok, " Error: ", err)
	}

	var mu sync.Mutex
	received := 0

	glog.Info("going to start listening for messages")

	err = c.Subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		mu.Lock()
		defer mu.Unlock()
		received++
		glog.Info("Got message: ", string(msg.Data))
		msg.Ack()

		req := &api.LightRequest{}
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			glog.Errorf("Failed to unmarshal light request: %v", err)
		}
		if err := c.SetLightIntensity(Secret(req.Token), req.Intensity); err != nil {
			glog.Errorf("Failed SetLightIntensity. %v", err)
		}
	})
	if err != nil {
		return err
	}

	glog.Info("done listening for messages")
	return nil
}
