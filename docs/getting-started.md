# 快速上手

本教程将构建一个基于 `Golang` 的 `Todo` 后端服务，提供标准的 `Graphql API` 以供前端应用调用。在我们开始之前，请确保您的机器上满足了以下前提条件。

## 前提条件

- [Make](setup-make.md)
- [Golang 1.16+](https://golang.org/doc/install)
- [Mysql 8.0+](https://dev.mysql.com/doc/refman/8.0/en/installing.html)
- [Heidou 0.1.10+](setup-local.md) 
- [Wire](https://github.com/google/wire)
- [gowatch](https://github.com/silenceper/gowatch) (可选)

## 数据字典

设计应用的数据字典 `todo.sql` 如下：

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

导入 `Mysql`

```bash
mysql -h <host> -u <user> < mysql.sql
```
## 创建项目

运行下列命令创建一个新的项目：

```bash
heidou init todo

cd todo

mv heidou-example.yml heidou.yml
```

此命令将在当前目录生成名为 `todo` 的项目，并生成配置文件模板。


修改 `heidou.yml` 内容如下：

```yaml
# 项目名称
projectName: "todo"

# 是否自动覆盖已有文件
overwrite: false

# 数据库类型
loader: "mysql"

# Extra 为用户自定义变量，以 map 结构读取，golang 中 map的key 不区分大小写，所以引用时全部以小写形式引用，例如 ： .Extra.pkgpath
extra:
  pkgPath: "github.com/ychengcloud/todo"

# Golang 的默认变量标识符与模板中的变量标识符相同时，需要修改成不同的
#delim:
#  left: "@@"
#  right: "@@"

# 数据库配置
db:
  dialect: "mysql"
  user: "<user>"
  password: "<password>"
  host: "127.0.0.1"
  port: 3306
  name: "todo"
  charset: "utf8mb4"

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
  - name: "todos"
    fields:
    - name: "id"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: "title"
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
  
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

修改配置模板内容如下：

```yaml
app:
  name: "todo"
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
    name: "todo"
    charset: "utf8mb4"

auth:
  skip: true

log:
  filename: "/tmp/todo.log"
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

#### 创建 Todo

```graphql
# Create
mutation ($input: TodoInput!) {
  todoCreate (input: $input) {
    todo {
      id
      title
      completed
    }
  }
}


# Variables
{
  "input": {
    "title": "My Todo 1",
    "completed": 0
	}
}
```

#### 查询 Todo
```graphql
# Query
query {
  todosOffsetBased {
    edges {
      node {
        id
        title
        completed
      }
      
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
    }
  }
}
```

#### 更新 Todo
```graphql
# Update
mutation($ids: [ID]!, $input: TodoInput!) {
  todoUpdate(ids: $ids, input: $input) {
    count
  }
}

# Variables
{
  "ids": ["595cc895-ca31-4c6f-9897-d376b1ec8bb2"],
  "input": {
    "title": "My Todo 1",
    "completed": 1
	}
}

```

#### 删除 Todo
```graphql
# Delete
mutation {
  todoDelete(ids: ["595cc895-ca31-4c6f-9897-d376b1ec8bb2"]) {
    count
  }
}
```