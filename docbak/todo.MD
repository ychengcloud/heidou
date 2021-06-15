# TodoMvc 

## 后端

### 创建数据表
mysql -h 127.0.0.1 -u root todo < db/schema.sql

mysql -h 127.0.0.1 -u root todo < db/table.sql

### 初始化项目

heidou init -p github.com/ychengcloud/heidou/example/todo backend


### 修改 heidou 配置文件， 生成项目代码
cd backend
cp heidou-example.yml heidou.yml
heidou generate -c heidou.yml

### 修改项目配置文件
config/server.yaml

### 编译
make build

### 启动
make run

## 前端

### Prerequisites
To complete this tutorial, you will need:

Node.js installed locally, which you can do by following [How to Install Node.js and Create a Local Development Environment](https://www.digitalocean.com/community/tutorial_series/how-to-install-node-js-and-create-a-local-development-environment).
Some familiarity with React will be beneficial, but is not required. You can take a look at our [How To Code in React.js series](https://www.digitalocean.com/community/tutorial_series/how-to-code-in-react-jscd react-to-do
).

样式使用

https://github.com/tastejs/todomvc-app-css

npm install todomvc-app-css

网络请求库
npm install axios classnames react-error-boundary react-query recoil todomvc-app-css