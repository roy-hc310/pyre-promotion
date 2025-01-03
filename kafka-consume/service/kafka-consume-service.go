package service

import (
	"context"
	"encoding/json"
	"fmt"
	"pyre-promotion/core-internal/utils"
	discount_model "pyre-promotion/feature-discount/v1/model"
	discount_serivce "pyre-promotion/feature-discount/v1/service"

	"github.com/rs/zerolog/log"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumeService struct {
	Client          *kgo.Client
	DiscountService *discount_serivce.DiscountService
}

func NewKafkaService(discountService *discount_serivce.DiscountService, client *kgo.Client) *KafkaConsumeService {
	return &KafkaConsumeService{
		Client:          client,
		DiscountService: discountService,
	}
}

func (k *KafkaConsumeService) Consume(record *kgo.Record) {

	go func() {
		for {
			fetches := k.Client.PollFetches(context.Background())
			errs := fetches.Errors()
			if len(errs) > 0 {
				log.Error().Msg(fmt.Sprint(errs))
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				k.HandleTopic(record)
			}
		}
	}()
}

func (k *KafkaConsumeService) HandleTopic(record *kgo.Record) {
	switch record.Topic {
	case utils.TopicCreateBulkDiscount:
		body := []discount_model.DiscountRequest{}
		err := json.Unmarshal(record.Value, &body)
		if err != nil {
			log.Error().Msg(err.Error())
		}

		ctx, _ := context.WithTimeout(context.Background(), utils.DefaultContextTimeOut)
		_, _, _, err = k.DiscountService.CreateBulkDiscount(ctx, body)
		if err != nil {
			log.Error().Msg(err.Error())
		}

	}
}
