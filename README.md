# ZYCache : An Ordered-Set K-V memory cache service

震寅Cache是一个轻量级的内存缓存服务，提供类似于Go2.0提案中的参数列表化泛型支持的KeyStruct类型索引以及定时缓存键值对功能。震寅Cache分为两个主要部分，ZhenCache用于提供缓存服务，功能目前基本完善；YinLog目前是一个日志系统，采用无锁队列进行日志缓存服务，记录所有的数据集查询以及添加操作，后续会为缓存持久化提供支持。

有序键值对集合提供可选的两种数据结构——经典的红黑树以及Redis中采用的跳表，另外也提供了原生的Map作为K-V缓存底层，可以根据不同的场景进行选择。

调用服务十分便捷，无需任何初始化。使用该数据库缓存只需两步：

1. *首先定义自己的KeyStruct，实例中给出string的例子，目前支持string/int32/int64 三种类型的键，后续很容易进行扩展*

   ```go
   //继承原来的默认Key并重写参数列表中想要的内容
   type TestKey struct {
   	keystruct.DefaultKey
   	key string
   }
   
   //重写比较函数
   func (key TestKey) CompareBiggerThan(other keystruct.KeyStruct) bool {
   	return key.key > other.KeyString()
   }
   
   //重写作为主键以及在比较中需要用到的键的索引
   func (key TestKey) KeyString() string {
   	return key.key
   }
   ```

   

2. *其次直接进行调用，函数内部会开一个Goroutine进行插入与查询的处理，多次操作是幂等的，仅有一个数据库实例对多线程操作进行响应*

   ```go
   //Get方法的使用
   val, err := store.GetValue(TestKey{keystruct.DefaultKey{}, "test"})
   if err != nil {
   	handle_err()
   }
   handle_val(val)
   /***************************/
   //Set方法的使用
   //这个键会保持15分钟的缓存时间，过期后将无法获得
   store.Set(TestKey{keystruct.DefaultKey{},"test"},"value of testkey",15*time.Minute)
   ```

预定的一些进阶功能：

- 缓存持久化防止数据丢失
- 客户端程序的引入
- 提供分布式扩展支持

9/24 更新
添加了缓存持久化，优化了日志系统