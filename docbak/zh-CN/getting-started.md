# 快速指南

## 安装

此命令会将 `heidou` 安装在 `/usr/local/bin`， 通过修改命令行参数，可安装在其他目录。

```console
curl -sfL https://raw.githubusercontent.com/ychengcloud/heidou/main/scripts/install.sh | sh -s -- -b /usr/local/bin
```

在安装 `heidou` 后，你应能在 `PATH` 中找到它。确认是否安装成功：

```console
heidou -v
```

## 数据字典

```sql

-- Create a database
CREATE DATABASE `todo` DEFAULT CHARACTER SET = `utf8mb4`;

USE `todo`;
DROP TABLE IF EXISTS `todos`;
CREATE TABLE `todos` (
  `id` char(36) NOT NULL,
  `title` varchar(32) NOT NULL DEFAULT '' COMMENT '标题',
  `completed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否完成',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT = 'Todo';


```

## 创建项目

此命令将在当前目录生成名为 `todo` 的项目：

```console
heidou init todo

cd todo

mv heidou-example.yml heidou.yml
```

修改 `heidou.yml` 内容

```yaml
# 项目名称
projectName: todo

# 是否自动覆盖已有文件
overwrite: false

# 数据库类型
loader: mysql

# Extra 为用户自定义变量，以 map 结构读取，golang 中 map的key 不区分大小写，所以引用时全部以小写形式引用，例如 ： .Extra.pkgpath
extra:
  pkgPath: github.com/ychengcloud/todo

# Golang 的默认变量标识符与模板中的变量标识符相同时，需要修改成不同的
#delim:
#  left: "@@"
#  right: "@@"

# 数据库配置
db:
  dialect: mysql
  user: root
  password: ""
  host: "127.0.0.1"
  port: 3306
  name: todo
  charset: utf8mb4

# NameFormat 目标路径
# Path 模板路径名，以 templates为相对路径
templates:
  - nameFormat: "internal/generated/models/%s.go"
    path:   "models/model.go.tmpl"
  - nameFormat: "internal/resolvers/%s.go"
    path:   "resolvers/resolver.go.tmpl"
  - nameFormat: "internal/generated/schemas/%s.gql"
    path:   "schemas/schema.gql.tmpl"
  - nameFormat: "internal/generated/services/%s.go"
    path:   "services/service.go.tmpl"

# 数据表配置
tables:
  - name: todos
    fields:
    - name: id
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: title
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
  
```

## 代码生成

    heidou generate -t ../contrib/graphql-server-template/ -c ./heidou.yml

## 构建

```console
make gql_gen
go mod tidy
make wire
make build

```

## 运行前配置

    mv configs/server-example.yml server.yml

```yaml
app:
  name: todo
http:
  mode: debug
  # mode: release
  host: 0.0.0.0
  port: 7779
  graphqlPath: graphql
  playgroundPath: playground
  isPlaygroundEnabled: true
  allowOrigins:
    - "*"
  allowMethods: 
    - "PUT"
    - "GET"
    - "POST"
    - "HEAD"
    - "PATCH"
    - "OPTIONS"
    - "DELETE"
  allowHeaders:
    - "*"
db:
  dialect: mysql
  debug: true
  autoMigrate: false

  mysql:
    user: root
    password: ""
    host: "127.0.0.1"
    port: 3306
    name: todo
    charset: utf8mb4
  sqlite:
    name: "gorm.db"

auth:
  authType: kratos
  skip: true
  
services:
  kratosHost: "hengha.ycheng.pro:4455"

log:
  filename: /tmp/.log
  maxSize: 500
  maxBackups: 3
  maxAge: 3
  level: "debug"
  stdout: false
jaeger:
  serviceName: admin
  reporter:
    localAgentHostPort: "jaeger-agent:6831"
  sampler:
    type: const
    param: 1
jwt:
  # dd if=/dev/urandom bs=1 count=32 2>/dev/null | base64 -w 0 | rev | cut -b 2- | rev
  # signingKey: GRuHhzxQm7z0H7jFBHxd0x2UEjvJHgt+286nnJCOHYw
  signingKey: YOUCHENG
  issuer: ycheng.pro
  claimKey: claim
  signingMethod: HS512
  # seconds
  expired: 1000000

```

## 运行

    make run

## 打开 Graphql Playground

    open http://localhost:7779/api/playground