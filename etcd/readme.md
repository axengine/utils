# sync
* sync Mutex为悲观锁，拿不到锁会一直阻塞
* sync New的时候会连接ETCD，失败返回nil；Lock时去set key，知道set成功或发生错误；UnLock时删除key；
* Lock是阻塞的，如果key一直没有expire或delete，会永远阻塞,因此set key时应正确设置ttl；
* 多个任务去获取同一把锁，会先后获取到锁

# Master
* Master采用锁机制申请为主，默认为slave，直到成功或发生错误才返回;
* 一但变为Master则一直为Master，知道程序coredump

# discovery