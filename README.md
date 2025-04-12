# 短链接项目：

- 核心思路：开发的一个长链接转短链接的通用组件，支持将网址长链接转为短链接，通过短链接跳转到原链接的项目。

- 希望用到的组件有mysql，redis，kafka，nginx
- 用到的框架有gin，gorm
- 采用Saas方式，后面希望可以使用上consul，etcd，nacos，grpc支持微服务
- 后期目标是使用docker，k8s做一个实际上线的部署


## 采用GPT4.5生成的框架，cursor和trae开发
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

使用Gin框架实现API服务，采用restful风格。

~~统一的HTTP接口对外提供服务：~~

~~转链接接口：POST /shorten~~

~~查链接接口：GET /:short_url~~

业务逻辑层

转链接服务：

特词过滤：检查链接合法性，避免敏感或非法链接。

循环转链检测：防止URL重复或循环重定向。

~~生成短链接ID：通过MySQL发号器服务进行ID生成~~。

生成短链接ID：0~7	时间戳的二进制值（防碰撞）
            8~15	高熵随机数（防重复）
            生成短链 key（Base62） ← 时间戳 + 随机数（共16字节） → 取前 N 位 → 冲突检查（checkExists）→ 返回

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

容器化部署（Docker ~~+ Kubernetes~~），实现服务水平扩展和高可用。

~~日志和监控：Prometheus + Grafana实现服务监控和报警。~~

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

``` bash
protoc --go_out=. --go-grpc_out=. proto/shortlinkpb/shortlink.proto
```

### nacos启动问题
  必须要把grpc的端口暴露出来
  主要原因是因为调用 configClient.GetConfig方法的时候会访问grpc服务，nacos2添加了grpc通信方式，所以需要把grpc的端口也打开

  docker启动的时候记得把9848和9849暴露出来，也就是把grpc打开

 TODO: 预生成短链池 + 异步批量写入数据库，减少实时生成压力

  使用 Redis 缓存热点长-短映射，降低数据库查询频率

  全量预热布隆：从数据库加载最近有效 key

分布式锁控制并发：使用 Redis 分布式锁（锁粒度为 URL hash）避免重复生成和数据库写入。

随机化短链：结合时间戳 + 随机因子，基于 Base62 编码生成短链，具备较强冲突熵，保障短链分布均匀。

布隆过滤器：在跳转时先通过布隆过滤器判断短链是否存在，防止非法请求穿透缓存打到数据库。

SingleFlight 并发合并：对并发跳转请求使用 Go 的 singleflight，合并同一时刻对同一短链的查询，防止缓存击穿。

缓存雪崩防护：所有 Redis 缓存设置随机过期时间（如 10min + rand(1~30s)），避免集中失效造成数据库雪崩。

## 测试shorten

### 第一次
  
- 最多同时并发 10 个请求
  总请求数: 500 
  成功短链数: 374
  唯一短链数: 374
  总耗时: 8.66 秒
  平均 QPS: 57.73

- 最多同时并发 100 个请求
  总请求数: 500
  成功短链数: 291
  唯一短链数: 291
  总耗时: 1.73 秒
  平均 QPS: 288.45

### 第二次
修改了生成id的方式，使用了时间戳的二进制值和随机数的方式，避免了碰撞和重复的问题
测试结果时：

- 最多并发200个请求
  总请求数: 1000
  成功短链数: 999
  唯一短链数: 999
  总耗时: 3.01 秒
  平均 QPS: 331.78
- 最多并发100个请求
  请求数: 10000
  成功短链数: 9923
  唯一短链数: 9923
  总耗时: 58.02 秒
  平均 QPS: 172.36

### 分布式锁后面不查询数据库，直接返回
- 最多并发50个请求
总请求数: 1348
成功短链数: 1000
唯一短链数: 1000
总耗时: 2.48 秒
平均 QPS: 543.28

- 最多并发200个请求
总请求数: 1348
成功短链数: 1002
唯一短链数: 1002
总耗时: 3.23 秒
平均 QPS: 417.05

超过200的并发才会出现重复的短链，50以内不会出现重复的短链，考虑分布式锁之后不加数据库查询，直接返回

**你完全可以在 Gin 框架层做限流**，特别是在你目前还没有引入 Nginx 或 Kong 这类独立网关组件的情况下。

---

## ✅ 一、Gin 层限流的适用场景

在你的架构中，Gin 承担了“网关 + 路由 + API业务”的多重角色。**在没有外部独立网关时，在 Gin 上做限流是合理且常见的选择。**

适合的使用场景包括：

- **防止单个用户/接口恶意刷请求**
- **保护后端服务资源（如数据库/缓存）**
- **确保高峰期服务稳定性**

