package logrus

import (
	"github.com/linmadan/egglib-go/log"
	"testing"
	"time"
)

const (
	kafkaHost = "192.168.139.129:9092"
	topic     = "pushMessage"
)

func TestKafkaWriter(t *testing.T) {
	w, err := NewKafkaWriter(kafkaHost, "pushMessage", false)
	defer w.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write([]byte(`{"@timestamp":48945123,"msg":"test","host":"127.0.0.1"}`))
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 100)
}

func BenchmarkKafkaWriter_Write(b *testing.B) {
	w, err := NewKafkaWriter(kafkaHost, "pushMessage", false)
	defer w.Close()
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = w.Write([]byte(`{"@timestamp":48945123,"msg":"test","host":"127.0.0.1"}`))
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	time.Sleep(time.Millisecond * 100)
}

func TestHookKafkaWriter(t *testing.T) {
	var Logger log.Logger
	Logger = NewLogrusLogger()
	Logger.SetServiceName("test_hook_kafka")
	Logger.SetLevel("debug")

	w, _ := NewKafkaWriter(kafkaHost, "pushMessage", false)
	Logger.AddHook(w)

	Logger.Debug("HTTP GET 1.0", map[string]interface{}{"host": "192.168.139.129"})
	time.Sleep(time.Millisecond * 100)
}
