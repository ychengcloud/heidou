
```yaml
# 项目名称
projectName: api

# 是否自动覆盖已有文件
overwrite: false

# 数据库类型
loader: mysql

# Extra 为用户自定义变量，以 map 结构读取，golang 中 map的key 不区分大小写，所以引用时全部以小写形式引用，例如 ： .Extra.pkgpath
extra:
  pkgPath: github.com/ychengcloud/api
  skipFields: 
    - CreatedAt
    - DeletedAt
    - UpdateBy
    - UpdatedAt
    - CreateBy

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
  name: auth
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

tables:
  - name: user_group
    fields:
    - name: name
      isRequired: true
      tags: 'binding:"required,max=64"'
  - name: org
    fields:
    - name: id
      isFilterable: true
      operations: ["Eq", "In"]
    - name: name
      isRequired: true
      isFilterable: true
  - name: org_node
    typeName: tree
    fields:
    - name: name
      isRequired: true
      isFilterable: true
      operations: ["Eq", "In"]
    - name: parent_id
      isFilterable: true
      operations: ["Eq", "In"]
    - name: org
      joinType: "BelongTo"
      tableName: org
    - name: parent
      joinType: "BelongTo"
      tableName: org_node
      foreignKey: "ParentId"
  - name: perm_action
    fields:
    - name: namespace
      joinType: "BelongTo"
      tableName: perm_namespace
      foreignKey: "NamespaceId"
  - name: perm_resource
    fields:
    - name: namespace
      joinType: "BelongTo"
      tableName: perm_namespace
      foreignKey: "NamespaceId"
  - name: perm_namespace
    fields:
    - name: resource
      joinType: "HasMany"
      tableName: perm_resource
      foreignKey: "NamespaceId"
```