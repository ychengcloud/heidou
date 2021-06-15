# HEIDOU

Heidou 是一个代码生成框架，用于高效生成各种业务代码，主要使用场景是生成工程框架和CRUD类代码，大大缩短业务从设计到上线的时间。在框架层面保证高度的灵活和可扩展性，不限制业务场景。

## 特性

- 解析数据库表结构，生成丰富的表述信息
- 支持自定义模板，结合不同模板快速生成不同业务场景的项目
- 支持 BelongTo 、HasOne 、HasMany 、ManyToMany 关联配置
- 支持配置校验信息
- 支持业务的灵活扩展
- 表结构更改可重复生成代码
- 使用简单，通过资源内嵌，实现一个二进制文件即可跨平台启动

## 文档

[中文](/README) | [English](/en_US/README)

## 模板

- 官方模板
  
  - [x] [graphql-server-template](https://github.com/ychengcloud/graphql-server-template) - 符合 Golang 设计哲学的工程框架,支持查询过滤、排序、分页、指定返回字段、批量查询、批量更新、批量删除。自动生成 Graphql 文档。
  - [ ] [react-antd-graphql-template]()