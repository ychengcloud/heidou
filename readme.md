# Heidou

代码生成框架，CRUD归我，业务逻辑归你

## 特性

- 解析数据库表结构，生成丰富的表述信息
- 支持自定义模板，结合不同模板快速生成不同业务场景的项目
- 支持 BelongTo \ HasOne \ HasMany \ ManyToMany 关联配置
- 支持配置校验信息
- 支持业务扩展
- 表结构更改可重复生成代码
- 资源内嵌，一个二进制文件即可启动

## 依赖

- make
- golang

## 编译

    make

## 安装

    make install

## 生成项目

### 初始化项目

初始化后，在 <project> 下会生成样例配置文件 heidou-example.yaml ，根据业务需要修改 heidou.yaml

```shell
heidou init <project>

cd <project>

mv heidou-example.yaml heidou.yaml
```

### 生成项目代码

    heidou generate -t <path/to/template> [-c heidou.yml] 

## 生成配置项说明

数据表配置项 : tables

| 名称    | 说明                   | 类型   | 默认值                                                | 例子               |
| ------- | ---------------------- | ------ | ----------------------------------------------------- | ------------------ |
| name    | 数据表名               | string | -                                                     | product            |
| isSkip  | 是否生成相应代码       | bool   | false                                                 | false              |
| extra   | 扩展配置               | map    | -                                                     |
| fields  | 字段数组               | Fields | -                                                     | -                  |
| errorCodes | 错误码 | string array | - |
| methods | 支持的 Api 方法 | array  | ["list", "update", "delete", "bulkGet", "bulkDelete"] | ["list", "update"] |

字段配置项 : fields:

| 名称           | 说明                                                      | 类型   | 默认值 | 例子                      |
| -------------- | --------------------------------------------------------- | ------ | ------ | ------------------------- |
| name           | 字段名                                                    | string | -      | id                        |
| alias           | 别名                                                    | string | -      | nameAlias                        |
| isSkip     | 是否忽略此字段                                              | bool   | false  | true                      |
| isRequired     | 是否必填字段                                              | bool   | false  | true                      |
| isSortable   | 是否可按此字段排序                                            | bool   | false  | true                      |
| isFilterable   | 是否可按此字段过滤                                            | bool   | false  | true                      |
| operations   | 排序时的可用操作                                            | enum   | false  | true                      |
| tags           | 扩展 struct tags                                          | string | ""     | binding:"required,max=64" |
| joinType       | 关联类型,取值 None, BelongTo, HasOne, HasMany, ManyToMany | enum   | None   | ManyToMany                |
| tableName      | 指定关联表表名                                            | string | ""     | category                  |
| joinTableName  | joinType 为ManyToMany时，指定连接表表名                   | string | ""     | product_category_relation |
| foreignKey     | 指定外键，可选，默认使用拥有者的类型名加上主字段名        | string | ""     | CategoryID                |
| references     | 指定引用                                                  | string | ""     | name                      |
| joinForeignKey | 指定连接表的外键                                          | string | ""     | productReferID            |
| joinReferences | 指定连接表的引用外键                                      | string | ""     | productRefer              |


关联关系相关字段 (foreignKey, references, joinForeignKey, joinReferences）的配置与 Gorm 保持一致，详见：[Gorm](https://gorm.io/zh_CN/docs) 的关联说明



## 模板

### 目录结构

```console
.
└── root
    ├── skeleton  
    └── templates
```

- skeleton，此目录下文件替换模板变量后会原样复制到目标目录,以.tmpl 或 .tpl为后缀的文件，生成的文件会自动去除后缀

- template，此目录下文件会根据数据表，每个表相应生成一个文件

## 官方模板

### Server

- [ ] [graphql-golang-server](https://github.com/ychengcloud/graphql-golang-server-template) 符合 Golang 设计哲学的工程框架，基于 GqlGen、GORM、Gin, 包括基础功能(JWT, OpenTracing, ZapLog, Promtheus)

- [ ] [restful-golang-server](https://github.com/ychengcloud/restful-golang-server-template) 符合 Golang 设计哲学的工程框架，基于 GORM、Gin, 包括基础功能(JWT, OpenTracing, ZapLog, Promtheus)

- [ ] [grpc-microservice-server]()

### React

- [ ] [graphql-react-admin-typescript-vite](https://github.com/ychengcloud/graphql-react-admin-template) 


## 社区模板

等待你的添加

## 鸣谢

本项目的产生离不开这些优秀项目的启发， 如有遗漏欢迎补充指正

https://github.com/LyricTian/gin-admin

https://github.com/cmelgarejo/go-gql-server

https://github.com/8treenet/freedom

https://github.com/jinzhu/gorm

https://github.com/wantedly/apig

https://github.com/facebook/ent

https://github.com/libragen/felix

https://github.com/smallnest/gen

https://github.com/webliupeng/gin-tonic