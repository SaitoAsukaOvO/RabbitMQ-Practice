# RabbitMQ

## Introduction：

MQ：在消息传输过程中保存信息的容器

```shell
Producer             Middleware      Consumer

System A  ------------> MQ -----------> System B
```

## Pros & Cons:

### Pros:

#### 应用解耦

直接远程调用：耦合性高，容错性低，可维护性低。

```shell
                             ---------------> System A
                             |
                             |
 user ------------------> MQ ---------------> System B
                             |
                             |
                             ---------------> System C
 mq接受到msg，就返回成功
 1. System A,B,C 只需要把msg拿出来去自己的系统里消费就可以
 2. if A or B or C fails，不影响其他系统，同时等系统恢复后再把msg从MQ中拿出来消费就可以 （容错性提高）
 3. 需要再加个System X，只需要X从MQ中拿出消息消费（可维护性提高）
```

#### 异步提速

```shell
                         --------300ms-------> System A
                         |
                         |
                         |
user ------------------->| -------300ms--------> System B
             |           |
             20ms        |
             |           |
             DB          |
                         --------300ms-------> System C
Latency = 20 + 300 + 300 + 300 = 920ms
```



    ```shell
                                --------5ms-------> System A (300ms)
                                |
                                |
                                |
    user --------|-----5ms-----> MQ -------5ms--------> System B (300ms)
                 20ms           |
                 |              |
                 DB             |
                                --------5ms-------> System C (300ms)
    Latency = 20 + 5 = 25ms
    ```



#### 削峰填谷

削峰：

```shell
                          ---------------> System A (QPS 1000)
                          |                   |
                          |                   |
                          |                   DB
 user ------------------->|
 
 瞬间请求增多：QPS=5000
```

```
                          ---------------> System A (QPS 1000) 从MQ中每秒拉取1000请求消费
                          |                   |
                          |                   |
                          |                   DB
 user ----MQ (5000 QPS)-->|
```



    填谷：
    
    积压的msg慢慢消费掉
    
    ==> **提高稳定性**

### Cons:

	可用性降低，复杂性提高，一致性问题

## RabbitMQ

### AMQP Components:

![Publish path from publisher to consumer via exchange and queue](https://www.rabbitmq.com/img/tutorials/intro/hello-world-example-routing.png)

**Exchange**：分发消息

**Queue**：存储消息

**Route**：该分发到哪个queue



### RabbitMQ Components:

![img](https://img2020.cnblogs.com/blog/1552936/202010/1552936-20201024103921637-693350551.png)

**Broker** : 消息队列服务 rabbitmq-server

**v-host** : `Virtual Host `虚拟主机。标识一批交换机、消息队列和相关对象。虚拟主机是共享相同的身份认证和加密环境的独立服务器域。每个vhost本质上就是一个mini版的RabbitMQ服务器，拥有自己的队列、交换器、绑定和权限机制。vhost是AMQP概念的基础，必须在链接时指定，RabbitMQ默认的vhost是 /。

**Binding** : 绑定，用于消息队列和交换机之间的关联。一个绑定就是基于路由键将交换机和消息队列连接起来的路由规则，所以可以将交换器理解成一个由绑定构成的路由表。

**Channel** : 信道，多路复用连接中的一条独立的双向数据流通道。信道是建立在真实的TCP连接内地虚拟链接，AMQP命令都是通过信道发出去的，不管是发布消息、订阅队列还是接收消息，这些动作都是通过信道完成。因为对于操作系统来说，建立和销毁TCP都是非常昂贵的开销，所以引入了信道的概念，以复用一条TCP连接。

**Connection** : 网络连接，比如一个TCP连接



## Practice:

see examples:
https://www.rabbitmq.com/getstarted.html

### Hello Word：

one producer one consumer:

![(P) -> [|||] -> (C)](https://www.rabbitmq.com/img/tutorials/python-one.png)



### Work Queue:

one producer multiple consumer (consumer为竞争关系)

![img](https://www.rabbitmq.com/img/tutorials/python-two.png)

#### 应用场景：

对于**任务过重**或者**任务较多**的情况使用work queue可以提高任务处理速度

#### Main Idea:

The main idea behind Work Queues (aka: *Task Queues*) is to avoid doing a resource-intensive task immediately and having to wait for it to complete. Instead we schedule the task to be done later. We encapsulate a *task* as a message and send it to a queue. A worker process running in the background will pop the tasks and eventually execute the job. When you run many workers the tasks will be shared between them

This concept is especially useful in web applications where it's impossible to handle a complex task during a short HTTP request window.

#### Realization:

We don't have a real-world task, like images to be resized or pdf files to be rendered, so let's fake it by just pretending we're busy - by using the `time.Sleep` function.

We'll take the number of dots in the string as its complexity; every dot will account for one second of "work". For example, a fake task described by `Hello...` will take three seconds.

#### Round-robin dispatching:

```shell
go run worker.go
go run worker.go
```

```shell
go run new_task.go First message.
go run new_task.go Second message..
go run new_task.go Third message...
go run new_task.go Fourth message....
go run new_task.go Fifth message.....
```



### Pub/Sub:

![img](https://www.rabbitmq.com/img/tutorials/exchanges.png)

#### exchange：

##### 	Fanout:

​		将消息交给所有绑定到的交换机队列

##### 	Direct:

​		把消息交给符合指定routing key的队列

##### 	Topic:

​		把消息交给符合routing pattern的队列

1. 接受P发送的消息
2. 知道如何处理消息，分发给某个特定队列，递交给所有队列，或者将消息丢弃

```shell
go run receive_log.go
go run receive_log.go > logs.log
go run send_log.go
```

#### 特性：

work queue: 很多个消费者监听同一个queue，只能有一个消费者收到

pub/sub: 很多消费者每个消费者监听自己的队列，消息来了之后，每个消费者都可以收到这个消息



### Routing：

![img](https://www.rabbitmq.com/img/tutorials/python-four.png)

只有queue的key和消息的key一致时，才会接收到消息

```shell
 go run receive_log_direct.go info warning error
 go run receive_log_direct.go warning error
 go run send_log_direct.go error "Run. Run. Or it will explode."
 go run send_log_direct.go
```



### Topic：

![img](https://www.rabbitmq.com/img/tutorials/python-five.png)



#### Topic exchange

Messages sent to a topic exchange can't have an arbitrary routing_key - it ***must be a list of words, delimited by dots.***

A few valid routing key examples: "stock.usd.nyse", "nyse.vmw", "quick.orange.rabbit". There can be as many words in the routing key as you like, ***up to the limit of 255 bytes***.

- ***\* (star) can substitute for exactly one word.***
- ***\# (hash) can substitute for zero or more words.***

These bindings in the Pic. can be summarised as:

- Q1 is interested in all the orange animals.

- Q2 wants to hear everything about rabbits, and everything about lazy animals.

  Message "quick.orange.rabbit" will be delivered to both queues.

  Message "lazy.orange.elephant" also will go to both of them.

  Message  "quick.orange.fox" will only go to the first queue, and "lazy.brown.fox" only to the second.

  Message "lazy.pink.rabbit" will be delivered to the second queue only once, even though it matches two bindings. "quick.brown.fox" doesn't match any binding so it will be discarded.

## Useful Links:

https://github.com/rabbitmq/internals

https://docs.vmware.com/en/VMware-Application-Catalog/services/tutorials/GUID-backup-restore-data-rabbitmq-kubernetes-index.html

https://tanzu.vmware.com/content/blog/kubernetes-tanzu-rabbitmq

