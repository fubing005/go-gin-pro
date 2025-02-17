# 多功能 Gin 框架项目文档

---

## 项目概述
基于 **Golang Gin 框架**封装的双端（用户端 + 后台管理端）解决方案，集成以下核心能力：
- **认证鉴权**：JWT 双端隔离鉴权、Redis 黑名单、Token 续签
- **分布式支持**：Redis 集群、分布式锁、RabbitMQ/Kafka 消息队列
- **企业级组件**：RBAC 权限模型、多语言支持、日志收集、文件存储模块
- **高可用设计**：限流中间件、Panic 恢复、多数据库集成（MySQL/ClickHouse/MongoDB/ES）

---

## 功能架构

### 核心模块
| 模块                | 功能描述                                                                 |
|---------------------|------------------------------------------------------------------------|
| **认证体系**         | 双端独立 JWT 签发/验证、Redis 黑名单注销机制、Token 自动续签（滑动窗口）   |
| **权限管理**         | RBAC 模型（角色-权限绑定）、接口级访问控制、动态路由鉴权                  |
| **多语言**           | 请求级语言切换、验证器错误自动翻译、多语言文件动态加载                    |
| **文件管理**         | 统一存储接口（支持本地/S3/OSS）、上传下载限速、文件哈希校验               |
| **监控审计**         | 请求日志记录（UserAgent/IP 解析）、Zap 日志分级存储、ES 日志聚合          |

---

## 技术栈整合

### 基础框架
- **Web 框架**: Gin
- **配置管理**: Viper
- **日志系统**: Zap + Lumberjack（日志切割）
- **依赖管理**: Go Modules

### 数据服务
| 服务             | 用途                               | 集群方案              |
|------------------|-----------------------------------|---------------------|
| MySQL 8.0        | 核心业务数据存储                   | 主从复制            |
| Redis 7.0        | Token 黑名单/分布式锁/缓存         | Cluster 模式        |
| Elasticsearch 7.8.1| 业务日志/操作审计分析               | 单节点          |
| ClickHouse 21.3.20  | 统计报表/分析型查询                 | 单节点           |
| RabbitMQ 3.12  | 订单通知/异步任务处理                 | 3节点集群     |
| Kafka 3.6  | 日志实时采集/用户行为追踪                 | 单节点     |
| Mongodb 6.0  | 分布式数据存储/数据集成和同步/实时事件处理        | 单节点     |

---

## 项目结构说明

### 核心目录
```bash
├── app
    ├── common
        ├── request
            ├── request_admin
            ├── request_api
            common.go     #通用request  data
            validator.go  #自定义json返回数据
        ├── response
            ├── reponse_admin
            ├── reponse_api
            response.go   #通用response data
    ├── controllers #控制器
        ├── admin #后台管理系统控制器
        ├── api #api控制器
    ├── middleware
        ├── adminAuth.go #后台管理系统鉴权中间件
        ├── cors.go #跨域中间件 
        ├── jwt.go #jwt鉴权中间件
        ├── rateLimit.go # 限流中间件
        ├── recovery.go #自定义gin框架panic中间件
        ├── requestLogger.go #请求日志中间件
        ├── setLang.go #翻译中间件
    ├── models #数据模型
    ├── services #业务逻辑
├── bin #项目打包后的二进制文件
├── bootstrap #服务及配置初始化
    ├── clickhouse.go
    ├── config.go
    ├── cron.go
    ├── db.go
    ├── elasticsearch.go
    ├── initService.go
    ├── log.go
    ├── mongodb.go
    ├── redis.go
    ├── router.go
    ├── storage.go
    ├── validator.go
├── config #定义服务以及配置所需要的参数
├── core
  ├── errors #全局错误配置
  ├── trans #翻译
├── docs #swagger文档
├── global #调用服务及配置
    ├── app.go
    ├── kafka.go
    ├── lock.go
    ├── queue.go
    ├── rabbitmq.go
├── go-storage #文件存储模块
├── lang #翻译文件
   ├── en
   ├── zh-cn
├── routes
    ├── admin.go
    ├── api.go
├── static #静态文件
├── storage #上传文件以及日志的存储
├── tmp #gin框架运行的临时文件
├── utils #工具函数
    ├── bcrypt.go
    ├── directory.go
    ├── helper.go
    ├── md5.go
    ├── str.go
    ├── useragent.go
    ├── validator.go
├── config.yaml #配置文件
├── main.go #启动文件
├── Readme.md #项目文档
```

### 关键设计
1. **双端隔离**  
   - 独立 JWT 密钥与签发者标识（`iss` 声明）
   - 分离的路由分组（`routes/api.go` 与 `routes/admin.go`）
   - 差异化的中间件链（管理端强制操作日志记录）

2. **安全防护**
   - 请求频率限制（滑动窗口算法）

3. **TOKEN**
   - 基于 Redis 的 LUA 实现
   - 自动续期实现

### 本地运行项目
```bash
#进入新项目根目录
go mod tidy
go run main.go 或者 fresh
```

# 打包项目[windows运行环境]
```bash
   make          # 构建应用
   make run        # 构建并运行应用
   make clean      # 清理生成的二进制文件
   docker-deploy   # 将构建后的应用部署到docker
```

# 容器化部署
```bash
# 构建镜像并启动容器
docker-compose up [-d] --build
```