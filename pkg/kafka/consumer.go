package kafka

import (
	"context"
	"errors"
	"fmt"
	"go-wal/config"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

type Processor interface {
	Processor(ctx context.Context, topic string, message []byte)
}

func NewConsumer(ctx context.Context, p Processor, cnf *config.Kafka, topic, groupId string) {
	fmt.Printf("NewConsumer: %v\n", cnf)
	cfg := sarama.NewConfig()

	cg, err := sarama.NewConsumerGroup(cnf.Brokers, groupId, cfg)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	consumer := processor{
		ready: make(chan bool),
		ctx:   ctx,
		p:     p,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			err := cg.Consume(ctx, []string{topic}, &consumer)
			if err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Fatalf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready
	// logger.Info(ctx, "Consumer started", zap.String("group_id", config.Config.Kafka.ConsumerGroupID), zap.Strings("topics", []string{constant.OrderServiceOrderStatusTopic}))
	log.Println("Consumer started", "group_id", groupId, "topics", []string{topic})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signals:
		// logger.Info(ctx, "Received shutdown signal. Gracefully shutting down consumer ...")
		log.Println("Received shutdown signal. Gracefully shutting down consumer ...")
	case err := <-cg.Errors():
		// logger.Fatal(ctx, err, zap.String("message", "Error from consumer"))
		log.Fatalf("Error from consumer: %v", err)
	}

	log.Println(ctx, "Start closing consumer")
	cancel()
	wg.Wait()
	log.Println(ctx, "Wait group done, and add delay 10s to consumer can finish in-flight request")
	// Add a delay after consumer is closed to allow in-flight requests to finish
	time.Sleep(10 * time.Second)
	err = cg.Close()
	if err != nil {
		log.Fatalf("Error closing consumer group: %v", err)
	} else {
		log.Println(ctx, "Consumer closed")
	}

}

type processor struct {
	ready chan bool
	ctx   context.Context
	p     Processor
}

// Setup implements sarama.ConsumerGroupHandler.
func (consumer *processor) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup implements sarama.ConsumerGroupHandler.
func (*processor) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implements sarama.ConsumerGroupHandler.
func (consumer *processor) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				fmt.Printf("claim.Messages() is closed\n")
				return nil
			}

			// consumer.ctx = context.WithValue(consumer.ctx, constant.XRequestID, logger.GenerateTraceId())
			ctxNoCancel := context.WithoutCancel(consumer.ctx)
			consumer.p.Processor(ctxNoCancel, message.Topic, message.Value)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
