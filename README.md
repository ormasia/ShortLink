# 短链接项目：

- 核心思路：开发的一个长链接转短链接的通用组件，支持将网址长链接转为短链接，通过短链接跳转到原链接的项目。

- 希望用到的组件有mysql，redis，kafka，nginx
- 用到的框架有gin，gorm
- 采用Saas方式，后面希望可以使用上consul，etcd，nacos，grpc支持微服务
- 后期目标是使用docker，k8s做一个实际上线的部署

- 计划用时：一周

## 采用GPT4.5生成的框架，可供参考
MINI短链项目架构设计

一、架构概述

MINI短链项目采用基于 Golang 的微服务架构，结合缓存、数据库、负载均衡和高效的过滤机制，实现了高性能、高可用的短链接服务。

二、技术栈

后端语言：Golang

数据库：MySQL

缓存服务：Redis

布隆过滤器：防止缓存穿透

Singleflight：防止缓存击穿

负载均衡器：Nginx

三、数据库设计

发号器表：实现短链接ID生成，高并发下保证唯一性和连续性。

字段：id（主键、自增）、创建时间、更新时间。

长短链接映射表：记录长链接与短链接的映射关系。

字段：short_url（主键）、original_url、创建时间、过期时间、访问次数。

四、服务架构设计

负载均衡层

使用Nginx作为反向代理，负载均衡多个服务实例，提升系统的并发处理能力和可靠性。

API服务层

统一的HTTP接口对外提供服务：

转链接接口：POST /shorten

查链接接口：GET /:short_url

业务逻辑层

转链接服务：

特词过滤：检查链接合法性，避免敏感或非法链接。

循环转链检测：防止URL重复或循环重定向。

生成短链接ID：通过MySQL发号器服务进行ID生成。

存储映射关系：长链接与生成的短链接存入MySQL数据库并同步到缓存。

查链接服务：

查询缓存：首先使用Redis缓存查询，提升响应速度。

缓存穿透防护：布隆过滤器检测短链接是否存在，防止大量无效查询击穿数据库。

缓存击穿防护：使用Singleflight合并高并发请求，避免大量请求同时穿透到数据库。

缓存层

Redis缓存存储短链接与长链接的映射关系，减轻数据库压力。

布隆过滤器缓存有效短链接标识，避免恶意或无效请求。

数据层

MySQL实现高可靠的数据存储。

通过数据库事务与主键生成确保数据一致性和服务高可用性。

五、性能优化

布隆过滤器：避免缓存穿透，减少对数据库无效查询的压力。

Singleflight：避免缓存击穿，防止高并发条件下重复数据库请求。

Redis缓存：提升热点链接访问性能。

六、部署方案

容器化部署（Docker + Kubernetes），实现服务水平扩展和高可用。

日志和监控：Prometheus + Grafana实现服务监控和报警。

架构优势

高性能：多层缓存和过滤机制有效提升系统性能。

高可用：微服务架构下，单服务异常不影响整体服务。

安全性高：特词过滤和循环转链检测保证系统链接安全性。

易维护扩展：微服务设计便于服务的快速迭代和水平扩展。


### 创建topic
``` bash
kafka-topics.sh --create \
  --bootstrap-server localhost:9092 \
  --replication-factor 1 \
  --partitions 3 \
  --topic shortlink-log
```
``` docker
docker-compose.yml
environment:
  KAFKA_CREATE_TOPICS: "shortlink-log:3:1"
```

### nacos启动问题
  必须要把grpc的端口暴露出来
  主要原因是因为调用 configClient.GetConfig方法的时候会访问grpc服务，nacos2添加了grpc通信方式，所以需要把grpc的端口也打开

  docker启动的时候记得把9848和9849暴露出来，也就是把grpc打开