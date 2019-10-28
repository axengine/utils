# utils

```
├── codec #一种混淆数据的编解码方法
│   ├── codec.go
│   ├── codec_test.go
│   └── README.md
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
│   ├── fifo #固定大小的fifo，特点是使用固定大小内存
│   │   ├── fifo.go
│   │   └── fifo_test.go
│   ├── omap #有序的map - 红黑树
│   │   ├── omap.go
│   │   └── omap_test.go
│   └── pqueue #优先级队列 - 小堆
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
│   └── sync #分布式锁
│       ├── sync.go
│       └── sync_test.go
├── go.mod
├── go.sum
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
│   ├── code #生成代码，通常用于生成邀请码
│   │   ├── code.go
│   │   └── code_test.go
│   └── uuid #uuid guid生成
│       ├── uuid.go
│       └── uuid_test.go
├── LICENSE
├── log #基于zap的日志库
│   ├── zap.go
│   └── zap_test.go
├── natx #基于nat的请求/应答模式封装
│   ├── defaultApp.go
│   ├── defaultApp_test.go
│   ├── msgpack_enc.go
│   ├── natc.go
│   ├── natx.go
│   └── README.md
├── nsq #nsq的生产者和消费者封装
│   ├── consumer.go
│   └── producer.go
├── pool 
│   ├── allocator #内存分配器，使用sync.Pool，内部分配16种大小规格pool
│   │   ├── alloc.go
│   │   └── alloc_test.go
│   ├── xbufio #
│   │   └── buffio.go
│   ├── xbytes #固定大小的内存分配
│   │   ├── bytes.go
│   │   └── writer.go
│   └── xtime #定时器池 小根堆
│       └── xtime.go
├── README.md
├── redis #redis基础指令+乐观锁
│   ├── mutex.go
│   ├── redis.go
│   └── redis_test.go
└── time
    ├── minheap #定时器池 小根堆
    │   ├── timer.go
    │   └── timer_test.go
    └── wheel #定时器池 时间轮
        ├── timer.go
        └── timer_test.go
``` 