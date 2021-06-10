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

- `templates`，根据数据表, 目录下文件会为数据表生成相应的文件。替换变量后的文件，会生成到 [配置文件中指定的位置](docs/guide/config#templates) 。例如 ：
  
  数据表有 `table1 table2`， 配置项中 `templates` 的目标路径配置为 `services`

    `service.tmpl` - 这个模板文件会生成 `services/table1.go services/table2.go` 两个文件。

```console
templates
    └──service.tmpl
    
```


## 官方模板

### Server

- [ ] [graphql-golang-server](https://github.com/ychengcloud/graphql-golang-server-template) 符合 Golang 设计哲学的工程框架，基于 GqlGen、GORM、Gin, 包括基础功能(JWT, OpenTracing, ZapLog, Promtheus)

- [ ] [restful-golang-server](https://github.com/ychengcloud/restful-golang-server-template) 符合 Golang 设计哲学的工程框架，基于 GORM、Gin, 包括基础功能(JWT, OpenTracing, ZapLog, Promtheus)

- [ ] [grpc-microservice-server]()

### React

- [ ] [graphql-react-admin-typescript-vite](https://github.com/ychengcloud/graphql-react-admin-template) 


## 社区模板

等待你的添加
