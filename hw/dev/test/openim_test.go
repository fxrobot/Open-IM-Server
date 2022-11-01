package test

import (
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
)

func TestKafka(t *testing.T) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = false
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer([]string{"116.63.8.212:9092"}, config)
	if err != nil {
		fmt.Println("producer close, err:", err)
		return
	}
	defer client.Close()

	msg := &sarama.ProducerMessage{}
	msg.Topic = "ms2ps_chat_dev"
	msg.Value = sarama.StringEncoder("hello kafka")
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed,", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)

	consumer, err := sarama.NewConsumer([]string{"116.63.8.212:9092"}, config)
	if err != nil {
		fmt.Printf("consumer_test create consumer error %s\n", err.Error())
		return
	}
	partitionList, err := consumer.Partitions("ms2ps_chat_dev") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	defer consumer.Close()
	for partition := range partitionList {
		partitionConsumer, err := consumer.ConsumePartition("ms2ps_chat_dev", int32(partition), sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("try create partition_consumer error %s\n", err.Error())
			return
		}
		defer partitionConsumer.Close()
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
					msg.Offset, msg.Partition, msg.Timestamp.String(), string(msg.Value))
			case err := <-partitionConsumer.Errors():
				fmt.Printf("err :%s\n", err.Error())
			}
		}
	}
}
