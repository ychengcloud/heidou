# 模板

`heidou` 接受 `--template` 参数，在 `heidou generate` 命令成功执行后，会生成一个与模板结构相同的目录。

## 目录结构

```console
.
├── skeleton  
└── templates
```

- `skeleton`，此目录下文件替换模板变量后会按原样复制到目标相同的目录下,以.tmpl 或 .tpl为后缀的文件，生成的文件会自动去除后缀。例如：

  `example/todo.go.tmpl` - 这个模板文件会生成 `example/todo.go` 的文件

- `templates`，根据数据表, 目录下文件会为数据表生成相应的文件。替换变量后的文件，会生成到配置文件中 `templates` 段指定的位置。例如 ：
  
  数据表有 `table1 table2`， 配置项中 `templates` 的目标路径配置为 `services`

    `service.tmpl` - 这个模板文件会生成 `services/table1.go services/table2.go` 两个文件。

```console
templates
    └──service.tmpl
    
```

## 官方模板

- [x] [graphql-server-template](templates-graphql-server.md) - 符合 Golang 设计哲学的工程框架,支持查询过滤、排序、分页、指定返回字段、批量查询、批量更新、批量删除。自动生成 Graphql 文档。
- [x] [react-admin-template](templates-react-admin.md)
- [ ] [restful-server-template]()
- [ ] [grpc-microservice-server]()

## 社区模板

等待你的添加

## 模板变量

- 根变量

| 名称        | 说明       | 类型   | 默认值 | 例子 |
| ----------- | ---------- | ------ | ------ | ---- |
| ProjectName | 项目名     | string | -      | -    |
| Extra       | 扩展信息   | map    | -      | -    |
| Tables      | 数据表信息 | array  | -      | -    |

- 数据表和字段都会导出的变量

| 名称                 | 说明     | 类型   | 默认值 | 例子           |
| -------------------- | -------- | ------ | ------ | -------------- |
| Name                 | 数据表名 | string | -      | product_table  |
| NameSnake            | 数据表名 | string | -      | product_table  |
| NameSnakePlural      | 数据表名 | string | -      | product_tables |
| NameCamel            | 数据表名 | string | -      | ProductTable   |
| NameCamelPlural      | 数据表名 | string | -      | ProductTables  |
| NameLowerCamel       | 数据表名 | string | -      | productTable   |
| NameLowerCamelPlural | 数据表名 | string | -      | productTables  |

- 数据表导出变量

| 名称            | 说明               | 类型   | 默认值 | 例子          |
| --------------- | ------------------ | ------ | ------ | ------------- |
| Name            | 数据表名           | string | -      | product_table |
| Description     | 数据表描述信息     | string | -      | -             |
| Extra           | 扩展配置           | map    | -      |
| Fields          | 字段数组           | array  | -      | -             |
| ErrorCodes      | 错误码             | array  | -      |
| PrimaryKeyField | 错误码             | Field  | -      |
| Filterable      | 是否有过滤字段     | bool   | false  | true          |
| Sortable        | 是否有排序字段     | bool   | false  | true          |
| HasErrorCode    | 是否有自定义错误码 | bool   | false  | true          |
| HasJoinField    | 是否有关联字段     | bool   | false  | true          |


- 字段导出变量

| 名称            | 说明                                                                              | 类型          | 默认值 | 例子      |
| --------------- | --------------------------------------------------------------------------------- | ------------- | ------ | --------- |
| Name            | 字段名                                                                            | string        | -      | id        |
| Alias           | 别名                                                                              | string        | -      | nameAlias |
| Description     | 描述信息                                                                          | string        | -      | -         |
| IsRequired      | 是否必填字段                                                                      | bool          | false  | true      |
| IsSortable      | 是否可按此字段排序                                                                | bool          | false  | true      |
| IsFilterable    | 是否可按此字段过滤                                                                | bool          | false  | true      |
| Operations      | 排序时的可用操作,取值 Eq,In,Gt,Gte,Lt,Lte,Contains,StartsWith,EndsWith,AND,OR,NOT | enum          | -      | true      | tags | 扩展 struct tags | string | "" | binding:"required,max=64" |
| IsPrimaryKey    | 是否主键                                                                          | bool          | false  | true      |
| IsForeignKey    | 是否外键                                                                          | bool          | false  | true      |
| IsAutoIncrement | 主键是否自增类型                                                                  | bool          | false  | true      |
| MetaType        | 字段类型元信息                                                                    | map           | -      |
| TagsHTML        | 字段标签信息                                                                      | template.HTML | -      |
| JoinTable       | 关联表                                                                            | Table         | -      | -         |
| MaxLength       | 字段最大长度                                                                      | number        | 0      | -         |
