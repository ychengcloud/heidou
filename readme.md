# Heidou

## 理念 
提效管理后台,自动化生成必要代码，CRUD交给我，业务交给你

## 特性

- 生成符合 Golang 设计哲学的工程框架，包括基础功能(JWT, OpenTracing, ZapLog, Promtheus)
- 依赖注入
- 根据数据库生成对应 CRUD 方法
- 支持查询过滤、排序、分页、指定返回字段、批量查询、批量更新、批量删除
- 支持 BelongTo \ HasOne \ HasMany \ ManyToMany 关联配置
- 支持配置校验信息
- 支持 Swagger 文档自动生成
- 支持业务扩展
- 表结构更改可重复生成代码
- 框架代码自主可控，初始生成后可根据业务需要灵活修改
- 资源内嵌，一个二进制文件即可启动

- 安全特性
    - JWT
    - Casbin Permission
    - CSRF
    - XSS
    - CORS


### TODO
- 批量更新
- 针对管理系统和应用的权限配置
- docker-compose

## 依赖

- make
- golang

## 编译

    make build

## 安装

    make install

## 生成项目

### 初始化项目

    heidou init -p {PackageName} {ProjectName}

### 配置项目

初始化后，会生成样例配置文件 heidou.yaml，根据业务修改

    cd {ProjectName}
    vim heidou.yaml

### 生成项目代码
    
    heidou generate -c heidou.yaml

## 运行项目

项目生成后，根据业务情况修改项目配置文件，即可构建运行

### 修改项目配置文件

生成代码后，会生成样例项目配置文件 config/server-example.yaml，根据业务修改

    cp config/server-example.yaml config/server.yaml
    vim config/server.yaml

### 编译
    make build

### 启动
    make run

## 生成配置项说明

数据表配置项 : tables

| 名称 | 说明 | 类型 | 默认值 | 例子 |
|-|-|-|-|-|-|
|name | 数据表名 | string | - | product |
|isSkip| 是否生成相应代码 | bool | false | false |
|fields| 字段数组 | string | - | - |
|methods| 支持的Restful Api 方法 | array | ["list", "update", "delete", "bulkGet", "bulkDelete"] | ["list", "update"] |


表字段配置项 : fields:

| 名称 | 说明 | 类型 | 默认值 | 例子 |
|-|-|-|-|-|
|name | 字段名 | string | - | id |
|isRequired| 是否必填字段 | bool | false | true |
|isFilterable| 是否可过滤字段 | bool | false | true |
|tags| 扩展 struct tags | string | "" | binding:"required,max=64" |
|joinType| 关联类型,取值 None, BelongTo, HasOne, HasMany, ManyToMany   | enum | None | ManyToMany |
|tableName| 指定关联表表名 | string | "" | category |
|joinTableName| joinType 为ManyToMany时，指定连接表表名 | string | "" | product_category_relation |
|foreignKey| 指定外键，可选，默认使用拥有者的类型名加上主字段名 | string | "" | CategoryID |
|references| 指定引用 | string | "" | name |
|joinForeignKey| 指定连接表的外键 | string | "" | productReferID |
|joinReferences| 指定连接表的引用外键 | string | "" | productRefer |


关联关系相关字段 (foreignKey, references, joinForeignKey, joinReferences）的配置与 Gorm 保持一致，详见：[Gorm](https://gorm.io/zh_CN/docs) 的关联说明

### 样例文件： 
    
    {ProjectName}/heidou-example.yml

## 运行配置项说明

## 鸣谢

本项目的产生离不开这些优秀项目的启发， 如有遗漏欢迎补充指正

https://github.com/cmelgarejo/go-gql-server

https://github.com/8treenet/freedom

https://github.com/jinzhu/gorm

https://github.com/wantedly/apig

https://github.com/facebook/ent

https://github.com/libragen/felix

https://github.com/smallnest/gen

https://github.com/webliupeng/gin-tonic