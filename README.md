# utils

```
├── codec # 一种数据混淆编码方法
│   ├── codec.go
│   ├── codec_test.go
│   └── README.md
├── container # 常用容器
│   ├── bitset
│   │   ├── bitset.go
│   │   └── bitset_test.go
│   ├── blocking_queue
│   │   ├── blocking_queue.go
│   │   └── blocking_queue_test.go
│   ├── bloom
│   │   ├── bloom9.go
│   │   └── bloom9_test.go
│   ├── fifo
│   │   ├── fifo.go
│   │   └── fifo_test.go
│   ├── omap
│   │   ├── omap.go
│   │   └── omap_test.go
│   └── pqueue
│       ├── priority_queue.go
│       ├── priority_queue_test.go
│       └── README.md
├── crypto # AES,DES加密解密
│   ├── aes256cbc.go
│   ├── aes256cbc_test.go
│   ├── aes.go
│   ├── PKCS.go
│   ├── tripledes_b_test.go
│   ├── tripledes.go
│   └── tripledes_test.go
├── etcdx # 基于etch的服务发行、选主、分布式锁
│   ├── discovery
│   │   ├── go.mod
│   │   ├── master.go
│   │   └── worker.go
│   ├── master
│   │   ├── go.mod
│   │   └── master.go
│   ├── readme.md
│   └── sync
│       ├── go.mod
│       ├── sync.go
│       └── sync_test.go
├── hash # 常用hash方法和一致性hash
│   ├── cityhash
│   │   ├── cityhash.go
│   │   └── cityhash_test.go
│   ├── hash.go
│   ├── ketama
│   │   ├── ketama.go
│   │   └── ketama_test.go
│   └── murmurhash3
│       ├── mmhash3.go
│       └── mmhash3_test.go
├── httpx # http sign
│   ├── go.mod
│   ├── README.md
│   ├── sign.go
│   └── sign_test.go
├── log # 封装基于zap.Logger的日志
│   ├── go.mod
│   ├── zap.go
│   └── zap_test.go
├── natx # nat使用封装
│   ├── defaultApp.go
│   ├── defaultApp_test.go
│   ├── go.mod
│   ├── msgpack_enc.go
│   ├── natc.go
│   ├── natx.go
│   └── README.md
├── nsqx # nsq使用封装
│   ├── consumer.go
│   ├── go.mod
│   └── producer.go
├── pool # 内存分配与优化
│   ├── allocator
│   │   ├── alloc.go
│   │   └── alloc_test.go
│   ├── xbufio
│   │   └── buffio.go
│   ├── xbytes
│   │   ├── bytes.go
│   │   └── writer.go
│   └── xtime # 基于内存小根堆定时器，扩展了timer func方法:func(interface{})
│       └── xtime.go
├── pprof # http pprof
│   ├── pprof_http.go
│   └── pprof_http_test.go
├── random # 随机字符串
│   └── string.go
├── README.md
├── redisx # redis常用封装和乐观锁
│   ├── go.mod
│   ├── mutex.go
│   ├── redis.go
│   └── redis_test.go
├── referral # 邀请码
│   ├── referral.go
│   └── refferal_test.go
├── sign # hmac签名
│   ├── hmac.go
│   └── hmac_test.go
├── timerx # 定时器
│   ├── minheap # 小根堆
│   │   ├── timer.go
│   │   └── timer_test.go
│   └── wheel # 时间轮
│       ├── timer.go
│       └── timer_test.go
└── validator # 参数验证
    └── validator.go
``` 