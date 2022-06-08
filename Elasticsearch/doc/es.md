一、Elasticsearch 的历史由来
一个名叫 Shay Banon 的程序员为了给学习厨师的妻子构建一个食谱的搜素引擎，他开始构建一个早期版本的 Lucene，直接基于 Lucene 工作会比较困难，所以 Shay 开始抽象 Lucene 代码为了方便 JAVA 程序员可以在应用中添加搜索功能，他发布了它的第一个开源项目“Compass”。

后来 Shay 工作在高性能和内存数据网络的分布式环境中，因此高性能的、实时的、分布式搜索引擎是必需的，然后他重写 Compass 库称为一个独立的服务叫做 Elasticsearch。

二、Elasticsearch 简介
1、简介
Elasticsearch 是一个建立在 Apache Lucene™ 基础上的实时的分布式搜索和分析引擎，是基于 Lucene 实现的、当今最先进，最高效的全功能开源搜索引擎框架。Elasticsearch 使用 Lucene 作为内部索引引擎，而实际使用中，我们只需要使用统一开发好的 API 即可，而不需要理解其背后复杂的 Lucene 工作原理。
Elasticsearch 并不仅仅是基于 Lucene 实现全文搜索功能。同时，还具有以下特性：
1、分布式实时文件存储，并将每一个字段都编入索引，使其可以被搜索。
2、实时分析的分布式搜索引擎。
3、可以扩展到上百台服务器，处理 PB 级别的结构化或非结构化数据。
Elasticsearch 的使用很简单。初学者只要简单配置以一些参数，其他均使用默认值即可。而且安装也比较简单，安装好启动就可使用，可以很大幅度的减少学习成本。
如果你研究的够深入，你会发现 Elasticsearch 还有更多的高级功能，它可以很灵活地进行配置。可以根据自己的需求，灵活的来定制自己的 Elasticsearch。
2、官网：
https://www.elastic.co/cn/

三、ES 的分布式架构原理是什么（是如何实现分布式的）？
Elasticsearch 设计的理念就是分布式搜索引擎，底层其实还是基于 lucene 的。核心思想是在多台机器上启动多个 ES 进程实例，组成了一个 ES 集群。ES 中存储数据的基本单位是索引，如要在 ES 中存储一些订单数据，就应该在 ES 中创建一个索引，order_idx，所有的订单数据就都写到这个索引里面去，一个索引差不多就是相当于是 mysql 里的一张表。ES 的层级如下：index -> type -> mapping -> document -> field。
1、index：mysql 里的一张表
2、type：没法跟 mysql 里去对比，一个 index 里可以有多个 type，每个 type 的字段都是差不多的，但是有一些略微的差别。
好比说，有一个 index，是订单 index，里面专门是放订单数据的。如在 mysql 中建表，有些订单是实物商品的订单，如一件衣服，一双鞋子；有些订单是虚拟商品的订单，如游戏点卡，话费充值。就两种订单大部分字段是一样的，但是少部分字段可能有略微的一些差别。
所以就会在订单 index 里，建两个 type，一个是实物商品订单 type，一个是虚拟商品订单 type，这两个 type 大部分字段是一样的，少部分字段是不一样的。
很多情况下，一个 index 里可能就一个 type，但是如果说是一个 index 里有多个 type 的情况，你可以认为 index 是一个类别的表，具体的每个 type 代表了具体的一个 mysql 中的表。
3、mapping 就代表了这个 type 的表结构的定义，定义了这个 type 中每个字段名称，字段是什么类型的，然后还有这个字段的各种配置
4、实际上你往 index 里的一个 type 里面写的一条数据，叫做一条 document，一条 document 就代表了 mysql 中某个表里的一行，每个 document 有多个 field，每个 field 就代表了这个 document 中的一个字段的值

