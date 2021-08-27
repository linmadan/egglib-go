package logrus

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
	"sync/atomic"
	"time"
)

var errOutOfMaxSize = errors.New("kafka msg  size is out of limit ")

type KafkaWriter struct {
	producer sarama.SyncProducer
	// 同步写 true:同步 false:异步
	syncWrite   bool
	topic       string
	msgChan     chan []byte
	maxSize     int32
	currentSize int32
	closeChan   chan struct{}
}

func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	if w.syncWrite {
		if _, _, err = w.producer.SendMessage(&sarama.ProducerMessage{Topic: w.topic, Value: sarama.ByteEncoder(p), Timestamp: time.Now()}); err == nil {
			n = len(p)
		}
		return
	}
	if w.currentSize >= w.maxSize {
		fmt.Println(errOutOfMaxSize.Error(), w.currentSize)
		return 0, errOutOfMaxSize
	}
	w.msgChan <- p
	atomic.AddInt32(&w.currentSize, 1)

	return len(p), nil
}

// syncWriteFlag 配置true, 同步写比较慢 55871 ns/op 需要异常错误时使用
func NewKafkaWriter(kafkaHosts string, topic string, syncWriteFlag bool) (*KafkaWriter, error) {
	writer := &KafkaWriter{
		syncWrite: syncWriteFlag,
		topic:     topic,
		maxSize:   10000,
		msgChan:   make(chan []byte, 10000),
		closeChan: make(chan struct{}),
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 2
	config.Producer.RequiredAcks = sarama.NoResponse
	brokerList := strings.Split(kafkaHosts, ",")
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	writer.producer = producer
	go writer.ConsumeMsg()
	return writer, nil
}

func (w *KafkaWriter) ConsumeMsg() {
	for {
		select {
		case <-w.closeChan:
			return
		case m, ok := <-w.msgChan:
			if ok {
				atomic.AddInt32(&w.currentSize, -1)
				if _, _, err := w.producer.SendMessage(&sarama.ProducerMessage{
					Topic: w.topic,
					Value: sarama.ByteEncoder(m),
				}); err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func (w *KafkaWriter) Close() {
	close(w.msgChan)
	w.closeChan <- struct{}{}
	w.producer.Close()
}
