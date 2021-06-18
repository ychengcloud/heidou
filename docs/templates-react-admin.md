# React Admin

全新技术栈的后台管理系统模板，结合 [Heidou](https://github.com/ychengcloud/heidou) 使用，自动生成具有CRUD功能的项目框架。

This Starter utilizes React, Recoil, React Graphql, React Hooks, Typescript And Vite.

[项目地址](https://github.com/ychengcloud/react-admin-template)

## 特性

- 全新技术栈
- 根据数据库生成对应 CRUD 功能
- 支持 BelongTo 、HasOne 、HasMany 、ManyToMany 关联配置
- 根据 Graphql 文档自动生成前端 API 调用代码。
- 开发阶段支持 Mock 和 Proxy Server 两种模式
- 框架代码自主可控，初始生成后可根据业务需要灵活修改

## 安装依赖

- [Yarn](https://yarnpkg.com/)
- [Node](http://nodejs.org/)
- [Heidou 0.1.10+](setup-local.md) 

## 项目生成

参见 [Heidou](https://docs.ycheng.pro/heidou)

## 安装模块

```bash
yarn 
```

## 生成 graphql 代码
根据 Graphql API 配置生成 Graphql 相关代码

```bash
yarn gql-gen
```
## 开发(Mock模式)

```bash
yarn mock 
```

## 开发(Proxy Server模式)

```bash
yarn dev 
```

## 构建

```bash
yarn build
```

## 发布
```bash
yarn dist
```

