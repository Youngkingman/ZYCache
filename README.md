# ZYCache : An Ordered-Set K-V Memory Cache Service

震寅Cache是一个轻量级的内存缓存服务，提供类似于Go2.0提案中的参数列表化泛型支持的KeyStruct类型索引以及定时缓存键值对功能。震寅Cache分为两个主要部分，ZhenCache用于提供缓存服务，功能目前基本完善；YinLog目前是一个日志系统，采用无锁队列进行日志缓存服务，记录所有的数据集查询以及添加操作，2021/9/30实现了基于Json的持久化（待改进）以及快照+日志恢复。

有序键值对集合提供可选的两种数据结构——经典的红黑树以及Redis中采用的跳表，另外也提供了原生的Map以及带LRU缓存的Map作为K-V缓存底层，可以根据不同的场景进行选择。

## 作为单机缓存使用：

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



## 作为分布式缓存使用：

目前ZYCache提供的分布式缓存服务类似于memcache（或者Go语言专属的GroupCache），仅仅支持string Key,。提供了基于简单的RPC客户端和服务端程序

### 客户端调用：

作为任意一个客户端，可以通过直接调用.\ZYCache\zhenCache\rpcdef\clientcall.go的两个函数来远程使用服务，***需要注意的是目前仍然需要客户自己手动将value序列化为字符串，在读取之后进行反序列化***，后续将会进行改进，暂定采用protobuf进行远程通信，配合即将要写的RPC框架（咕咕咕）一起使用。severAddr参数指示了缓存服务器所在地址。

```go
//get some value by cli
//assert return interface as marshalled string
func Get(key string, serverAddr string) (interface{}, error) {
	....
}

//set some value by cli
//need value to be marshal
func Set(key string, value string, expire time.Duration, serverAddr string) error {
	....
}
```

### 服务端配置：

服务端采用了一致性哈希算法进行负载均衡，每一个服务器自带客户端，访问任意一个服务器都会正确地将键值对存储在合适的服务器上

<img src="E:\myproject\ZYCache\doc\1008_12.jpg" style="zoom:33%;" />

如图所示Client1访问Peer2服务端插入的键值对{"Zhenyin":1234}，Peer2对字符串"Zhenyin"经过一致性哈希计算被放置在Peer3服务器中;Client2访问Peer1查询键值{"Zhenyin"}，Peer1计算"Zhenyin"哈希值继而向Peer3获取该值。每一个Peer既是客户端也是服务端，相互之间的通信通过RPC调用实现，在一开始需要全部进行注册，实现上图的注册需要如下代码，放置于server.go中：

```go
func main() {
	var host string
	flag.IntVar(&host, "host",  "ZhenYin serve host")
	flag.Parse()

	addrMap := map[int]string{
		1: "10.0.8.22:9090",
		2: "10.0.8.18:9090",
		3: "10.0.8.33:9090",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	cood := rpcdef.New(host,"9090",addrs)
    cood.CoodinatorServe()
    select{}
}
```

Peer1的启动脚本,其他Peer简单替换一下host名即可

```shell
#!/bin/bash
trap "rm server;kill 0" EXIT

go build -o server
./server -host="10.0.8.22"
```
