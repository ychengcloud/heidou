# TodoMvc 

## 后端

### 创建数据表
mysql -h 127.0.0.1 -u root todo < db/schema.sql

mysql -h 127.0.0.1 -u root todo < db/table.sql

### 生成项目代码

heidou-linux-amd64 init -p github.com/decker502/heidou/example/todo backend
cd backend
heidou-linux-amd64 generate -c configs/server.yml

## 前端
样式使用

https://github.com/tastejs/todomvc-app-css

npm install todomvc-app-css