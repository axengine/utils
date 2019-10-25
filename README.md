# utils

```
.
├── container
│   ├── bitset #位图
│   │   ├── bitset.go
│   │   └── bitset_test.go
│   ├── blocking_queue #阻塞队列
│   │   ├── blocking_queue.go
│   │   └── blocking_queue_test.go
│   ├── bloom #布隆过滤器
│   │   ├── bloom9.go
│   │   └── bloom9_test.go
│   ├── fifo #fifo
│   │   ├── fifo.go
│   │   └── fifo_test.go
│   ├── omap #有序map 红黑树
│   │   ├── omap.go
│   │   └── omap_test.go
│   └── pqueue #优先级队列 minheap
│       ├── priority_queue.go
│       ├── priority_queue_test.go
│       └── README.md
├── crypto
│   ├── aes256cbc.go
│   ├── aes256cbc_test.go
│   ├── aes.go
│   ├── PKCS.go
│   ├── tripledes_b_test.go
│   ├── tripledes.go
│   └── tripledes_test.go
├── etcd
│   ├── discovery #服务发现
│   │   ├── master.go
│   │   └── worker.go
│   ├── master #选主
│   │   └── master.go
│   ├── readme.md
│   └── sync #锁
│       ├── sync.go
│       └── sync_test.go
├── go.mod
├── hash
│   ├── cityhash
│   │   ├── cityhash.go
│   │   └── cityhash_test.go
│   ├── ketama
│   │   ├── ketama.go
│   │   └── ketama_test.go
│   └── murmurhash3
│       ├── mmhash3.go
│       └── mmhash3_test.go
├── id
│   ├── code #用于生成类似邀请码的东东
│   │   ├── code.go
│   │   └── code_test.go
│   └── uuid #uuid guid生成
│       ├── uuid.go
│       └── uuid_test.go
├── LICENSE
├── log #zap logger的简单使用
│   ├── zap.go
│   └── zap_test.go
├── nsq #nsq client
│   ├── consumer.go
│   └── producer.go
├── README.md
├── redis #redis基础命令+锁
│   ├── mutex.go
│   ├── redis.go
│   └── redis_test.go
└── time
    ├── minheap #基于最小堆的定时器
    │   ├── timer.go
    │   └── timer_test.go
    └── wheel #基于时间轮的定时器
        ├── timer.go
        └── timer_test.go
``` 