# Heidou

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