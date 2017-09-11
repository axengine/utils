package nsq

import (
	"log"

	"github.com/nsqio/go-nsq"
)

type dispatcher struct {
	topic  string
	handle func(string, interface{})
}

type consumerDispatcher struct {
	F *dispatcher
	C *nsq.Consumer
}

type TopicDiscoverer struct {
	topics map[string]*consumerDispatcher
	cfg    *nsq.Config
}

//NewTopicDiscoverer 生成多个主题发现的消费者
//topics 订阅主题列表
//channel channel名称
//maxInFlight NSQD在有maxInFlight*25条消息时向下推送
//lookupdHTTPAddrs NSQLOOKUPs的http地址
//handle 消息处理句柄
//PS:对channel不同的topic应该生成不同的消费者
func NewTopicDiscoverer(topics []string, channel string, maxInFlight int, lookupdHTTPAddrs []string, handle func(string, interface{})) *TopicDiscoverer {
	var discoverer TopicDiscoverer
	discoverer.topics = make(map[string]*consumerDispatcher)

	cfg := nsq.NewConfig()
	cfg.UserAgent = USER_AGENT
	cfg.MaxInFlight = maxInFlight

	discoverer.cfg = cfg
	for _, v := range topics {
		cd, err := newConsumerDispatcher(v, channel, cfg, lookupdHTTPAddrs, handle)
		if err != nil {
			log.Fatal(err)
		}
		discoverer.topics[v] = cd
	}
	return &discoverer
}

func newConsumerDispatcher(topic string, channel string, cfg *nsq.Config, lookupdHTTPAddrs []string, handle func(string, interface{})) (*consumerDispatcher, error) {
	var err error
	r := newDispatcher(topic, handle)
	if err != nil {
		return nil, err
	}

	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		return nil, err
	}
	consumer.AddHandler(r)
	err = consumer.ConnectToNSQLookupds(lookupdHTTPAddrs)
	if err != nil {
		log.Fatal(err)
	}
	return &consumerDispatcher{
		C: consumer,
		F: r,
	}, nil
}

func newDispatcher(topic string, handle func(string, interface{})) *dispatcher {
	var disp dispatcher
	disp.topic = topic
	disp.handle = handle
	return &disp
}

//HandleMessage 消息处理句柄 实现接口Handler
func (p *dispatcher) HandleMessage(msg *nsq.Message) error {
	p.handle(p.topic, msg.Body)
	return nil
}