如上图，ES 中的一个索引可以拆分成多个 shard，每个 shard 存储部分数据。这个 shard 的数据实际是有多个备份，就是说每个 shard 都有一个 primary shard，负责写入数据，但是还有几个 replica shard。primary shard 写入数据之后，会将数据同步到其他几个 replica shard 上去。
通过这个 replica 的方案，每个 shard 的数据都有多个备份，如果某个机器宕机了，没关系啊，还有别的数据副本在别的机器上呢，这样就实现了高可用。
ES 集群多个节点，会自动选举一个节点为 master 节点，这个 master 节点其实就是干一些管理的工作的，比如维护索引元数据拉，负责切换 primary shard 和 replica shard 身份拉，之类的。要是 master 节点宕机了，那么会重新选举一个节点为 master 节点。
如果是非 master 节点宕机了，那么会由 master 节点，让那个宕机节点上的 primary shard 的身份转移到其他机器上的 replica shard。接着如果修复了那个宕机机器，重启之后，master 节点会控制将缺失的 replica shard 分配过去，同步后续修改的数据之类的，让集群恢复正常。

四、ES 写入、读取、查询、删除的原理是什么？
写数据过程
1）客户端选择一个 node 发送请求过去，这个 node 就变为 coordinating node（协调节点）
2）coordinating node，对 document 进行路由，将请求转发给对应的 node（有 primary shard）
3）实际的 node 上的 primary shard 处理请求，然后将数据同步到 replica node
4）coordinating node，如果发现 primary node 和所有 replica node 都搞定之后，就返回响应结果给客户端
读数据过程
查询，GET 某一条数据，写入了某个 document，这个 document 会自动给你分配一个全局唯一的 id，doc id，同时也是根据 doc id 进行 hash 路由到对应的 primary shard 上面去。也可以手动指定 doc id，比如用订单 id，用户 id。
1）客户端发送请求到任意一个 node，成为 coordinate node
2）coordinate node 对 document 进行路由，将请求转发到对应的 node，此时会使用 round-robin 随机轮询算法，在 primary shard 以及其所有 replica 中随机选择一个，让读请求负载均衡
3）接收请求的 node 返回 document 给 coordinate node
4）coordinate node 返回 document 给客户端
搜索数据过程
1）客户端发送请求到一个 coordinate node
2）协调节点将搜索请求转发到所有的 shard 对应的 primary shard 或 replica shard 也可以
3）query phase：每个 shard 将自己的搜索结果（其实就是一些 doc id），返回给协调节点，由协调节点进行数据的合并、排序、分页等操作，产出最终结果
4）fetch phase：接着由协调节点，根据 doc id 去各个节点上拉取实际的 document 数据，最终返回给客户端
写数据底层原理
1）先写入 buffer，在 buffer 里的时候数据是搜索不到的；同时将数据写入 translog 日志文件
2）如果 buffer 快满了，或者到一定时间（每隔 1 秒钟），就会将 buffer 数据 refresh 到一个新的 segment file 中，但是此时数据不是直接进入 segment file 的磁盘文件的，而是先进入 os cache 的。这个过程就是 refresh 。
为什么叫 es 是准实时的？NRT，near real-time，准实时。默认是每隔 1 秒 refresh 一次的，所以 es 是准实时的，因为写入的数据 1 秒之后才能被看到。
可以通过 es 的 restful api 或者 java api，手动执行一次 refresh 操作，就是手动将 buffer 中的数据刷入 os cache 中，让数据立马就可以被搜索到。
只要数据被输入 os cache 中，buffer 就会被清空了，因为不需要保留 buffer 了，数据在 translog 里面已经持久化到磁盘去一份了
3）只要数据进入 os cache，此时就可以让这个 segment file 的数据对外提供搜索了
4）重复 1~3 步骤，新的数据不断进入 buffer 和 translog，不断将 buffer 数据写入一个又一个新的 segment file 中去，每次 refresh 完 buffer 清空，translog 保留。随着这个过程推进，translog 会变得越来越大。当 translog 达到一定长度的时候，就会触发 commit 操作 。
5）commit 操作发生第一步，就是将 buffer 中现有数据 refresh 到 os cache 中去，清空 buffer
6）将一个 commit point 写入磁盘文件，里面标识着这个 commit point 对应的所有 segment file
7）强行将 os cache 中目前所有的数据都 fsync 到磁盘文件中去
8）将现有的 translog 清空，然后再次重启启用一个 translog，此时 commit 操作完成。默认每隔 30 分钟会自动执行一次 commit，但是如果 translog 过大，也会触发 commit。整个 commit 的过程，叫做 flush 操作。我们可以手动执行 flush 操作，就是将所有 os cache 数据刷到磁盘文件中去。
补充：
A、translog 日志文件的作用是什么？
在你执行 commit 操作之前，数据要么是停留在 buffer 中，要么是停留在 os cache 中，无论是 buffer 还是 os cache 都是内存，一旦这台机器死了，内存中的数据就全丢了。所以需要将数据对应的操作写入一个专门的日志文件，translog 日志文件中，一旦此时机器宕机，再次重启的时候，es 会自动读取 translog 日志文件中的数据，恢复到内存 buffer 和 os cache 中去。
B、commit 操作三步走？
写 commit point===>将 os cache 数据 fsync 强刷到磁盘上去===>清空 translog 日志文件
C、flush 与 commit 喵喵喵？
ES 中的 flush 操作，就对应着 commit 的全过程。我们也可以通过 es api，手动执行 flush 操作，手动将 os cache 中的数据 fsync 强刷到磁盘上去，记录一个 commit point，清空 translog 日志文件。
9）translog 其实也是先写入 os cache 的，默认每隔 5 秒刷一次到磁盘中去，所以默认情况下，可能有 5 秒的数据会仅仅停留在 buffer 或者 translog 文件的 os cache 中，如果此时机器挂了，会丢失 5 秒钟 的数据。但是这样性能比较好，最多丢 5 秒的数据。也可以将 translog 设置成每次写操作必须是直接 fsync 到磁盘，但是性能会差很多。
10）如果是删除操作，commit 的时候会生成一个.del 文件，里面将某个 doc 标识为 deleted 状态，那么搜索的时候根据.del 文件就知道这个 doc 被删除了
11）如果是更新操作，就是将原来的 doc 标识为 deleted 状态，然后新写入一条数据
12）buffer 每次 refresh 一次，就会产生一个 segment file，所以默认情况下是 1 秒钟一个 segment file，segment file 会越来越多，此时会定期执行 merge。当 segment file 多到一定程度的时候也会自动触发 merge 操作。
13）每次 merge 的时候，会将多个 segment file 合并成一个，同时这里会将标识为 deleted 的 doc 给物理删除掉，然后将新的 segment file 写入磁盘，这里会写一个 commit point，标识所有新的 segment file，然后打开 segment file 供搜索使用，同时删除旧的 segment file。
ES 里的写流程，有 4 个底层的核心概念， refresh、flush、translog、merge

