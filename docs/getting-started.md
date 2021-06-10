# 快速指南

## 安装

此命令会将 `heidou` 安装在 `/usr/local/bin`， 通过修改命令行参数，可安装在其他目录。

```console
curl -sfL https://raw.githubusercontent.com/ychengcloud/heidou/main/scripts/install.sh | sh -s -- -b /usr/local/bin
```

在安装 `heidou` 后，你应能在 `PATH` 中找到它。确认是否安装成功：

```console
heidou -v
```

## 创建项目

转到要生成项目的目录并运行：

```console
go run github.com/ychengcloud/heidou/cmd/heidou init todo
```

此命令将在当前目录生成名为 `todo` 的项目

```console

```