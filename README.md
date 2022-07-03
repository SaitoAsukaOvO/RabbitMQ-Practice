# RabbitMQ

## Introduction：

MQ：在消息传输过程中保存信息的容器

```shell
Producer             Middleware      Consumer

System A  ------------> MQ -----------> System B
```

## Pros & Cons:

### Pros:

1. 应用解耦

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

2. 异步提速

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



3. 削峰填谷

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

​	可用性降低，复杂性提高，一致性问题

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

## Useful Reference:

https://github.com/rabbitmq/internals

https://docs.vmware.com/en/VMware-Application-Catalog/services/tutorials/GUID-backup-restore-data-rabbitmq-kubernetes-index.html

https://tanzu.vmware.com/content/blog/kubernetes-tanzu-rabbitmq