五、ES 在数据量很大的情况下（数十亿级别）如何提高查询效率？
es 性能优化是没有什么银弹的，不要期待着随手调一个参数，就可以万能的应对所有的性能慢的场景。也许有的场景是你换个参数，或者调整一下语法，就可以搞定，但是绝对不是所有场景都可以这样。
（1）性能优化的杀手锏——filesystem cache

es 的搜索引擎严重依赖于底层的 filesystem cache，你如果给 filesystem cache 更多的内存，尽量让内存可以容纳所有的 indx segment file 索引数据文件，那么你搜索的时候就基本都是走内存的，性能会非常高。
要让 es 性能要好，最佳的情况下，就是你的机器的内存，至少可以容纳你的总数据量的一半，比如说，你一共要在 es 中存储 1T 的数据，那么你的多台机器留个 filesystem cache 的内存加起来综合，至少要到 512G，至少半数的情况下，搜索是走内存的，性能一般可以到几秒钟。
如果最佳的情况下，是仅仅在 es 中只存少量的数据，就是你要用来搜索的那些索引，内存留给 filesystem cache 的，就 100G，那么你就控制在 100G 以内，相当于是，你的数据几乎全部走内存来搜索，性能非常之高，一般可以在 1 秒以内。
比如说你现在有一行数据
id name age ….30 个字段
但是你现在搜索，只需要根据 id name age 三个字段来搜索
如果你傻乎乎的往 es 里写入一行数据所有的字段，就会导致 70%的数据是不用来搜索的，结果硬是占据了 es 机器上的 filesystem cache 的空间，单条数据的数据量越大，会导致 filesystem cahce 能缓存的数据就越少
仅仅只是写入 es 中要用来检索的少数几个字段就可以了，比如说，就写入 es id name age 三个字段就可以了，然后你可以把其他的字段数据存在 mysql 里面，一般是建议用 es + hbase 的这么一个架构。
hbase 的特点是适用于海量数据的在线存储，从 es 中根据 name 和 age 去搜索，拿到的结果可能就 20 个 doc id，然后根据 doc id 到 hbase 里去查询每个 doc id 对应的完整的数据，给查出来，再返回给前端。
（2）数据预热
对于那些你觉得比较热的，经常会有人访问的数据，最好做一个专门的缓存预热子系统，就是对热数据，每隔一段时间，你就提前访问一下，让数据进入 filesystem cache 里面去。这样期待下次别人访问的时候，一定性能会好一些。
举个例子，就比如说，微博，你可以把一些大 v，平时看的人很多的数据给提前你自己后台搞个系统，每隔一会儿，你自己的后台系统去搜索一下热数据，刷到 filesystem cache 里去，后面用户实际上来看这个热数据的时候，他们就是直接从内存里搜索了，很快。
电商，你可以将平时查看最多的一些商品，比如说 iphone XS max，热数据提前后台搞个程序，每隔 1 分钟自己主动访问一次，刷到 filesystem cache 里去。
（3）冷热分离
关于 ES 性能优化，数据拆分，之前说将大量不搜索的字段，拆分到别的存储中去，类似于 mysql 分库分表的垂直拆分。
ES 也可以做类似于 mysql 的水平拆分，将大量的访问很少，频率很低的数据，单独写一个索引，然后将访问很频繁的热数据单独写一个索引。这样可以确保热数据在被预热之后，尽量都让他们留在 filesystem os cache 里，不让冷数据给冲刷掉。
（4）document 模型设计
ES 里面的复杂的关联查询，复杂的查询语法（join 之类的），尽量别用，一旦用了性能一般都不太好。
很多复杂的乱七八糟的一些操作，如何执行
两个思路，在搜索/查询的时候，要执行一些业务强相关的特别复杂的操作：
1）在写入数据的时候，就设计好模型，加几个字段，把处理好的数据写入加的字段里面
2）自己用 java 程序封装，es 能做的，用 es 来做，搜索出来的数据，在 java 程序里面去做，用 java 封装一些特别复杂的操作
比如，如果一定要执行 join 操作，可以采用如下策略：
写入 es 的时候，搞成两个索引，order 索引，orderItem 索引。order 索引，里面就包含 id order_code total_price。orderItem 索引，里面写入进去的时候，就完成 join 操作，id order_code total_price id order_id goods_id purchase_count price。
即写入 es 的 java 系统里，就完成关联，将关联好的数据直接写入 es 中，搜索的时候，就不需要利用 es 的搜索语法去完成 join 来搜索了。
（5）分页性能优化
es 的分页是较坑的，为啥呢？举个例子吧，假如你每页是 10 条数据，你现在要查询第 100 页，实际上是会把每个 shard 上存储的前 1000 条数据都查到一个协调节点上，如果你有个 5 个 shard，那么就有 5000 条数据，接着协调节点对这 5000 条数据进行一些合并、处理，再获取到最终第 100 页的 10 条数据。
翻页的时候，翻的越深，每个 shard 返回的数据就越多，而且协调节点处理的时间越长。非常坑爹。所以用 es 做分页的时候，会发现越翻到后面，就越是慢。
那么该如何解决呢？
1）不允许深度分页或者告诉用户/PM 默认深度分页性能很惨
2）类似于 app 里的推荐商品不断下拉出来一页一页的，可以用 scroll api
scroll 会一次性给你生成所有数据的一个快照，然后每次翻页就是通过游标移动，获取下一页下一页这样子，性能会比上面说的那种分页性能也高很多很多，无论翻多少页，性能基本上都是毫秒级的。
但是唯一的一点就是，这个适合于那种类似微博下拉翻页的，不能随意跳到任何一页的场景。同时这个 scroll 是要保留一段时间内的数据快照的，你需要确保用户不会持续不断翻页翻几个小时。
因为 scroll api 是只能一页一页往后翻的，是不能说，先进入第 10 页，然后去 120 页，回到 58 页，不能随意乱跳页。所以现在很多产品，都是不允许你随意翻页的，app，也有一些网站，做的就是你只能往下拉，一页一页的翻。

