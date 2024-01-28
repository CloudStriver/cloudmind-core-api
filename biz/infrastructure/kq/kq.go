package kq

import (
	"context"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/config"
	"github.com/CloudStriver/cloudmind-core-api/biz/infrastructure/consts"
	"github.com/CloudStriver/cloudmind-mq/app/util/message"
	"github.com/CloudStriver/go-pkg/utils/pconvertor"
	"github.com/CloudStriver/go-pkg/utils/util/batcher"
	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"hash/crc32"
)

type CreateNotificationsKq struct {
	*kq.Pusher
	*batcher.Batcher
}

type ReadNotificationsKq struct {
	*kq.Pusher
	*batcher.Batcher
}

type CleanNotificationKq struct {
	*kq.Pusher
	*batcher.Batcher
}

func NewCreateNotificationsKq(c *config.Config) *CreateNotificationsKq {
	crc := crc32.MakeTable(0xD5828281)
	pusher := kq.NewPusher(c.CreateNotificationsKq.Brokers, c.CreateNotificationsKq.Topic)
	b := batcher.New(
		batcher.WithSize(consts.BatcherSize),
		batcher.WithBuffer(consts.BatcherBuffer),
		batcher.WithWorker(consts.BatcherWorker),
		batcher.WithInterval(consts.BatcherInterval),
	)
	b.Sharding = func(key string) int {
		pid := crc32.Checksum(pconvertor.String2Bytes(key), crc)
		return int(pid) % consts.BatcherWorker
	}
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		var msgs []*message.CreateNotificationsMessage
		for _, vs := range val {
			for _, v := range vs {
				msgs = append(msgs, v.(*message.CreateNotificationsMessage))
			}
		}
		kd, err := sonic.Marshal(msgs)
		if err != nil {
			logx.Errorf("Batcher.Do json.Marshal msgs: %v error: %v", msgs, err)
		}
		if err = pusher.Push(pconvertor.Bytes2String(kd)); err != nil {
			logx.Errorf("KafkaPusher.Push kd: %s error: %v", pconvertor.Bytes2String(kd), err)
		}
	}
	b.Start()
	return &CreateNotificationsKq{
		Pusher:  pusher,
		Batcher: b,
	}
}

func NewReadNotificationsKq(c *config.Config) *ReadNotificationsKq {
	crc := crc32.MakeTable(0xD5828281)
	pusher := kq.NewPusher(c.CreateNotificationsKq.Brokers, c.CreateNotificationsKq.Topic)
	b := batcher.New(
		batcher.WithSize(consts.BatcherSize),
		batcher.WithBuffer(consts.BatcherBuffer),
		batcher.WithWorker(consts.BatcherWorker),
		batcher.WithInterval(consts.BatcherInterval),
	)
	b.Sharding = func(key string) int {
		pid := crc32.Checksum(pconvertor.String2Bytes(key), crc)
		return int(pid) % consts.BatcherWorker
	}
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		var msgs []*message.ReadNotificationsMessage
		for _, vs := range val {
			for _, v := range vs {
				msgs = append(msgs, v.(*message.ReadNotificationsMessage))
			}
		}
		kd, err := sonic.Marshal(msgs)
		if err != nil {
			logx.Errorf("Batcher.Do json.Marshal msgs: %v error: %v", msgs, err)
		}
		if err = pusher.Push(pconvertor.Bytes2String(kd)); err != nil {
			logx.Errorf("KafkaPusher.Push kd: %s error: %v", pconvertor.Bytes2String(kd), err)
		}
	}
	b.Start()
	return &ReadNotificationsKq{
		Pusher:  pusher,
		Batcher: b,
	}
}
func NewCleanNotificationKq(c *config.Config) *CleanNotificationKq {
	crc := crc32.MakeTable(0xD5828281)
	pusher := kq.NewPusher(c.CreateNotificationsKq.Brokers, c.CreateNotificationsKq.Topic)
	b := batcher.New(
		batcher.WithSize(consts.BatcherSize),
		batcher.WithBuffer(consts.BatcherBuffer),
		batcher.WithWorker(consts.BatcherWorker),
		batcher.WithInterval(consts.BatcherInterval),
	)
	b.Sharding = func(key string) int {
		pid := crc32.Checksum(pconvertor.String2Bytes(key), crc)
		return int(pid) % consts.BatcherWorker
	}
	b.Do = func(ctx context.Context, val map[string][]interface{}) {
		var msgs []*message.CleanNotificationMessage
		for _, vs := range val {
			for _, v := range vs {
				msgs = append(msgs, v.(*message.CleanNotificationMessage))
			}
		}
		data, err := sonic.Marshal(msgs)
		if err != nil {
			logx.Errorf("Batcher.Do json.Marshal msgs: %v error: %v", msgs, err)
		}
		if err = pusher.Push(pconvertor.Bytes2String(data)); err != nil {
			logx.Errorf("KafkaPusher.Push kd: %s error: %v", pconvertor.Bytes2String(data), err)
		}
	}
	b.Start()
	return &CleanNotificationKq{
		Pusher:  pusher,
		Batcher: b,
	}
}
