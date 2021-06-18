# API

基于 Golang & Graphql 实现 `xiaokeai` 的后端 `API` 服务。

## 依赖

参见 [模板依赖](templates-graphql-server?id=安装依赖)

## 创建项目

运行下列命令创建一个新的项目：

```bash
heidou init xkaapi

cd xkaapi

mv heidou-example.yml heidou.yml
```

此命令将在当前目录生成名为 `xkaapi` 的项目，并生成配置文件模板。

修改 `heidou.yml` 内容如下：

```yaml
# 项目名称
projectName: "xkaapi"

# 是否自动覆盖已有文件
overwrite: false

# 数据库类型
loader: "mysql"

# Extra 为用户自定义变量，以 map 结构读取，
# 特别注意： golang 中 map的key 不区分大小写，所以引用时全部以小写形式引用，例如 ： .Extra.pkgpath
extra:
  pkgPath: "github.com/ychengcloud/tutorial/xkaapi"

# Golang 的默认变量标识符与模板中的变量标识符相同时，需要修改成不同的
# delim:
#   left: "@@"
#   right: "@@"

# 数据库配置
db:
  dialect: "mysql"
  user: "<user>"
  password: "<password>"
  host: "127.0.0.1"
  port: 3306
  name: "xiaokeai"
  charset: "utf8mb4"

# 生成名称格式，用于指定 templates.nameFormat 配置项中动态替换部分的名称格式，默认使用数据库表名相同的格式，可选值如下：
# camel eg: VariableName
# lowerCamel eg: variableName
# camelPlural  eg: VariableNames
# lowerCamelPlural  eg: variableNames
# snake  eg: variable_name
tmplNameFormat: "lowerCamel"

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
  - name: "user"
    fields:
    - name: "name"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
      tags: 'binding:"required,max=64"'
    - name: group
      joinType: "ManyToMany"
      tableName: group  
      joinTableName: user_group
  - name: "group"
    fields:
    - name: "name"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: user
      joinType: "ManyToMany"
      tableName: user  
      joinTableName: user_group
  - name: "product"
    fields:
    - name: "name"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]  
  - name: "order"
    fields:
    - name: "id"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: "user_id"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: order_item
      joinType: "HasMany"
      tableName: order_item
  - name: "order_item"
    fields:
    - name: "order_id"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: order
      joinType: "BelongTo"
      tableName: order
    - name: product
      joinType: "HasOne"
      tableName: product
  - name: user_group
    isSkip: true  
```

## 生成代码

- 下载 `graphql-server` 模板

```bash
mkdir -p $HOME/.heidou/graphql-server-template
git clone https://github.com/ychengcloud/graphql-server-template $HOME/.heidou/graphql-server-template

```

- 基于模板和配置文件生成代码

```bash
heidou generate -t $HOME/.heidou/graphql-server-template/ -c ./heidou.yml
```

## 构建

```bash
# 执行 gqlgen
make gql_gen

# go modules
go mod tidy

# 依赖注入
make wire

# 构建
make build

```

## 运行前配置

```bash
mv configs/server-example.yml server.yml
```

## 运行前配置

```bash
mv configs/server-example.yml configs/server.yml
```

修改配置模板内容如下：

```yaml
app:
  name: "xkaapi"
http:
  mode: "debug"
  # mode: release
  host: 0.0.0.0
  port: 7779
  graphqlPath: "graphql"
  playgroundPath: "playground"
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
  dialect: "mysql"
  debug: true
  autoMigrate: false

  mysql:
    user: "<user>"
    password: "<password>"
    host: "127.0.0.1"
    port: 3306
    name: "xiaokeai"
    charset: "utf8mb4"

auth:
  skip: true

log:
  filename: "/tmp/xiaokeai.log"
  maxSize: 500
  maxBackups: 3
  maxAge: 3
  level: "debug"
  stdout: false
jaeger:
  serviceName: "admin"
  reporter:
    localAgentHostPort: "jaeger-agent:6831"
  sampler:
    type: "const"
    param: 1
jwt:
  signingKey: "YOUCHENG"
  issuer: "ycheng.pro"
  claimKey: "claim"
  signingMethod: "HS512"
  # seconds
  expired: 1000000

```

## 运行

```bash
make run
```

## 体验

浏览器打开 [Graphql Playground](http://localhost:7779/api/playground)
