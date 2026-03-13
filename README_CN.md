# agent-rss

[English](./README.md)

面向 AI Agent 的 CLI RSS 工具。订阅 RSS 源，拉取 RSS/Atom 内容，支持按时间和关键字过滤。

## 安装

### npm

```bash
npm install -g @atopos31/agent-rss
```

### Go

```bash
go install github.com/atopos31/agent-rss/cmd/agent-rss@latest
```

### 从源码构建

```bash
git clone https://github.com/atopos31/agent-rss.git
cd agent-rss
go build -o agent-rss ./cmd/agent-rss
```

## 使用方法

### 订阅管理

```bash
# 添加订阅
agent-rss add hn https://news.ycombinator.com/rss

# 列出所有订阅
agent-rss list

# 获取指定订阅
agent-rss get hn

# 更新订阅
agent-rss update hn --src https://news.ycombinator.com/rss

# 删除订阅
agent-rss remove hn
```

### 拉取 RSS

```bash
# 拉取指定订阅
agent-rss fetch --name hn

# 拉取所有订阅
agent-rss fetch --all

# 输出为 JSON 数组（默认为 NDJSON）
agent-rss fetch --all --format json
```

### 过滤

```bash
# 按时间范围过滤
agent-rss fetch --all --since 2026-03-12
agent-rss fetch --all --since 2026-03-12T08:00:00+08:00 --until 2026-03-12T18:00:00+08:00

# 按标题关键字过滤
agent-rss fetch --all --title "AI"

# 按内容关键字过滤
agent-rss fetch --all --content "机器学习"

# 组合过滤
agent-rss fetch --all --since 2026-03-12 --title "Go" --title "Rust"
```

### 全局选项

```bash
# 使用自定义订阅文件
agent-rss --file /path/to/feeds.txt list
```

## AI Agent 使用最佳实践

许多 AI Agent 环境（如 Claude Code、OpenClaw 等）对 bash 命令的**输出大小有限制**。当拉取包含大量内容的 RSS 时，输出可能会被截断。

**推荐做法：** 将输出写入文件，然后使用 Agent 的文件读取功能获取完整内容。

```bash
# 将 RSS 输出写入临时文件
agent-rss fetch --all --since 2026-03-12 > /tmp/rss-output.json

# 然后使用 Agent 的 Read 工具读取完整内容
# Agent 可以无限制地读取 /tmp/rss-output.json
```

这种方式可以确保：
- RSS 内容不会被截断
- 完整访问所有拉取的条目
- 更好地处理大型订阅源

## 订阅文件格式

订阅存储在 `~/.config/agent-rss/feeds.txt`：

```
# 以 # 开头的是注释
hn https://news.ycombinator.com/rss
golang https://blog.golang.org/feed.atom
```

## 输出格式

### NDJSON（默认）

```json
{"name":"hn","src":"https://...","time":"2026-03-12T15:30:00+08:00","title":"...","content":"...","link":"...","id":"..."}
{"name":"hn","src":"https://...","time":"2026-03-12T14:20:00+08:00","title":"...","content":"...","link":"...","id":"..."}
```

### JSON

```json
[
  {"name":"hn","src":"https://...","time":"2026-03-12T15:30:00+08:00","title":"...","content":"...","link":"...","id":"..."},
  {"name":"hn","src":"https://...","time":"2026-03-12T14:20:00+08:00","title":"...","content":"...","link":"...","id":"..."}
]
```

## 寻找 RSS 源

想找 RSS 源订阅？查看 [awesome-rsshub-routes](https://github.com/JackyST0/awesome-rsshub-routes)，这里有各类 RSS 源的精选列表。

## 许可证

MIT
