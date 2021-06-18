# 管理系统

基于 React & Graphql 实现 `xiaokeai` 的前端管理系统。

## 依赖

- [ApiServer](tutorial-api.md) - 需要先启动 Api Server 才能自动生成前端 Graphql 层
- [模板依赖](templates-react-admin?id=安装依赖)

## 创建项目

运行下列命令创建一个新的项目：

```bash
heidou init xkaadmin

cd xkaadmin

mv heidou-example.yml heidou.yml
```

此命令将在当前目录生成名为 `xkaadmin` 的项目，并生成配置文件模板。

修改 `heidou.yml` 内容如下：

```yaml
# 项目名称
projectName: "xkaadmin"

# 是否自动覆盖已有文件
overwrite: false

# 数据库类型
loader: "mysql"

# Extra 为用户自定义变量，以 map 结构读取，
# 特别注意： golang 中 map的key 不区分大小写，所以引用时全部以小写形式引用，例如 ： .Extra.pkgpath
extra:
  skipFieldsInGql: 
    - "CreatedAt"
    - "UpdatedAt"
    - "DeletedAt"
  skipFieldsInTable: 
    - "CreatedAt"
    - "UpdatedAt"
    - "DeletedAt"
  skipFieldsInForm: 
    - "CreatedAt"
    - "UpdatedAt"
    - "DeletedAt"
    - "CreateBy"
    - "UpdateBy"

# Golang 的默认变量标识符与模板中的变量标识符相同时，需要修改成不同的
delim:
  left: "@@"
  right: "@@"

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
  - NameFormat: "src/pages/%s"
    Path:   "pages"
  - NameFormat: "src/mocks/%s.ts"
    Path:   "mock.ts.tmpl"
  - NameFormat: "src/stores/%s.ts"
    Path:   "store.ts.tmpl"
  - NameFormat: "src/graphql/%s.gql"
    Path:   "graphql.gql.tmpl"

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

- 下载 `react-admin` 模板

```bash
mkdir -p $HOME/.heidou/graphql-server-template
git clone https://github.com/ychengcloud/react-admin-template $HOME/.heidou/react-admin-template

```

- 基于模板和配置文件生成代码

```bash
heidou generate -t $HOME/.heidou/react-admin-template/ -c ./heidou.yml
```

## 生成 graphql 代码
根据 Graphql API 配置生成 Graphql 相关代码

```bash
yarn gql-gen
```

## 预览(Mock模式)

`Mock` 模式数据为本地模拟生成，不需要 Api 调用 ，基础功能测试通过后，再进行前后端 [联调](tutorial-integrated.md)

```bash
yarn mock 
```
