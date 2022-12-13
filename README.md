# Hello

ID: `com.aimerneige.hello`

Module for [MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)

## 功能

1. 收到好友请求或加群请求时自动通过，并发送指定消息。
2. 当黑名单用户尝试发送加群或好友请求时会自动拒绝。

## 使用方法

在适当位置引用本包

```go
package example

imports (
    // ...

    _ "github.com/yukichan-bot-module/MiraiGo-module-hello"

    // ...
)

// ...
```

在你的 `application.yaml` 里填入配置：

```yaml
aimerneige:
  hello:
    path: "./hello.yaml" # 配置文件路径，未设置默认为 `./hello.yaml`
```

编辑你的配置文件：

```yaml
blacklist: # 这里填黑名单用户
  - 1781924496
group:
  msg: "群聊的消息"
  image: # 如果需要发送图片可以像下面这样使用图片的路径
    - "/path/to/your/image/sample-1.png"
    - "/path/to/your/image/sample-2.png"
friend:
  msg: "私聊的消息"
  image: # 当然你可以选择不添加图片
```

## LICENSE

<a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
<img src="https://www.gnu.org/graphics/agplv3-155x51.png">
</a>

本项目使用 `AGPLv3` 协议开源，您可以在 [GitHub](https://github.com/yukichan-bot-module/MiraiGo-module-hello) 获取本项目源代码。为了整个社区的良性发展，我们强烈建议您做到以下几点：

- **间接接触（包括但不限于使用 `Http API` 或 跨进程技术）到本项目的软件使用 `AGPLv3` 开源**
- **不鼓励，不支持一切商业使用**