六、Elasticsearch 安装
https://www.cnblogs.com/tianyiliang/p/10291305.html

七、可视化工具
可视化工具比较多，各有优缺点，选择适合自己的就可以，目前用的比较多的是 Kibana
1、 Kibana
下载
https://www.elastic.co/cn/downloads/past-releases/kibana-6-8-0
安装
https://blog.csdn.net/weixin_34727238/article/details/81200071
2、 Head
Linux: https://blog.csdn.net/sinat_37690778/article/details/78905390
Windows ： https://blog.csdn.net/tzconn/article/details/83016494
https://blog.csdn.net/u013456370/article/details/79608365
3、ElasticHD
https://github.com/360EntSecGroup-Skylar/ElasticHD

八、SpringCloud 接入 Elastic Search
1、pom.xml
<dependencies>
<dependency>
<groupId>org.springframework.boot</groupId>
<artifactId>spring-boot-starter</artifactId>
</dependency>

 <dependency>
 <groupId>org.springframework.boot</groupId>
 <artifactId>spring-boot-starter-test</artifactId>
 <scope>test</scope>
 </dependency>

 <!--rest-->
 <dependency>
 <groupId>org.elasticsearch.client</groupId>
 <artifactId>elasticsearch-rest-client</artifactId>
 <version>6.4.0</version>
 </dependency>
 <dependency>
 <groupId>org.elasticsearch.client</groupId>
 <artifactId>elasticsearch-rest-high-level-client</artifactId>
 <version>6.4.0</version>
 </dependency>

 </dependencies>
