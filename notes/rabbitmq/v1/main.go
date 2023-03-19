package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const MQ_URI = "amqp://guest:guest@127.0.0.1:5672/test?connection_attempts=10"

func main() {
	var (
		ctx = context.Background()
	)
	producers := NewProducers(ctx, Conf{
		addr:          MQ_URI,
		exchange:      "parser.excel.headers",
		exchangeType:  amqp.ExchangeHeaders,
		queue:         "parser.excel.headers.xls",
		routingKey:    "",
		consumerTag:   "parser.excel.headers.xls",
		prefetchCount: 1,
		prefetchSize:  0,
		exchangeArgs:  nil,
		QueueArgs:     nil,
		//QueueBindArgs: amqp.Table{
		//	"cc:"
		//},
	})
	producers.Start()
	defer producers.Stop()
	fmt.Println(producers.Publish([]byte("aaa"), "parser.excel.headers", "aaaa", amqp.Table{
		"name":  "aa",
		"title": "bb",
	}))
}

type Conf struct {
	addr          string
	exchange      string
	exchangeType  string
	queue         string
	routingKey    string
	consumerTag   string
	prefetchCount int
	prefetchSize  int
	exchangeArgs  amqp.Table
	QueueArgs     amqp.Table
	QueueBindArgs amqp.Table
}
type Producers struct {
	ctx           context.Context
	conn          *amqp.Connection
	channel       *amqp.Channel
	connNotify    chan *amqp.Error
	channelNotify chan *amqp.Error
	Conf
}

func NewProducers(ctx context.Context, conf Conf) Producers {
	return Producers{
		ctx:  ctx,
		Conf: conf,
	}
}

// Start 启动 mq consumer
func (p *Producers) Start() {
	if err := p.run(); err != nil {
		log.Print("failed to scripts consumer")
		panic(err)
	}
	go p.reConnect(p.ctx)
}

func (p *Producers) init() (err error) {
	if p.conn, err = amqp.Dial(p.addr); err != nil {
		return err
	}
	if p.channel, err = p.conn.Channel(); err != nil {
		return err
	}
	//if err = p.channel.ExchangeDeclare(p.exchange, p.exchangeType, true, false, false, false, p.exchangeArgs); err != nil {
	//	p.Stop()
	//	return
	//}
	//if _, err = p.channel.QueueDeclare(p.queue, true, false, false, false, p.QueueArgs); err != nil {
	//	p.Stop()
	//	return
	//}
	//_ = p.channel.Qos(p.prefetchCount, p.prefetchSize, true)
	//if err = p.channel.QueueBind(p.queue, p.routingKey, p.exchange, false, p.QueueBindArgs); err != nil {
	//	p.Stop()
	//	return
	//}
	return nil
}

func (p *Producers) Stop() {
	if p.conn == nil {
		return
	}
	if !p.conn.IsClosed() {
		// 关闭 SubMsg message delivery
		if p.channel != nil {
			_ = p.channel.Cancel(p.consumerTag, true)
		}
		_ = p.conn.Close()
	}
}

func (p *Producers) run() (err error) {
	if err = p.init(); err != nil {
		return
	}
	p.connNotify = p.conn.NotifyClose(make(chan *amqp.Error))
	p.channelNotify = p.channel.NotifyClose(make(chan *amqp.Error))
	return
}

func (p *Producers) reConnect(ctx context.Context) {
	select {
	case err := <-p.connNotify:
		if err != nil {
			log.Print("rabbitmq consumer - connection NotifyClose: ", err)
		}
	case err := <-p.channelNotify:
		if err != nil {
			log.Print("rabbitmq consumer - channel NotifyClose: ", err)
		}
	case <-ctx.Done():
		return
	}

	// backstop
	if !p.conn.IsClosed() {
		// close message delivery
		if err := p.channel.Cancel(p.consumerTag, true); err != nil {
			log.Print("rabbitmq consumer - channel cancel failed: ", err)
		}
		if err := p.conn.Close(); err != nil {
			log.Print("rabbitmq consumer - conn close failed: ", err)
		}
	}

	// IMPORTANT: 必须清空 Notify，否则死连接不会释放
	for range p.channelNotify {
	}
	for range p.connNotify {
	}

	for i := 0; i < 1000; i++ {
		fmt.Println("rabbitmq consumer - reconnected....")
		if err := p.run(); err != nil {
			log.Print("rabbitmq consumer - failCheck: ", err)
			// sleep 5s reconnect
			time.Sleep(time.Second * 5)
			continue
		}
		fmt.Println("mq run success")
		return
	}
	log.Print("mq consumer reconnect error")
}

func (p Producers) Publish(body []byte, exchange, routingKey string, headers amqp.Table) error {
	publishing := amqp.Publishing{
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		DeliveryMode:    2,
		Priority:        0,
		Body:            body,
		Headers:         headers,
	}
	if err := p.channel.Publish(exchange, routingKey, true, false, publishing); err != nil {
		return errors.WithMessage(err, "rabbitmq publish error")
	}
	return nil
}
