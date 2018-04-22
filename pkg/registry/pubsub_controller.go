package registry

import (
	"context"
	"log"

	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/golang/glog"
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
	})
	if err != nil {
		return err
	}

	glog.Info("done listening for messages")
	return nil
}