2、application.properties
spring.elasticsearch.rest.uris= http://127.0.0.1:9200

3、使用
@Resource
private RestHighLevelClient restHighLevelClient;

九、es 数据类型
1、字段类型概述

一级分类 二级分类 具体类型
核心类型 字符串类型 string,text,keyword
整数类型 integer,long,short,byte
浮点类型 double,float,half_float,scaled_float
逻辑类型 boolean
日期类型 date
范围类型 range
二进制类型 binary
复合类型 数组类型 array
对象类型 object
嵌套类型 nested
地理类型 地理坐标类型 geo_point
地理地图 geo_shape
特殊类型 IP 类型 ip
范围类型 completion
令牌计数类型 token_count
附件类型 attachment
抽取类型 percolator

2、字符串类型

（1）string
string 类型在 ElasticSearch 旧版本中使用较多，从 ElasticSearch 5.x 开始不再支持 string，由 text 和 keyword 类型替代。
（2）text
当一个字段是要被全文搜索的，比如 Email 内容、产品描述，应该使用 text 类型。设置 text 类型以后，字段内容会被分析，在生成倒排索引以前，字符串会被分析器分成一个一个词项。text 类型的字段不用于排序，很少用于聚合。
（3）keyword
keyword 类型适用于索引结构化的字段，比如 email 地址、主机名、状态码和标签。如果字段需要进行过滤(比如查找已发布博客中 status 属性为 published 的文章)、排序、聚合。keyword 类型的字段只能通过精确值搜索到。
3、整数类型

类型 取值范围
byte -128~127
short -32768~32767
integer -2^31~2^31-1
long -2^63~2^63-1
在满足需求的情况下，尽可能选择范围小的数据类型。比如，某个字段的取值最大值不会超过 100，那么选择 byte 类型即可。迄今为止吉尼斯记录的人类的年龄的最大值为 134 岁，对于年龄字段，short 足矣。字段的长度越短，索引和搜索的效率越高。
4、浮点类型

类型 取值范围
doule 64 位双精度 IEEE 754 浮点类型
float 32 位单精度 IEEE 754 浮点类型
half_float 16 位半精度 IEEE 754 浮点类型
scaled_float 缩放类型的的浮点数
对于 float、half_float 和 scaled_float,-0.0 和+0.0 是不同的值，使用 term 查询查找-0.0 不会匹配+0.0，同样 range 查询中上边界是-0.0 不会匹配+0.0，下边界是+0.0 不会匹配-0.0。
其中 scaled_float，比如价格只需要精确到分，price 为 57.34 的字段缩放因子为 100，存起来就是 5734
优先考虑使用带缩放因子的 scaled_float 浮点类型。
5、boolean 类型

