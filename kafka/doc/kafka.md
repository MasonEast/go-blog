## 基本概念

### 主题层

`Topic` 是用来存储生产者(`producer`)生产的消息事件，是由一个或多个 `partition` 组成的。

![image-20220602104413408](https://raw.githubusercontent.com/MasonEast/go-blog/master/kafka/doc/%E4%B8%BB%E9%A2%98%E4%B8%8E%E5%88%86%E5%8C%BA.png)

`offset`是用来标记一个消费者在某个`partition`上读到了那一条消息。

### 分区层

在实际应用中，我们往往将`partition`分配在不同的磁盘上，利用多磁盘来增加多写效率。既然是分布式，必然会有多个机器，而一个机器，我们称为一个`broker`（节点）

> tips: 这里的机器不一定指的是物理机器，多节点不一定要在不同的机器上

![image-20220602104413408](https://raw.githubusercontent.com/MasonEast/go-blog/master/kafka/doc/broker%E4%B8%8E%E9%9B%86%E7%BE%A4.png)

每个 broker 都有一套冗余数据，称为 repliaction(副本)。

### zookeeper

假如我们有三个`broker`，客户端不知道该和哪个节点连接，这就需要`zookeeper`选取一个`leader`，剩下的两个为`follower`，生产者和消费者只和`leader`交互，`follower`负责跟随，同步`leader`的消息，等待`leader`挂了它顶上去。

`zookeeper`可以理解为`kafka`的管家，它负责管理所有的`broker`的`IP`地址。

### 消息层

这一层主要是存储信息和消费者`consumer`的`offset`。

## docker 启动 kafka

### 配置 yml

```yml
version: "3.8"
services:
  zookeeper:
    container_name: zookeeper
    image: wurstmeister/zookeeper
    ports:
      - "2181:2181"
  kafka:
    container_name: kafka
    image: wurstmeister/kafka
    depends_on: [zookeeper]
    ports:
      - "9092:9092"
    environment:
      KAFKA_ADVERTISED_HOST_NAME: 127.0.0.1
      KAFKA_ADVERTISED_PORT: 9092
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
```

通过`docker-compose up`启动。

### bash 测试

启动两个 bash

- bash1

```bash
#创建单分区单副本的 topic demo：
bin/kafka-topics.sh --create --zookeeper zookeeper_apoint:2181 --replication-factor 1 --partitions 1 --topic demo

#查看 topic 列表：
bin/kafka-topics.sh --list --zookeeper zookeeper_apoint:2181

#发送消息【生产者】
bin/kafka-console-producer.sh --broker-list kafka_apoint:9092 --topic demo

#查看描述 topics 信息
bin/kafka-topics.sh --describe --zookeeper zookeeper_apoint:2181 --topic demo

Topic:demo      PartitionCount:1        ReplicationFactor:1     Configs:
        Topic: demo     Partition: 0    Leader: 1       Replicas: 1     Isr: 1


# 第一行给出了所有分区的摘要，每个附加行给出了关于一个分区的信息。 由于我们只有一个分区，所以只有一行。

# “Leader”: 是负责给定分区的所有读取和写入的节点。
# “Replicas”: 是复制此分区日志的节点列表，无论它们是否是领导者，或者即使他们当前处于活动状态。
# “Isr”: 是一组“同步”副本。这是复制品列表的子集，当前活着并被引导到领导者

```

- bash2

```bash
#接收消息并在终端打印：
bin/kafka-console-consumer.sh --bootstrap-server kafka_apoint:9092 --topic demo --from-beginning
```

### 外网访问

配置：server.properties

```bash
# 修改成内网地址
listeners=PLAINTEXT://0.0.0.0:9092
# 改成外网地址
advertised.listeners=PLAINTEXT://
```
