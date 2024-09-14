package wal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/jackc/pgx"
	"go-wal/config"
	"go-wal/pkg/kafka"
	"go-wal/pkg/logger"
	"strconv"
)

const (
	defaultPostgresTimeline = -1
	columnNameForTopic      = "topic"
	columnNameForPayload    = "payload"
)

type PGWalListener struct {
	slot    string
	lsn     string
	lastWal uint64
	buffer  []byte
	metrics *Metrics
	chunks  string

	kafkaProducer *kafka.Producer
}

func NewPgWalListener(kafkaProducer *kafka.Producer) *PGWalListener {
	// Init dependencies from app level => Avoid import cycle
	return &PGWalListener{
		kafkaProducer: kafkaProducer,
		slot:          "my_slot",
		lsn:           "0/0",
		metrics:       NewMetrics(),
		chunks:        "1",
	}
}

func (p *PGWalListener) Start(ctx context.Context) error {
	// Start the listener
	pgxCfg := pgx.ConnConfig{
		Host:     config.Config.Database.Host,
		Port:     uint16(config.Config.Database.Port),
		User:     config.Config.Database.User,
		Password: config.Config.Database.Password,
		Database: config.Config.Database.Name,
	}
	conn, err := pgx.ReplicationConnect(pgxCfg)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to connect to replication")
		return err
	}
	defer func(conn *pgx.ReplicationConn) {
		err := conn.Close()
		if err != nil {
			logger.Error(ctx).Err(err).Msg("Failed to close replication connection")
		}
	}(conn)
	lsnAsInt, err := pgx.ParseLSN(p.lsn)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to parse LSN")
		return err
	}
	var outputParams []string
	if p.chunks != "" {
		outputParams = append(outputParams, fmt.Sprintf("\"write-in-chunks\" '%s'", p.chunks))
	}
	if err := conn.StartReplication(p.slot, lsnAsInt, defaultPostgresTimeline, outputParams...); err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to start replication")
		return err
	}
	logger.Info(ctx).Msg("Started wal listener")
	p.listen(ctx, conn)
	return nil
}

// listen is main infinite process for handle this replication slot
func (p *PGWalListener) listen(ctx context.Context, rc *pgx.ReplicationConn) {
	for {
		r, err := rc.WaitForReplicationMessage(context.Background())
		if err != nil {
			logger.Error(ctx).Err(err).Msg("Failed to wait for replication message")
		}

		if r != nil {
			if r.ServerHeartbeat != nil {
				p.handleHeartbeat(ctx, rc, r.ServerHeartbeat.ServerWalEnd)
			} else if r.WalMessage != nil {
				p.handleMessage(ctx, rc, r.WalMessage.WalData, r.WalMessage.WalStart, r.WalMessage.ServerWalEnd)
			}
		}
	}
}

// handleMessage parses WAL data, send message to Kafka, and sends standby status
func (p *PGWalListener) handleMessage(ctx context.Context, rc *pgx.ReplicationConn, walData []byte, walStart uint64, walEnd uint64) {
	p.lastWal = walStart
	p.metrics.lastCommittedWal = p.lastWal

	//buffering all input bytes
	p.buffer = append(p.buffer, walData...)
	p.metrics.bufferSize = len(p.buffer)

	var change Wal2JsonMessage
	//trying to deserialize to JSON
	if err := json.Unmarshal(p.buffer, &change); err == nil {
		logger.Info(ctx).Any("change", change).Msg("Received WAL message")
		messages := p.getAsMessages(ctx, change)
		p.produceMessage(ctx, messages)

		p.buffer = []byte{}
		p.metrics.bufferSize = 0
	}

	if err := p.sendStandBy(rc); err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to send standby status")
		p.metrics.totalErrors++
	}
}

// produceMessage iterates over parsed messages and publish it to Kafka
func (p *PGWalListener) produceMessage(ctx context.Context, messages *[]Message) {
	if messages == nil || len(*messages) == 0 {
		return
	}

	var saramaMessages []*sarama.ProducerMessage
	for _, message := range *messages {
		logger.Info(ctx).Msg(fmt.Sprintf("sending message %+v", message))
		msgData, err := json.Marshal(message.Value)
		if err != nil {
			logger.Error(ctx).Err(err).Msg("Failed to marshal message value")
			continue
		}
		msg := p.kafkaProducer.BuildMessage(ctx, message.Topic, "", msgData)
		saramaMessages = append(saramaMessages, &msg)
	}
	if err := p.kafkaProducer.Publish2KafkaMessages(ctx, saramaMessages); err != nil {
		p.metrics.totalErrors++
	}
	p.metrics.totalMessages++
}

// getAsMessages parses topic and value from Wal2Json data
func (p *PGWalListener) getAsMessages(ctx context.Context, change Wal2JsonMessage) *[]Message {
	var messages []Message

	if len(change.Change) == 0 {
		return &messages
	}

	for _, item := range change.Change {
		//watch only inserts
		logger.Info(ctx).Any("item", item).Msg("change.Change.item")
		if item.Kind == "insert" {
			var topic, value string

			for i, name := range item.ColumnNames {
				if name == columnNameForTopic {
					topic = p.getParsedValue(ctx, item.ColumnValues[i])
				} else if name == columnNameForPayload {
					value = p.getParsedValue(ctx, item.ColumnValues[i])
				}
			}
			messages = append(messages, Message{Topic: topic, Value: value})
		}
	}

	return &messages
}

func (p *PGWalListener) getParsedValue(ctx context.Context, input interface{}) string {
	switch v := input.(type) {
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case nil:
		return "null"
	default:
		logger.Error(ctx).Err(fmt.Errorf("unknown type %T", v)).Msg("getParsedValue")
	}

	return ""
}

// handleHeartbeat sends standby status
func (p *PGWalListener) handleHeartbeat(ctx context.Context, rc *pgx.ReplicationConn, walEnd uint64) {
	p.lastWal = walEnd
	p.metrics.totalHeartbeats++
	p.metrics.lastCommittedWal = p.lastWal

	if err := p.sendStandBy(rc); err != nil {
		logger.Error(ctx).Err(err).Msg("Failed to send standby status")
		p.metrics.totalErrors++
	}
}

func (p *PGWalListener) sendStandBy(rc *pgx.ReplicationConn) error {
	status, err := pgx.NewStandbyStatus(p.lastWal)
	if err != nil {
		return err
	}

	err = rc.SendStandbyStatus(status)
	return err
}