---

## 🧰 二、Gin 框架实现限流的方式

### 方式一：使用中间件 + 令牌桶算法（推荐）

你可以用 Go 的 [`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate) 实现一个 **基于令牌桶的限流中间件**，按 IP / 用户 / 接口维度限流。

#### 示例：每秒只允许 5 个请求
```go
import (
	"time"
	"golang.org/x/time/rate"
	"github.com/gin-gonic/gin"
)

var limiter = rate.NewLimiter(5, 10) // 速率: 每秒5个请求，突发容量10个

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"message": "Too Many Requests"})
			return
		}
		c.Next()
	}
}
```

注册到 Gin：
```go
r := gin.Default()
r.Use(RateLimitMiddleware())
```

---

### 方式二：使用第三方库（封装更好）

#### 推荐库：`github.com/ulule/limiter/v3`

支持：
- 内存、Redis、PostgreSQL 存储
- 灵活的限流策略（固定窗口、滑动窗口、漏桶等）
- 多维度限流（IP、Header、用户Token）

GitHub: https://github.com/ulule/limiter

---

### 方式三：自己实现简单的滑动窗口计数器

如果你希望对某些接口做用户级别限流（例如 `/api/shorten` 每个用户每分钟最多调用5次），也可以通过一个内存Map或Redis实现滑动窗口计数逻辑。

---

## 🧠 三、Gin 层限流 vs 网关层限流（Nginx）

| 维度             | Gin 层限流                     | Nginx / Kong 网关限流                   |
|------------------|--------------------------------|-----------------------------------------|
| 实现灵活度       | ✅ 非常高（可按用户/IP/接口定制） | ❌（Nginx 规则配置较死板）               |
| 执行位置         | 应用层（业务代码中）            | 网络入口层（还未进入应用）               |
| 性能开销         | 较高，进了应用再判断            | 较低，提前拒绝                          |
| 可观测性/监控     | 手动实现                       | Kong/APISIX 等自带监控和仪表盘          |
| 与RBAC等功能集成 | ✅ 容易一起封装在中间件里         | ❌ 难以访问业务Token/权限信息           |

---

## 🧩 四、建议策略（你的项目现阶段）

| 项目阶段       | 限流建议                                      |
|----------------|-----------------------------------------------|
| 当前阶段（Gin为主） | ✅ 使用 Gin 中间件方式实现限流                  |
| 后期接入网关（Nginx/Kong） | ✅ 推荐将基础限流（如IP限流）前置到网关，业务限流留在Gin |

---

## 🚀 五、加分项：基于用户限流 + Redis 持久化

如果你支持登录、Token认证，可实现基于用户的限流策略。例如：

- `userID = 123`：每分钟创建短链上限为10
- 存储结构：`SETEX rate_limit:user:123 60 <count>`

这样一来，每个用户的限流规则都独立控制，不会被其他用户影响。

---

1. Gin 限流中间件代码（支持 IP / 用户 ID / 路径维度）
2. 基于 Redis 的高可用限流实现
3. 各种限流策略组合的配置样例


[架构](https://i-blog.csdnimg.cn/img_convert/e0b5e75d945af7f21f8f94f58fad5cbc.png)

[功能](https://i-blog.csdnimg.cn/img_convert/98fa3ad51bdfcac86752312c93272b4d.png)

### api网关限流
全局限流（RateLimitMiddleware）：
  使用令牌桶算法
  基于客户端 IP 进行限流
  每秒最多处理 100 个请求
  适用于所有接口
批量创建限流（BatchRateLimitMiddleware）：
  使用滑动窗口算法
  基于用户 ID 进行限流
  每分钟最多处理 10 个批量请求
  专门用于 /api/v1/links/batch 接口
这样的限流策略可以：
  防止单个 IP 的恶意请求
  保护批量创建接口不被滥用
  确保系统在高并发下的稳定性
您可以根据实际需求调整以下参数：
  MaxRequests：每个时间窗口允许的最大请求数
  WindowSize：时间窗口大小
  Rate：令牌产生速率
  Capacity：令牌桶容量

  是的，在 Gin 框架中获取客户端 IP 地址需要注意以下几点：

1. `c.ClientIP()` 方法会按以下顺序尝试获取真实 IP：
   - X-Forwarded-For 头
   - X-Real-IP 头
   - RemoteAddr

但是，由于您的应用可能部署在代理（如 Nginx）后面，我们需要确保获取到真实的客户端 IP。让我修改一下限流中间件的实现：


抱歉，让我重新组织一下回答。

在 Gin 框架中获取客户端 IP 地址时，需要考虑以下几种情况：

1. 直接访问（无代理）：
   - 使用 `c.RemoteIP()` 可以获取客户端 IP

2. 通过代理访问（如 Nginx）：
   - `X-Forwarded-For`: 包含请求经过的所有代理 IP，格式为 `client, proxy1, proxy2`
   - `X-Real-IP`: 通常由代理服务器设置的真实客户端 IP

建议的获取真实 IP 的逻辑如下：

```go
// 获取真实客户端IP
func getRealIP(c *gin.Context) string {
    // 尝试从 X-Forwarded-For 获取
    xForwardedFor := c.GetHeader("X-Forwarded-For")
    if xForwardedFor != "" {
        // X-Forwarded-For 格式为：client, proxy1, proxy2
        // 第一个 IP 是真实客户端 IP
        ips := strings.Split(xForwardedFor, ",")
        if len(ips) > 0 {
            return strings.TrimSpace(ips[0])
        }
    }

    // 尝试从 X-Real-IP 获取
    xRealIP := c.GetHeader("X-Real-IP")
    if xRealIP != "" {
        return xRealIP
    }

    // 如果都没有，使用 RemoteIP
    return c.RemoteIP()
}
```

如果您使用 Nginx 作为反向代理，需要在 Nginx 配置中添加以下配置来传递客户端 IP：

```nginx
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
```

这样可以确保在各种部署场景下都能正确获取到真实的客户端 IP。

您是否需要我帮您：
1. 修改 Nginx 配置文件？
2. 添加更多的 IP 地址验证逻辑（如检查 IP 格式）？
3. 实现 IP 黑名单功能？


## 注销用户
基本符合一个**可靠用户注销流程**的最佳实践。我们可以把你的五个步骤稍微补充整理一下，形成更系统的执行流程：

``` text
注销用户流程：
│
├── 1. 删除用户（数据库操作，标记为已注销 / 删除行）
│
├── 2. 删除 Token 缓存（如 Redis 中的 token-key）
│
├── 3. 异步删除用户短链（Kafka 发布注销事件）
│       └─ 消费方监听事件 → 删除短链
│
├── 4. 异步操作失败处理（Kafka失败重试 + 死信队列 DLQ）
│
└── 5. 补偿机制（定时任务扫描未删除成功记录，重试）
```

### ✅ 用户注销操作推荐五步法：

#### **1️⃣ 删除用户主数据（数据库）**

- 逻辑删除或硬删除用户表中的记录
- 若有用户状态字段（如 `status`），可置为 `deleted`，避免误删
- 若涉及多表外键引用，推荐**软删除**

```go
db.Model(&User{}).Where("id = ?", uid).Update("status", "deleted")
```

---

#### **2️⃣ 清除登录状态（Token 缓存）**

- 删除 Redis 缓存中的 JWT Token
- 可按用户 ID 删除所有 token（例如用 Set 存储该用户所有 token）
- 防止注销后还能访问接口

```go
cache.Del("token:{token}")
```

---

#### **3️⃣ 异步删除资源（Kafka 通知短链服务删除用户短链）**

- **Kafka 发布消息**（Topic: `user-deleted`）
- **短链服务监听并删除该用户的短链记录**

> 为什么用 Kafka？
> - 异步解耦，提高性能
> - 失败后可以重试 / 进 DLQ（死信队列）

---

#### **4️⃣ 记录失败事件（用于重试或报警）**

- 消费失败的记录（如短链删除失败）要打日志或写入失败表
- 日志信息中记录：用户 ID、失败原因、时间、失败阶段等

```go
logger.Warn("短链删除失败", zap.Uint("userID", uid), zap.Error(err))
```

---

#### **5️⃣ 补偿机制（定时任务 + 死信队列）**

- 可用 cron 定时扫描“失败删除任务表”重新触发
- 或消费 Kafka 的死信队列重试
- 防止用户数据“删了一半”

---

### ✅ 总结成表格如下：

| 步骤 | 描述 | 是否异步 | 失败处理 |
|------|------|-----------|------------|
| 1. 删除用户记录 | MySQL 中删除/标记 | 否 | 返回错误 |
| 2. 删除 token 缓存 | Redis 中删除登录信息 | 否 | 忽略或重试 |
| 3. 删除短链资源 | Kafka 通知 shortlink 服务删除 | ✅ 是 | 写入失败记录 |
| 4. 记录删除失败 | 记录日志或失败任务表 | 否 | — |
| 5. 失败补偿 | 定时任务或 DLQ 重试 | ✅ 是 | 保证一致性 |

- Kafka `UserDeletedEvent` 定义
- 短链服务消费 `user-deleted` 的代码
- 删除失败记录表 + 定时补偿逻辑
