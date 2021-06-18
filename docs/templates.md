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

