package main

import (
	"context"
	"log"
	"os"

	"flag"

	"os/signal"
	"syscall"

	"cloud.google.com/go/pubsub"
	"github.com/golang/glog"
)

var options struct {
	ProjectID string
	Topic     string
}

func init() {
	flag.StringVar(&options.ProjectID, "projectId", "", "specify the gcp projectId")
	flag.StringVar(&options.Topic, "topic", "", "specify the pub/sub topic")
	flag.Parse()
}

func main() {
	if err := run(); err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		glog.Fatalln(err)
	}
}

func run() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go cancelOnInterrupt(ctx, cancelFunc)

	return runWithContext(ctx)
}

var (
	topic  *pubsub.Topic
	client *pubsub.Client
)

func runWithContext(ctx context.Context) error {

	var err error
	client, err = pubsub.NewClient(ctx, options.ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	topic = client.Topic(options.Topic)
	if _, err := topic.Exists(ctx); err != nil {
		log.Fatal(err)
	}

	send(topic)

	return nil
}

func cancelOnInterrupt(ctx context.Context, f context.CancelFunc) {
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-term:
			glog.Infof("Received SIGTERM, exiting gracefully...")
			f()
			os.Exit(0)
		case <-ctx.Done():
			os.Exit(0)
		}
	}
}

func send(topic *pubsub.Topic) {
	ctx := context.Background()

	msg := &pubsub.Message{
		Data: []byte("hello"),
	}

	if _, err := topic.Publish(ctx, msg).Get(ctx); err != nil {
		glog.Error("Could not publish message:", err)
		return
	}

	glog.Info("Message published.")
}
