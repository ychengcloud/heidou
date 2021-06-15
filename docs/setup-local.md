# 安装

## 本地安装

此命令会将 `heidou` 安装在 `/usr/local/bin`， 也可通过修改命令行参数安装在其他目录。

请打开终端/控制台窗口，输入如下命令：

```bash
curl -sfL https://raw.githubusercontent.com/ychengcloud/heidou/main/scripts/install.sh | sh -s -- -b /usr/local/bin
```

在安装 `heidou` 后，你应能在 `PATH` 中找到它。确认是否安装成功：

```bash
heidou -v
```