逻辑类型（布尔类型）可以接受 true/false/”true”/”false”值
PUT test
{
"mappings" :{
"my" :{
"properties" :{
"empty" :{
"type" : "boolean"
}
}
}
}
}
6、date 类型

我们人类使用的计时系统是相当复杂的：秒是基本单位, 60 秒为 1 分钟, 60 分钟为 1 小时, 24 小时是一天……如果计算机也使用相同的方式来计时, 那显然就要用多个变量来分别存放年月日时分秒, 不停的进行进位运算, 而且还要处理偶尔的闰年和闰秒以及协调不同的时区. 基于”追求简单”的设计理念, UNIX 在内部采用了一种最简单的计时方式：
计算从 UNIX 诞生的 UTC 时间 1970 年 1 月 1 日 0 时 0 分 0 秒起, 流逝的秒数.
UTC 时间 1970 年 1 月 1 日 0 时 0 分 0 秒就是 UNIX 时间 0, UTC 时间 1970 年 1 月 2 日 0 时 0 分 0 秒就是 UNIX 时间 86400.
这个计时系统被所有的 UNIX 和类 UNIX 系统继承了下来, 而且影响了许多非 UNIX 系统.
日期类型表示格式可以是以下几种：
（1）日期格式的字符串，比如 “2018-01-13” 或 “2018-01-13 12:10:30”
（2）long 类型的毫秒数( milliseconds-since-the-epoch，epoch 就是指 UNIX 诞生的 UTC 时间 1970 年 1 月 1 日 0 时 0 分 0 秒)
（3）integer 的秒数(seconds-since-the-epoch)
ElasticSearch 内部会将日期数据转换为 UTC，并存储为 milliseconds-since-the-epoch 的 long 型整数。
PUT test
{
"mappings" :{
"my" :{
"properties" :{
"postdate" :{
"type" : "date" ,
"format" : "yyyy-MM-ddHH:mm:ss||yyyy-MM-dd||epoch_millis"
}
}
}
}
}
7、binary 类型
二进制字段是指用 base64 来表示索引中存储的二进制数据，可用来存储二进制形式的数据，例如图像。默认情况下，该类型的字段只存储不索引。二进制类型只支持 index_name 属性。

8、array 类型
在 ElasticSearch 中，没有专门的数组（Array）数据类型，但是，在默认情况下，任意一个字段都可以包含 0 或多个值，这意味着每个字段默认都是数组类型，只不过，数组类型的各个元素值的数据类型必须相同。在 ElasticSearch 中，数组是开箱即用的（out of box），不需要进行任何配置，就可以直接使用。
在同一个数组中，数组元素的数据类型是相同的，ElasticSearch 不支持元素为多个数据类型：[ 10, “some string” ]，常用的数组类型是：
（1）字符数组: [ “one”, “two” ]
（2）整数数组: productid:[ 1, 2 ]
（3）对象（文档）数组: “user”:[ { “name”: “Mary”, “age”: 12 }, { “name”: “John”, “age”: 10 }]，ElasticSearch 内部把对象数组展开为 {“ user.name ”: [“Mary”, “John”], “user.age”: [12,10]}

9、object 类型
JSON 天生具有层级关系，文档会包含嵌套的对象
PUT test
PUT test/my/1
{
"employee" :{
"age" : 30 ,
"fullname" :{
"first" : "hadron" ,
"last" : "cheng"
}
}
}
上面文档整体是一个 JSON，JSON 中包含一个 employee,employee 又包含一个 fullname
GET test/\_mapping
{
"test" :{
"mappings" :{
"my" :{
"properties" :{
"employee" :{
"properties" :{
"age" :{
"type" : "long"
},
"fullname" :{
"properties" :{
"first" :{
"type" : "text" ,
"fields" :{
"keyword" :{
"type" : "keyword" ,
"ignore_above" : 256
}
}
},
"last" :{
"type" : "text" ,
"fields" :{
"keyword" :{
"type" : "keyword" ,
"ignore_above" : 256
}
}
}
}
}
}
}
}
}
}
}
}

十、API
