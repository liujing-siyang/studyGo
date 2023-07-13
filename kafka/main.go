package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

// 基于sarama第三方库开发的kafka client

func main() {
	// producer()
	consumer()
	// MyConsumerGroup()
}

func producer() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	config.Net.MaxOpenRequests = 1                            // 在broker没有响应请求前能发送的最大请求数，默认为5，开启幂等性需要
	config.Producer.Idempotent = true                         // 开启幂等性，事务依赖幂等性
	config.Producer.Transaction.ID = "transaction.id"         // 开启事务需指定事务id

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "first"
	msg.Value = sarama.StringEncoder("Transaction test")
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"101.37.25.231:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer func() {
		if err != nil {
			client.AbortTxn()
		}
		client.Close()
	}()

	//开启事务
	client.BeginTxn()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)

	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	//提交事务，没有提交事务数据还是发送过去了？？
	err = client.CommitTxn()
	// if err == nil {
	// 	err = fmt.Errorf("事务提交失败测试")
	// }

	fmt.Printf("pid:%v offset:%v\n", pid, offset)
}

func consumer() {
	consumer, err := sarama.NewConsumer([]string{"101.37.25.231:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("second") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	var wg sync.WaitGroup
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("second", int32(partition), sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
			wg.Done()
		}(pc)
	}
	wg.Wait()
}

type ConsumerGroup struct {
	brokers           []string
	topics            []string
	startOffset       int64
	version           string
	group             string
	channelBufferSize int
	assignor          string
}

var brokers = []string{"121.199.71.218:9092"}
var topics = []string{"operate_cumulative_commission"}
var group = "kafkagroup"
var assignor = "roundrobin"

// var brokers = []string{"101.37.25.231:9092"}
// var topics = []string{"second"}
// var group = "test"
// var assignor = "roundrobin"

type Consumer struct {
	id    int
	ready chan bool
}

func NewConsumerGroup() *ConsumerGroup {
	return &ConsumerGroup{
		brokers:           brokers,
		topics:            topics,
		group:             group,
		channelBufferSize: 1000,
		version:           "2.8.1",
		assignor:          assignor,
	}
}

func (cg *ConsumerGroup) Connect(consumer  Consumer) func() {
	log.Infoln("kafka init...")

	version, err := sarama.ParseKafkaVersion(cg.version)
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version
	// 分区分配策略
	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.ChannelBufferSize = cg.channelBufferSize // channel长度
	config.Consumer.Offsets.AutoCommit.Enable = false
	//会话管理
	// config.Consumer.Group.Session.Timeout = 10 * time.Second
	// config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	//设置消费者拉取一批次的时间,一次最小拉取8字节，如果没有8字节，最多等待500ms
	// config.Consumer.Fetch.Min = 8
	// config.Consumer.MaxWaitTime = 500 * time.Microsecond

	// 创建client
	newClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	// 获取所有的topic
	topics, err := newClient.Topics()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("topics: ", topics)

	// 根据client创建consumerGroup
	client, err := sarama.NewConsumerGroupFromClient(cg.group, newClient)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	partitions, err := newClient.Partitions("operate_cumulative_commission")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("partitions: ", partitions)

	// go TimerSync(ctx)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, cg.topics, consumer); err != nil {
				// 当setup失败的时候，error会返回到这里
				log.Errorf("Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	// <-consumer.ready
	log.Info("running")
	// 保证在系统退出时，通道里面的消息被消费
	return func() {
		log.Info("kafka close")
		cancel()
		wg.Wait()
		if err = client.Close(); err != nil {
			log.Errorf("Error closing client: %v", err)
		}
	}
}

func (cg *ConsumerGroup) Connect2() func() {
	log.Infoln("kafka init...")

	version, err := sarama.ParseKafkaVersion(cg.version)
	if err != nil {
		log.Fatalf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version
	// 分区分配策略
	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", assignor)
	}

	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.ChannelBufferSize = cg.channelBufferSize // channel长度
	// config.Consumer.Offsets.AutoCommit.Enable = false
	//会话管理
	// config.Consumer.Group.Session.Timeout = 10 * time.Second
	// config.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	//设置消费者拉取一批次的时间,一次最小拉取8字节，如果没有8字节，最多等待500ms
	// config.Consumer.Fetch.Min = 8
	// config.Consumer.MaxWaitTime = 500 * time.Microsecond

	// 创建client
	newClient, err := sarama.NewClient(brokers, config)
	if err != nil {
		log.Fatal(err)
	}
	// 获取所有的topic
	topics, err := newClient.Topics()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("topics: ", topics)

	partitions, err := newClient.Partitions("second")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("partitions: ", partitions)

	consumer := Consumer{
		id:    1,
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	// 根据client创建consumerGroup
	client, err := sarama.NewConsumerGroupFromClient(cg.group, newClient)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//当只有一个消费者时，读取历史数据为何没有读取所有分区的数据
		for {
			if err := client.Consume(ctx, cg.topics, &consumer); err != nil {
				// 当setup失败的时候，error会返回到这里
				log.Errorf("Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	log.Info("111")
	time.Sleep(2 * time.Second)
	consumer2 := Consumer{
		id:    2,
		ready: make(chan bool),
	}
	log.Info("create 2")
	// client2, err := sarama.NewConsumerGroupFromClient(cg.group, newClient)
	// if err != nil {
	// 	log.Fatalf("Error creating consumer group client: %v", err)
	// }
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, cg.topics, &consumer2); err != nil {
				// 当setup失败的时候，error会返回到这里
				log.Errorf("Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			consumer2.ready = make(chan bool)
		}
	}()
	<-consumer2.ready
	log.Info("222")
	// 保证在系统退出时，通道里面的消息被消费
	return func() {
		log.Info("kafka close")
		cancel()
		wg.Wait()
		if err = client.Close(); err != nil {
			log.Errorf("Error closing client: %v", err)
		}
	}
}

func (cg *ConsumerGroup) NewConsumer(newClient sarama.Client, wg *sync.WaitGroup, ctx context.Context) sarama.ConsumerGroup {
	consumer := Consumer{
		ready: make(chan bool),
	}
	// 根据client创建consumerGroup
	client, err := sarama.NewConsumerGroupFromClient(cg.group, newClient)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		//当只有一个消费者时，读取历史数据为何没有读取所有分区的数据
		for {
			if err := client.Consume(ctx, cg.topics, &consumer); err != nil {
				// 当setup失败的时候，error会返回到这里
				log.Errorf("Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println(ctx.Err())
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	log.Infoln("xfz")
	return client
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (k Consumer) Setup(session sarama.ConsumerGroupSession) error {

	// session.ResetOffset("operate_cumulative_commission", 0, 6995, "")
	// session.ResetOffset("second", 0, 0, "")
	// session.ResetOffset("second", 1, 0, "")
	// session.ResetOffset("second", 2, 0, "")
	log.Infof("setup [id:%d] [memberid:%v] [claim:%s]", k.id, session.Claims(), session.MemberID())
	// Mark the consumer as ready
	// close(k.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (k Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Info("cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (k Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// <https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29>
	// 具体消费消息
	for message := range claim.Messages() {

		// log.Infof("[id:%d] [Consume:%s] [topic:%s] [partiton:%d] [offset:%d] [value:%s] [time:%v]", k.id, session.MemberID(),
		// 	message.Topic, message.Partition, message.Offset, string(message.Value), message.Timestamp)
		log.Infof("[id:%d] [Consume:%s] [topic:%s] [partiton:%d] [offset:%d]", k.id, session.MemberID(),
		message.Topic, message.Partition, message.Offset)
		// go Do(message)
		// 更新位移
		session.MarkMessage(message, "")
		//提交 offset
		// session.Commit()
	}
	return nil
}

// 并发安全的map
var m = sync.Map{}
var ch = make(chan *AffectRow, 10)

func Do(message *sarama.ConsumerMessage) (row AffectRow, err error) {
	err = json.Unmarshal(message.Value, &row)
	if err != nil {
		log.Errorf("[%s]反序列化失败：%v\n", string(message.Value), err)
		return
	}
	m.Store(fmt.Sprintf("%s,%s,%s,%s", row.TableName, row.SQLType, row.Columns.StatisticsMonth, row.Columns.FollowCity), row)

	// if !ok {
	// 	err = fmt.Errorf("[%s]记录数据失败", string(message.Value))
	// 	log.Error(err)
	// }
	log.Infof("offset:%d,%v",message.Offset,row)
	ch <- &row
	return
}

func TimerSync(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			// 这里处理定时器触发时的逻辑
			m.Range(updateForm)
			m.Range(func(k, v interface{}) bool {
				m.Delete(k)
				return true
			})
		case <-ch:
			ticker.Reset(1 * time.Minute)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func updateForm (key, value interface{}) bool{
	log.Infof("key:[%v],value:[%v]",key,value)
	return true
}

type AffectRow struct {
	TableName string `json:"tableName"`
	SQLType   string `json:"sqlType"`
	Columns   struct {
		ContCode        string `json:"cont_code"`
		FollowCity      string `json:"follow_city"`
		StatisticsMonth string `json:"statistics_month"`
	} `json:"columns"`
}

func MyConsumerGroup() {
	cg := NewConsumerGroup()
	consumer := Consumer{
		id:    1,
		ready: make(chan bool),
	}
	c := cg.Connect(consumer)
	// consumer2 := Consumer{
	// 	id:    2,
	// 	ready: make(chan bool),
	// }
	// c2 := cg.Connect(consumer2)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Warnln("terminating: via signal")
	}
	c()
	// c2()
}
