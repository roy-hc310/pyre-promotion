package infrastructure

import (
	"pyre-promotion/core-internal/utils"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaInfra struct {
	Client *kgo.Client
}

func NewKafkaInfra() *KafkaInfra {

	client, err := kgo.NewClient(
		kgo.AllowAutoTopicCreation(),
		kgo.SeedBrokers(strings.Split(utils.GlobalEnv.Kafka.HostPort, ",")...),
		kgo.ConsumerGroup(utils.GlobalEnv.Kafka.ConsumerGroup),
		kgo.ConsumeTopics(
			utils.TopicCreateBulkDiscount,
		),
	)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	return &KafkaInfra{
		Client: client,
	}
}
