# Distributed API Quota Service

纯Go实现的分布式API流量策略与配额服务。当前默认使用内存适配器，便于本地零依赖启动；`internal/application`通过接口保留PostgreSQL、Redis和事件总线替换边界。数据库结构位于`migrations/`。

## 快速开始

```bash
go test ./...
go run ./cmd/api-quota-service
```

默认监听`:8084`，可通过`API_QUOTA_ADDRESS`、`API_QUOTA_SHUTDOWN_TIMEOUT`、`API_QUOTA_LOG_LEVEL`和`API_QUOTA_STORAGE_DRIVER`覆盖配置。

## 核心流程

1. `POST /v1/services`注册服务。
2. `POST /v1/policies`创建草稿策略。
3. `POST /v1/policies/{id}/publish`发布策略。
4. `POST /v1/decisions`执行策略决策，超限返回HTTP 429。
5. `GET /v1/events`查询超限事件。

支持`fixed_window`、`sliding_window`和`token_bucket`。策略可按服务、路径、HTTP方法、来源标签和环境匹配，按priority选择最高优先级策略。

## 运维端点

`/healthz`、`/readyz`和`/metrics`提供健康、就绪和Prometheus文本指标。服务响应SIGINT/SIGTERM并等待HTTP请求完成后退出。

## 目录

`cmd`启动入口；`internal/domain`领域模型和算法；`internal/application`用例；`internal/http`接口适配器；`internal/infrastructure`内存实现；`api` OpenAPI；`configs`配置示例；`migrations`数据库迁移；`deploy`容器文件；`scripts`本地检查脚本。

源码按职责拆分，非测试Go源码超过2000行目标时应继续扩展实际适配器和持久化实现，禁止通过重复代码凑行数。
