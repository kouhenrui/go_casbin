# Go Casbin 项目

基于Go语言和Casbin的权限管理系统

## 项目结构

```
go_casbin/
├── api/                    # API层
│   ├── router.go          # 路由配置
│   └── api.go             # API接口定义
├── cmd/                   # 应用入口
│   └── main.go           # 主程序入口
├── configs/               # 配置文件
│   ├── config.dev.yaml   # 开发环境配置
│   ├── model.conf        # Casbin模型配置
│   └── policy.csv        # Casbin策略配置
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── controller/       # 控制器层
│   ├── dto/              # 数据传输对象
│   ├── do/               # 数据对象
│   ├── handler/          # 处理器
│   ├── logger/           # 日志管理
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   └── service/          # 业务逻辑层
├── pkg/                   # 公共包
│   ├── casbin/           # Casbin封装
│   ├── database/         # 数据库操作
│   ├── encrypt/          # 加密工具
│   ├── etcd/             # Etcd客户端
│   ├── jwt/              # JWT工具
│   ├── path/             # 路径工具
│   ├── redis/            # Redis客户端
│   └── util/             # 通用工具
├── k8s/                   # Kubernetes配置
│   └── deployment.yaml   # K8s部署文件
├── logs/                  # 日志文件
├── scripts/               # 脚本文件
├── test/                  # 测试文件
├── go.mod                # Go模块文件
└── go.sum                # Go依赖校验文件
```

## 快速开始

1. 构建镜像
```bash
docker build -t go-casbin:latest .
```

2. 部署到K8s
```bash
kubectl apply -f k8s/deployment.yaml
```

## 技术栈

- **框架**: Gin
- **数据库**: PostgreSQL + GORM
- **缓存**: Redis
- **权限**: Casbin
- **配置**: Viper
- **日志**: Zap
- **认证**: JWT
- **配置中心**: Etcd
